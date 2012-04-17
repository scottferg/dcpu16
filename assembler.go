package main

import (
    "os"
    "io/ioutil"
    "regexp"
    "strings"
    "strconv"
    "encoding/binary"
    "bytes"
    "fmt"
    "encoding/hex"
)

var (
    labels map[string]int

    buffer = new(bytes.Buffer)

    LABEL               string = "^:[A-Za-z]+[^ ]"
    HEX                 string = "^0x[0-9a-zA-Z]+$"
    MEMORY              string = "^\\[0x[0-9a-zA-Z]+\\]$"
    MEMORY_AND_REGISTER string = "^\\[0x[0-9a-zA-Z]+\\+[A-Za-z]+\\]$"
    LABEL_AND_REGISTER  string = "^\\[[A-Za-z]+\\+[A-Za-z]+\\]$"
)

func valueFromHex(operand string) (op Word) {
    value := operand[2:]

    if len(value) % 2 == 1 {
        value = fmt.Sprintf("0%s", value)
    }

    literal, err := hex.DecodeString(value)
    if err != nil {
        fmt.Println(err)
    }

    switch {
    case len(literal) > 1:
        op = (Word(literal[0]) << 8) + Word(literal[1])
    default:
        op = Word(literal[0])
    }

    return
}

func getOperand(operand string) (op Word, nextWord Word) {

    if _, ok := labels[operand]; ok {
        return 0x1f, Word(labels[operand])
    }

    memoryRegisterMatched, _ := regexp.MatchString(MEMORY_AND_REGISTER, operand)
    memoryLabelMatched, _ := regexp.MatchString(LABEL_AND_REGISTER, operand)

    if memoryRegisterMatched || memoryLabelMatched {
        switch strings.ToUpper(operand[len(operand) - 2:len(operand) - 1]) {
        case "A":
            op = 0x10
        case "B":
            op = 0x11
        case "C":
            op = 0x12
        case "X":
            op = 0x13
        case "Y":
            op = 0x14
        case "Z":
            op = 0x15
        case "I":
            op = 0x16
        case "J":
            op = 0x17
        }

        hex, _ := regexp.Compile("0x[0-9a-zA-Z]+")
        label, _ := regexp.Compile("[A-Za-z]+[^+]")

        var match string
        if match = hex.FindString(operand); len(match) > 0 {
            nextWord = valueFromHex(match)
        } else if matched, _ := regexp.MatchString(operand, LABEL_AND_REGISTER); matched {
            match = label.FindString(operand)
            nextWord = Word(labels[match])
        }

        return
    }

    if matched, _ := regexp.MatchString(MEMORY, operand); matched {
        return 0x1e, valueFromHex(strings.Trim(operand, "[]"))
    }

    switch strings.ToUpper(operand) {
    case "A":
        op = 0x00
    case "B":
        op = 0x01
    case "C":
        op = 0x02
    case "X":
        op = 0x03
    case "Y":
        op = 0x04
    case "Z":
        op = 0x05
    case "I":
        op = 0x06
    case "J":
        op = 0x07
    case "[A]":
        op = 0x08
    case "[B]":
        op = 0x09
    case "[C]":
        op = 0x0a
    case "[X]":
        op = 0x0b
    case "[Y]":
        op = 0x0c
    case "[Z]":
        op = 0x0d
    case "[I]":
        op = 0x0e
    case "[J]":
        op = 0x0f
    case "POP":
        op = 0x18
    case "PEEK":
        op = 0x19
    case "PUSH":
        op = 0x1a
    case "SP":
        op = 0x1b
    case "PC":
        op = 0x1c
    case "O":
        op = 0x1d
    default:
        var literal Word

        if matched, _ := regexp.MatchString(HEX, operand); matched {
            literal = valueFromHex(operand)

            if literal < 32 {
                op = literal + 32
            } else {
                nextWord = literal
                op = 0x1f
            }
        } else {
            value, _ := strconv.Atoi(operand)
            literal = Word(value)

            if literal < 32 {
                op = literal + 32
            } else {
                nextWord = literal + 32
                op = 0x1f
            }
        }
    }

    return
}

func getOpcode(opcode string) (op Word, e error) {
    switch strings.ToUpper(opcode) {
    case "SET":
        op = 0x1
    case "ADD":
        op = 0x2
    case "SUB":
        op = 0x3
    case "MUL":
        op = 0x4
    case "DIV":
        op = 0x5
    case "MOD":
        op = 0x6
    case "SHL":
        op = 0x7
    case "SHR":
        op = 0x8
    case "AND":
        op = 0x9
    case "BOR":
        op = 0xa
    case "XOR":
        op = 0xb
    case "IFE":
        op = 0xc
    case "IFN":
        op = 0xd
    case "IFG":
        op = 0xe
    case "IFB":
        op = 0xf
    }

    return
}

func parseDatString(dat string) {
    fields := strings.Split(dat, ", ")

    var color Word
    writeString := func(str string) {
        for _, ch := range str {
            if color != 0 {
                buffer.WriteByte(byte(color))
            } else {
                buffer.WriteByte(0x00)
            }

            buffer.WriteByte(byte(ch))
        }
    }

    for _, field := range fields {
        field = strings.Trim(field, "\"")
        if matched, _ := regexp.MatchString(HEX, field); matched {
            color = valueFromHex(field) << 7
        } else if field == "0" {
            buffer.WriteByte(byte(0x00))
            buffer.WriteByte(byte(0x00))
        } else {
            writeString(field)
        }
    }
}

func getWord(word Word) ([]byte) {
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.BigEndian, word); err != nil {
        fmt.Println("binary.Write failed:", err)
    }

    return buf.Bytes()
}

func writeToBuffer(ins Word, nextA Word, nextB Word) {
    buffer.WriteByte(byte(ins >> 8))
    buffer.WriteByte(byte(ins))

    if nextA > 0 {
        buffer.Write(getWord(nextA))
    }

    if nextB > 0 {
        buffer.Write(getWord(nextB))
    }
}

func isMultiwordOperand(operand string) (multi bool) {
    op, _ := getOperand(operand)
    switch op {
    case 0x1e:
        multi = true
    case 0x1f:
        multi = true
    default:
        multi = false
    }

    if matched, _ := regexp.MatchString(MEMORY_AND_REGISTER, operand); matched {
        multi = true
    }

    if matched, _ := regexp.MatchString(LABEL_AND_REGISTER, operand); matched {
        multi = true
    }

    return
}

func scanForLabelAddresses(source []string) {
    wordCount := 0
    for _, sourceLine := range source {
        line := strings.Split(strings.Trim(sourceLine, " "), ";")[0]

        if len(line) == 0 {
            continue
        }

        var label string
        expression, _ := regexp.Compile(LABEL)
        if match := expression.FindString(line); match != "" {
            index := strings.Index(line, " ")

            if index == -1 {
                label = line[1:]
            } else {
                label = line[1:index]
            }

            line = strings.Trim(strings.Replace(line, match, "", -1), " ")

            labels[label] = wordCount
        }

        if len(line) == 0 {
            continue
        }

        countDatLine := func(dat string) {
            fields := strings.Split(dat, "\", ")
            for _, field := range fields {
                field = strings.Trim(field, "\"")
                if matched, _ := regexp.MatchString(HEX, string(field)); !matched {
                    wordCount += len(string(field))
                }
            }
        }

        if strings.ToLower(line[:3]) == "dat" {
            countDatLine(strings.Trim(line[3:], " "))

            continue
        }

        operands := strings.Split(line[3:], ", ")

        aaaa := strings.Trim(operands[0], " ")
        if isMultiwordOperand(aaaa) {
            wordCount++
        }

        if len(operands) > 1 {
            bbbb := strings.Trim(operands[1], " ")
            if isMultiwordOperand(bbbb) {
                wordCount++
            }
        }

        wordCount++
    }
}

func scanForLabels(source []string) {
    for _, sourceLine := range source {
        line := strings.Split(strings.Trim(sourceLine, " "), ";")[0]

        if len(line) == 0 {
            continue
        }

        var label string
        if matched, _ := regexp.MatchString(LABEL, line); matched {
            index := strings.Index(line, " ")

            if index == -1 {
                label = line[1:]
            } else {
                label = line[1:index]
            }

            labels[label] = 0x0
        }
    }
}

func main() {
    labels = make(map[string]int)

    if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
        source := strings.Split(string(contents), "\n")

        scanForLabels(source)
        scanForLabelAddresses(source)

        lineNumber := 1
        for _, sourceLine := range source {
            line := strings.Split(strings.Trim(sourceLine, " "), ";")[0]

            if len(line) == 0 {
                continue
            }

            // Shave off the label and whitespace
            expression, _ := regexp.Compile(LABEL)
            if match := expression.FindString(line); len(match) > 0 {
                line = strings.Trim(strings.Replace(line, match, "", -1), " ")
            }

            if len(line) == 0 {
                continue
            }

            opcode := line[:3]

            op, _ := getOpcode(opcode)

            if opcode == "dat" {
                parseDatString(strings.Trim(line[3:], " "))
                continue
            }

            operands := strings.Split(line[3:], ", ")

            aaaa := strings.Trim(operands[0], " ")
            a, nextAWord := getOperand(aaaa)

            if opcode == "JSR" {
                op = 0x01 << 4
                a = 0x1f << 10

                writeToBuffer(op + a, nextAWord, 0)

                continue
            }

            bbbb := strings.Trim(operands[1], " ")
            b, nextBWord := getOperand(bbbb)

            instruction := op + (a << 4) + (b << 10)

            // fmt.Printf("%s: 0x%x\n", line, instruction)

            writeToBuffer(instruction, nextAWord, nextBWord)

            lineNumber++
        }

        ioutil.WriteFile("output.dcpx", buffer.Bytes(), 0666)
    }
}
