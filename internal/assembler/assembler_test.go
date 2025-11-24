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

func TestAssembleLabelDefinition(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("loop:\naddi x1, x0, 1")

	if len(bytecode) != 4 {
		t.Errorf("expected 4 bytes (just addi), got %d", len(bytecode))
	}
}

func TestAssembleBranchWithLabel(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("loop:\naddi x1, x0, 1\nblt x1, x2, loop")

	expected := []int{
		int(opcodes.ADDI), 1, 0, 1,
		int(opcodes.BLT), 1, 2, -4,
	}
	assert.Equal(t, expected, bytecode, "branch should resolve label to PC-relative offset")
}

func TestAssembleBranchWithLabelNotAtPositionFour(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("addi x1, x0, 1\nloop:\naddi x2, x0, 2\nblt x1, x2, loop")

	expected := []int{
		int(opcodes.ADDI), 1, 0, 1,
		int(opcodes.ADDI), 2, 0, 2,
		int(opcodes.BLT), 1, 2, -4,
	}
	assert.Equal(t, expected, bytecode, "branch at position 8 should jump back to position 4")
}

func TestAssemblePseudoInstructions(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []int
		message  string
	}{
		{"li", "li x1, 42", []int{int(opcodes.ADDI), 1, 0, 42}, "li should expand to addi rd, x0, imm"},
		{"mv", "mv x1, x2", []int{int(opcodes.ADDI), 1, 2, 0}, "mv should expand to addi rd, rs, 0"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			asm := NewAssembler()
			bytecode := asm.Assemble(tc.input)
			assert.Equal(t, tc.expected, bytecode, tc.message)
		})
	}
}

func TestAssembleJALWithLabel(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("addi x2, x0, 3\naddi x1, x0, 1\nloop:\nadd x1, x1, x1\nblt x2, x1, 8\njal x0, loop")

	expected := []int{
		int(opcodes.ADDI), 2, 0, 3,
		int(opcodes.ADDI), 1, 0, 1,
		int(opcodes.ADD), 1, 1, 1,
		int(opcodes.BLT), 2, 1, 8,
		int(opcodes.JAL), 0, 0, -8,
	}
	assert.Equal(t, expected, bytecode, "branch should resolve label to PC-relative offset")
}
