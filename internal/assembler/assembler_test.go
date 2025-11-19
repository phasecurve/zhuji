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
blt x2, x3, 12`

	bytecode := asm.Assemble(program)

	expected := []int{
		int(opcodes.ADDI), 1, 0, 0,
		int(opcodes.ADDI), 2, 0, 1,
		int(opcodes.ADDI), 3, 0, 6,
		int(opcodes.ADD), 1, 1, 2,
		int(opcodes.ADDI), 2, 2, 1,
		int(opcodes.BLT), 2, 3, 12,
	}
	assert.Equal(t, expected, bytecode)
}
