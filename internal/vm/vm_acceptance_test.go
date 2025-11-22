package vm

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/memory"
	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/phasecurve/zhuji/internal/registers"
	"github.com/stretchr/testify/assert"
)

func TestSumOneToFive(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 0,
		int(opcodes.ADDI), 2, 0, 1,
		int(opcodes.ADDI), 3, 0, 6,

		int(opcodes.ADD), 1, 1, 2,
		int(opcodes.ADDI), 2, 2, 1,
		int(opcodes.BLT), 2, 3, -8,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(15), rs.Read(1), "sum of 1+2+3+4+5 should be 15")
	assert.Equal(t, int32(6), rs.Read(2), "counter should stop at limit")
}

func TestFibonacci(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 0, // set x1=0
		int(opcodes.ADDI), 2, 0, 1, // set x2=1
		int(opcodes.ADDI), 3, 0, 9, // set x3=9 - target to reach to branch to the exit
		int(opcodes.ADDI), 4, 0, 0, // set x4=0 the counter to increment for each recursion

		int(opcodes.ADD), 5, 1, 2, // set x5=x1+x2
		int(opcodes.ADD), 1, 2, 0, // set x1=x2+x0 or x1=x2+0 since x0 is hardcoded to 0 in risc-v
		int(opcodes.ADD), 2, 5, 0, // set x2=x5+x0
		int(opcodes.ADDI), 4, 4, 1, // set x4=x4+1 - increment the counter
		int(opcodes.BLT), 4, 3, -16, // check if x4<x3 - has the counter reached 9 then don't jump, continue to exit, otherwise repeat
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(55), rs.Read(2), "10th fibonacci number should be 55")
}
