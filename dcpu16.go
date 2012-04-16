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

    Memory []Word
)

func (cpu *Dcpu) dumpRegisters() {
    fmt.Printf("A: %d B: %x C: %d X: %d Y: %d Z: %d I: %d J: %d PC: %d SP: %d\n", cpu.A, cpu.B, cpu.C, cpu.X, cpu.Y, cpu.Z, cpu.I, cpu.J, cpu.PC - 1, cpu.SP)
    /*
    Memory[0x8000] = 65 + (0x35 << 8)
    Memory[0x8001] = 58
    Memory[0x8002] = 32
    a := []byte(fmt.Sprintf("%d", cpu.A))
    Memory[0x8003] = Word(a[0])
    Memory[0x8004] = 32

    Memory[0x8005] = 88 + (0x42 << 8)
    Memory[0x8006] = 58
    Memory[0x8007] = 32
    x := []byte(fmt.Sprintf("%d", cpu.X))
    Memory[0x8008] = Word(x[0])
    Memory[0x8009] = 32
    */
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

func (cpu *Dcpu) dumpVideoRam() {
    for i := 0x8000; i < 0x8200; i++ {
        if Memory[i] > 0 {
            fmt.Printf("D: %d Register I: %d\n", Memory[i] & 0x7f, cpu.I)
        }
    }
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
    // cpu.dumpVideoRam()
}
