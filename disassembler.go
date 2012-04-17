package main

import (
    "io/ioutil"
    "os"
    "fmt"
)

var (
    program []Word
    PC = 0
)

func processOperand(operand Word) (value string) {
    switch operand {
    case 0:
        value = "A"
    case 1:
        value = "B"
    case 2:
        value = "C"
    case 3:
        value = "X"
    case 4:
        value = "Y"
    case 5:
        value = "Z"
    case 6:
        value = "I"
    case 7:
        value = "J"
    case 8:
        value = "[A]"
    case 9:
        value = "[B]"
    case 10:
        value = "[C]"
    case 11:
        value = "[X]"
    case 12:
        value = "[Y]"
    case 13:
        value = "[Z]"
    case 14:
        value = "[I]"
    case 15:
        value = "[J]"
    case 16:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "+A]"
        PC++
    case 17:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "+B]"
        PC++
    case 18:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "+C]"
        PC++
    case 19:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "+X]"
        PC++
    case 20:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "+Y]"
        PC++
    case 21:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "+Z]"
        PC++
    case 22:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "+I]"
        PC++
    case 23:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "+J]"
        PC++
    case 24:
        value = "POP"
    case 25:
        value = "PEEK"
    case 26:
        value = "PUSH"
    case 27:
        value = "SP"
    case 28:
        value = "PC"
    case 29:
        value = "O"
    case 30:
        value = "[" + fmt.Sprintf("0x%x", program[PC]) + "]"
        PC++
    case 31:
        value = fmt.Sprintf("0x%x", program[PC])
        PC++
    default:
        value = fmt.Sprintf("0x%x", operand)
    }

    return
}

func readInstruction(instruction Word) {
    PC++

    opcode := instruction & 0xf
    a := (instruction >> 4) & 0x3f
    b := (instruction >> 10) & 0x3f

    var aaaa string
    var bbbb string
    if opcode != 0 {
        aaaa = processOperand(a)
        bbbb = processOperand(b)
    }

    switch opcode {
    case 0:
        opcode, a = a, b

        switch instruction {
        case 1:
            fmt.Printf("JSR label")
        }
    case 1:
        fmt.Printf("SET %s %s\n", aaaa, bbbb)
    case 2:
        fmt.Printf("ADD %s %s\n", aaaa, bbbb)
    case 3:
        fmt.Printf("SUB %s %s\n", aaaa, bbbb)
    case 4:
        fmt.Printf("MUL %s %s\n", aaaa, bbbb)
    case 5:
        fmt.Printf("DIV %s %s\n", aaaa, bbbb)
    case 6:
        fmt.Printf("MOD %s %s\n", aaaa, bbbb)
    case 7:
        fmt.Printf("SHL %s %s\n", aaaa, bbbb)
    case 8:
        fmt.Printf("SHR %s %s\n", aaaa, bbbb)
    case 9:
        fmt.Printf("AND %s %s\n", aaaa, bbbb)
    case 10:
        fmt.Printf("BOR %s %s\n", aaaa, bbbb)
    case 11:
        fmt.Printf("XOR %s %s\n", aaaa, bbbb)
    case 12:
        fmt.Printf("IFE %s %s\n", aaaa, bbbb)
    case 13:
        fmt.Printf("IFN %s %s\n", aaaa, bbbb)
    case 14:
        fmt.Printf("IFG %s %s\n", aaaa, bbbb)
    case 15:
        fmt.Printf("IFB %s %s\n", aaaa, bbbb)
    }
}

func main() {
    if contents, err := ioutil.ReadFile(os.Args[1]); err == nil {
        fileLength := len(contents)

        for i := 0; i < fileLength; i += 2 {
            program = append(program, (Word(contents[i]) << 8) + Word(contents[i + 1]))
        }

        for int(PC) < len(program) {
            readInstruction(program[PC])
        }
    }
}
