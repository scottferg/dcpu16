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
    MEMORY_AND_REGISTER string = "^\\[0x[0-9a-zA-Z]+\\+[A-Z]+\\]$"
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
        fmt.Println("Got a label: " + operand)
        return 0x1f, Word(labels[operand])
    }

    if matched, _ := regexp.MatchString(MEMORY_AND_REGISTER, operand); matched {
        switch operand[len(operand) - 2:len(operand) - 1] {
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

        expression, _ := regexp.Compile("0x[0-9a-zA-Z]+")

        var match string
        if match = expression.FindString(operand); len(match) > 0 {
            nextWord = valueFromHex(match)
        }

        return
    }

    if matched, _ := regexp.MatchString(MEMORY, operand); matched {
        return 0x1e, valueFromHex(strings.Trim(operand, "[]"))
    }

    switch operand {
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
    switch opcode {
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

func isMultiwordOperand(operand Word) (multi bool) {
    switch operand {
    case 0x1e:
        multi = true
    case 0x1f:
        multi = true
    default:
        multi = false
    }

    return
}

func scanForLabels(source []string) {
    wordCount := 0
    for _, sourceLine := range source {
        line := strings.Split(strings.Trim(sourceLine, " "), ";")[0]

        if len(line) == 0 {
            continue
        }

        operands := strings.Split(line[3:], ", ")

        aaaa := strings.Trim(operands[0], " ")
        a, _ := getOperand(aaaa)

        if isMultiwordOperand(a) {
            fmt.Println("Operand: " + aaaa)
            wordCount++
        }

        if len(operands) > 1 {
            bbbb := strings.Trim(operands[1], " ")
            b, _ := getOperand(bbbb)

            if isMultiwordOperand(b) {
                fmt.Println("Operand: " + bbbb)
                wordCount++
            }
        }

        if matched, _ := regexp.MatchString(LABEL, line); matched {
            index := strings.Index(line, " ")
            label := line[1:index]

            labels[label] = wordCount
        }

        wordCount++
    }
}

func main() {
    labels = make(map[string]int)

    if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
        source := strings.Split(string(contents), "\n")

        scanForLabels(source)
        scanForLabels(source)

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

            opcode := line[:3]

            op, _ := getOpcode(opcode)

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

            writeToBuffer(instruction, nextAWord, nextBWord)

            lineNumber++
        }

        ioutil.WriteFile("output.dasm", buffer.Bytes(), 0666)
    }
}
