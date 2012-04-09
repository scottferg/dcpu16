package main

import (
    "fmt"
)

type Word uint16

type Dcpu struct {
    A    Word
    B    Word
    C    Word
    X    Word
    Y    Word
    Z    Word
    I    Word
    J    Word
    PC   Word
    SP   Word
    O    Word
    Skip bool
}

var (
    cpu Dcpu

    Memory [0x10000]Word
)

func (cpu *Dcpu) dumpRegisters() {
    fmt.Println("  A  |  B  |  C  |  X  |  Y  |  Z  |  I  |  J  |  PC  ")
    fmt.Printf("  %d  |  %d  |  %d  |  %d  |  %d  |  %d  |  %d  |  %d  |  %d  \n", cpu.A, cpu.B, cpu.C, cpu.X, cpu.Y, cpu.Z, cpu.I, cpu.J, cpu.PC)
}

func (cpu *Dcpu) processOperand(operand Word) (value Word, assignee *Word) {
    switch operand {
    case 0:
        assignee = &cpu.A
    case 1:
        assignee = &cpu.B
    case 2:
        assignee = &cpu.C
    case 3:
        assignee = &cpu.X
    case 4:
        assignee = &cpu.Y
    case 5:
        assignee = &cpu.Z
    case 6:
        assignee = &cpu.I
    case 7:
        assignee = &cpu.J
    case 8:
        value = cpu.A
    case 9:
        value = cpu.B
    case 10:
        value = cpu.C
    case 11:
        value = cpu.X
    case 12:
        value = cpu.Y
    case 13:
        value = cpu.Z
    case 14:
        value = cpu.I
    case 15:
        value = cpu.J
    case 16:
        assignee = &Memory[cpu.PC + cpu.A]
        cpu.PC++
    case 17:
        assignee = &Memory[cpu.PC + cpu.B]
        cpu.PC++
    case 18:
        assignee = &Memory[cpu.PC + cpu.C]
        cpu.PC++
    case 19:
        assignee = &Memory[cpu.PC + cpu.X]
        cpu.PC++
    case 20:
        assignee = &Memory[cpu.PC + cpu.Y]
        cpu.PC++
    case 21:
        assignee = &Memory[cpu.PC + cpu.Z]
        cpu.PC++
    case 22:
        assignee = &Memory[cpu.PC + cpu.I]
        cpu.PC++
    case 23:
        assignee = &Memory[cpu.PC + cpu.J]
        cpu.PC++
    case 24:
        assignee = &Memory[cpu.SP]
        cpu.SP++
    case 25:
        assignee = &Memory[cpu.SP]
    case 26:
        cpu.SP--
        assignee = &Memory[cpu.SP]
    case 27:
        assignee = &cpu.SP
    case 28:
        assignee = &cpu.PC
    case 29:
        assignee = &cpu.O
    case 30:
        assignee = &Memory[Memory[cpu.PC]]
        cpu.PC++
    case 31:
        value = Memory[cpu.PC]
        cpu.PC++
    default:
        value = operand
    }

    if (assignee != nil) {
        value = *assignee
    }

    return
}

func (cpu *Dcpu) Step() {
    opcode := Memory[cpu.PC]
    cpu.PC++

    if (cpu.Skip) {
        cpu.PC++
        cpu.Skip = false
    }

    instruction := opcode & 0xf
    a := (opcode >> 4) & 0x3f
    b := (opcode >> 10) & 0x3f

    var assignable *Word

    a, assignable = cpu.processOperand(a)
    b, _ = cpu.processOperand(b)

    switch instruction {
    case 0:
        // Do nothing yet
    case 1:
        // SET a, b
        *assignable = b
    case 2:
        // ADD a, b
        *assignable = a + b
    case 3:
        // SUB a, b
        *assignable = a - b
    case 4:
        // MUL a, b
        *assignable = a * b
    case 5:
        // DIV a, b
        *assignable = a / b
    case 6:
        // MOD a, b
        *assignable = a % b
    case 7:
        // SHL a, b
        *assignable = a << b
    case 8:
        // SHR a, b
        *assignable = a >> b
    case 9:
        // AND a, b
        *assignable = a & b
    case 10:
        // BOR a, b
        *assignable = a | b
    case 11:
        // XOR a, b
        *assignable = a ^ b
    case 12:
        // IFE a, b
        if (a == b) {
            cpu.Skip = true
        }
    case 13:
        // IFN a, b
        if (a != b) {
            cpu.Skip = true
        }
    case 14:
        // IFG a, b
        if (a > b) {
            cpu.Skip = true
        }
    case 15:
        // IFB a, b
        if (a & b != 0) {
            cpu.Skip = true
        }
    }

    cpu.dumpRegisters()
}

func main() {
    program := []Word{
        0x7c01, 0x0030, 0x7de1, 0x1000, 0x0020, 0x7803, 0x1000,
    }

    cpu := new(Dcpu)
    cpu.PC = 0

    for index, value := range program {
        Memory[index] = value
    }

    for ; int(cpu.PC) < len(program); {
        cpu.Step()
    }
}
