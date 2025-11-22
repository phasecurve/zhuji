package assembler

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, expected, bytecode, "sum program should assemble all instructions in sequence")
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
	assert.Equal(t, expected, bytecode, "fibonacci program should assemble with correct branch offsets")
}
