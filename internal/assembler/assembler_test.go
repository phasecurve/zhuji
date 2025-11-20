package assembler

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/stretchr/testify/assert"
)

func TestAssembleAddi(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("addi x1, x0, 42")

	expected := []int{
		int(opcodes.ADDI), 1, 0, 42,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleAdd(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("add x3, x1, x2")

	expected := []int{
		int(opcodes.ADD), 3, 1, 2,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleSub(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("sub x3, x1, x2")

	expected := []int{
		int(opcodes.SUB), 3, 1, 2,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleLw(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("lw x1, 0(x0)")

	expected := []int{
		int(opcodes.LW), 1, 0, 0,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleSw(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("sw x1, 4(x2)")

	expected := []int{
		int(opcodes.SW), 1, 4, 2,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleBlt(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("blt x1, x2, 12")

	expected := []int{
		int(opcodes.BLT), 1, 2, 12,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleSumProgram(t *testing.T) {
	asm := NewAssembler()

	program := `addi x1, x0, 0
addi x2, x0, 1
addi x3, x0, 6
add x1, x1, x2
addi x2, x2, 1
blt x2, x3, -8`

	bytecode := asm.Assemble(program)

	expected := []int{
		int(opcodes.ADDI), 1, 0, 0,
		int(opcodes.ADDI), 2, 0, 1,
		int(opcodes.ADDI), 3, 0, 6,
		int(opcodes.ADD), 1, 1, 2,
		int(opcodes.ADDI), 2, 2, 1,
		int(opcodes.BLT), 2, 3, -8,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleBeq(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("beq x1, x2, 12")

	expected := []int{
		int(opcodes.BEQ), 1, 2, 12,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleBne(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("bne x1, x2, 12")

	expected := []int{
		int(opcodes.BNE), 1, 2, 12,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleBge(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("bge x1, x2, 12")

	expected := []int{
		int(opcodes.BGE), 1, 2, 12,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleMul(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("mul x3, x1, x2")

	expected := []int{
		int(opcodes.MUL), 3, 1, 2,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleDiv(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("div x3, x1, x2")

	expected := []int{
		int(opcodes.DIV), 3, 1, 2,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleMod(t *testing.T) {
	asm := NewAssembler()

	bytecode := asm.Assemble("mod x3, x1, x2")

	expected := []int{
		int(opcodes.MOD), 3, 1, 2,
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleFibonacci(t *testing.T) {
	asm := NewAssembler()

	program := `addi x1, x0, 0
addi x2, x0, 1
addi x3, x0, 9
addi x4, x0, 0
bge x4, x3, 40
add x5, x1, x2
add x1, x2, x0
add x2, x5, x0
addi x4, x4, 1
blt x4, x3, 16`

	bytecode := asm.Assemble(program)

	expected := []int{
		int(opcodes.ADDI), 1, 0, 0,
		int(opcodes.ADDI), 2, 0, 1,
		int(opcodes.ADDI), 3, 0, 9,
		int(opcodes.ADDI), 4, 0, 0,
		int(opcodes.BGE), 4, 3, 40,
		int(opcodes.ADD), 5, 1, 2,
		int(opcodes.ADD), 1, 2, 0,
		int(opcodes.ADD), 2, 5, 0,
		int(opcodes.ADDI), 4, 4, 1,
		int(opcodes.BLT), 4, 3, 16,
	}
	assert.Equal(t, expected, bytecode)
}
