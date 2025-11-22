package assembler

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/stretchr/testify/assert"
)

func TestAssembleArithmeticInstructions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []int
		message  string
	}{
		{"add", "add x3, x1, x2", []int{int(opcodes.ADD), 3, 1, 2}, "add should encode destination and two source registers"},
		{"sub", "sub x3, x1, x2", []int{int(opcodes.SUB), 3, 1, 2}, "sub should encode destination and two source registers"},
		{"mul", "mul x3, x1, x2", []int{int(opcodes.MUL), 3, 1, 2}, "mul should encode destination and two source registers"},
		{"div", "div x3, x1, x2", []int{int(opcodes.DIV), 3, 1, 2}, "div should encode destination and two source registers"},
		{"mod", "mod x3, x1, x2", []int{int(opcodes.MOD), 3, 1, 2}, "mod should encode destination and two source registers"},
		{"addi", "addi x1, x0, 42", []int{int(opcodes.ADDI), 1, 0, 42}, "addi should encode destination, source register, and immediate value"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			asm := NewAssembler()
			bytecode := asm.Assemble(tc.input)
			assert.Equal(t, tc.expected, bytecode, tc.message)
		})
	}
}

func TestAssembleBranchInstructions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []int
		message  string
	}{
		{"blt", "blt x1, x2, 12", []int{int(opcodes.BLT), 1, 2, 12}, "blt should encode two registers and offset for branch if less than"},
		{"beq", "beq x1, x2, 12", []int{int(opcodes.BEQ), 1, 2, 12}, "beq should encode two registers and offset for branch if equal"},
		{"bne", "bne x1, x2, 12", []int{int(opcodes.BNE), 1, 2, 12}, "bne should encode two registers and offset for branch if not equal"},
		{"bge", "bge x1, x2, 12", []int{int(opcodes.BGE), 1, 2, 12}, "bge should encode two registers and offset for branch if greater or equal"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			asm := NewAssembler()
			bytecode := asm.Assemble(tc.input)
			assert.Equal(t, tc.expected, bytecode, tc.message)
		})
	}
}

func TestAssembleMemoryInstructions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []int
		message  string
	}{
		{"lw", "lw x1, 0(x0)", []int{int(opcodes.LW), 1, 0, 0}, "lw should encode destination, offset, and base register for load word"},
		{"sw", "sw x1, 4(x2)", []int{int(opcodes.SW), 1, 4, 2}, "sw should encode source, offset, and base register for store word"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			asm := NewAssembler()
			bytecode := asm.Assemble(tc.input)
			assert.Equal(t, tc.expected, bytecode, tc.message)
		})
	}
}
