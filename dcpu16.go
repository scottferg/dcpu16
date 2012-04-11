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
}

var (
    cpu Dcpu

    Memory [0xffff]Word
)

/********************************************************************************************************************************************

0x0000 - 0x7FFD - program space. Can be used as pleased by DCPU software.
0x7FFFE         - PC Speaker word. Least significant bit acts as flag of whether or not PC Speaker is playing, rest set frequency.
0x7FFF          - screen mode word. 0x6FC means GraPhiCs mode, any other value means screen is in standard 32x16 text mode.
0x8000 - 0x8200 - text buffer. Can be used as pleased if 0x7FFF == 0x6FC (graphics mode).
0x8201          - FDD chunk ID. Each FDD chunk has 3582(dec) words of data and there is 211(dec) chunks (211th not full) that form 1.44MB of floppy
0x8202 - 0x8FFF - FDD data
0x9000          - keyboard word. Preserves last pressed key.
0x9001 - 0x94FF - preserved for COMM peripherals
0x9500 - 0xAD00 - VRAM. Used when screen mode word is set to 0x6FC (Graphics mode). Graphics resolution is 192x128 and each pixel is represented by 5 bits. Palette is fixed and has 16 colors (where 0x0 is black and 0xF is white). Can be used as pleased if screen mode flag (0x7FFF) is set to any different value than 0x6FC.
0xAD01          - D411 preserved for NAVI peripherals.
0xD412 - 0xE79A - preserved for COMM peripherals
0xE79B - 0xFFFF - Stack

******************************************************************************************************************************************/

func (cpu *Dcpu) dumpRegisters() {
    fmt.Printf("A: %d B: %d C: %d X: %d Y: %d Z: %d I: %d J: %d PC: %d\n", cpu.A, cpu.B, cpu.C, cpu.X, cpu.Y, cpu.Z, cpu.I, cpu.J, cpu.PC - 1)
}

func (cpu *Dcpu) WordCount(opcode Word) Word {
    a := (opcode >> 4) & 0x3f
    b := (opcode >> 10) & 0x3f

    a, _ = cpu.processOperand(a)
    b, _ = cpu.processOperand(b)

	count := Word(1)

	switch {
	case a >= 16 && a <= 23:
	case a == 30:
	case a == 31:
		count++
	}

	switch {
	case b >= 16 && b <= 23:
	case b == 30:
	case b == 31:
		count++
	}

	return count
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
        assignee = &Memory[Memory[cpu.PC] + cpu.A]
        cpu.PC++
    case 17:
        assignee = &Memory[Memory[cpu.PC] + cpu.B]
        cpu.PC++
    case 18:
        assignee = &Memory[Memory[cpu.PC] + cpu.C]
        cpu.PC++
    case 19:
        assignee = &Memory[Memory[cpu.PC] + cpu.X]
        cpu.PC++
    case 20:
        assignee = &Memory[Memory[cpu.PC] + cpu.Y]
        cpu.PC++
    case 21:
        assignee = &Memory[Memory[cpu.PC] + cpu.Z]
        cpu.PC++
    case 22:
        assignee = &Memory[Memory[cpu.PC] + cpu.I]
        cpu.PC++
    case 23:
        assignee = &Memory[Memory[cpu.PC] + cpu.J]
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
        value = operand - 32
    }

    if assignee != nil {
        value = *assignee
    }

    return
}

func (cpu *Dcpu) Step() {
    opcode := Memory[cpu.PC]
    cpu.PC++

    instruction := opcode & 0xf
    a := (opcode >> 4) & 0x3f
    b := (opcode >> 10) & 0x3f

    var assignable *Word

    if instruction != 0 {
        a, assignable = cpu.processOperand(a)
        b, _ = cpu.processOperand(b)
    }

    switch instruction {
    case 0:
        instruction, a = a, b

        switch instruction {
        case 1:
            _, assignable = cpu.processOperand(0x1a)
            a, _ = cpu.processOperand(b)

            *assignable = cpu.PC
            cpu.PC = a
        }
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
        if a != b {
            cpu.PC += cpu.WordCount(Memory[cpu.PC])
        }
    case 13:
        // IFN a, b
        if a == b {
            cpu.PC += cpu.WordCount(Memory[cpu.PC])
        }
    case 14:
        // IFG a, b
        if !(a > b) {
            cpu.PC += cpu.WordCount(Memory[cpu.PC])
        }
    case 15:
        // IFB a, b
        if a & b == 0 {
            cpu.PC += cpu.WordCount(Memory[cpu.PC])
        }
    }

    cpu.dumpRegisters()
}

func main() {
    program := []Word{
        0x7c01, 0x0030, 0x7de1, 0x1000, 0x0020, 0x7803, 0x1000, 0xc00d,
        0x7dc1, 0x001a, 0xa861, 0x7c01, 0x2000, 0x2161, 0x2000, 0x8463,
        0x806d, 0x7dc1, 0x000d, 0x9031, 0x7c10, 0x0018, 0x7dc1, 0x001a,
        0x9037, 0x61c1, 0x7dc1, 0x001a, 0x0000, 0x0000, 0x0000, 0x0000,
    }

    cpu := new(Dcpu)
    cpu.PC = 0
    cpu.SP = 0xffff

    for index, value := range program {
        Memory[index] = value
    }

    for int(cpu.PC) < len(program) {
        cpu.Step()
    }
}
