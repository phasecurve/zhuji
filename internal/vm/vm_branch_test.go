package vm

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/memory"
	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/phasecurve/zhuji/internal/registers"
	"github.com/stretchr/testify/assert"
)

func TestBranchIfLessThanTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BLT), 1, 2, 8,
		int(opcodes.ADDI), 3, 0, 99,
		int(opcodes.ADDI), 4, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(0), rs.Read(3), "skipped instruction should not execute when branch taken")
	assert.Equal(t, int32(42), rs.Read(4), "instruction after skip should execute")
}

func TestBranchIfLessThanNotTakenWhenEqual(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BLT), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(99), rs.Read(3), "next instruction should execute when values are equal")
}

func TestBranchIfLessThanWithNegatives(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, -5,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BLT), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(0), rs.Read(3), "branch should be taken when negative is less than positive")
}

func TestBranchIfEqualTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BEQ), 1, 2, 8,
		int(opcodes.ADDI), 3, 0, 99,
		int(opcodes.ADDI), 4, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(0), rs.Read(3), "skipped instruction should not execute when values are equal")
	assert.Equal(t, int32(42), rs.Read(4), "instruction after skip should execute")
}

func TestBranchIfNotEqualTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BNE), 1, 2, 8,
		int(opcodes.ADDI), 3, 0, 99,
		int(opcodes.ADDI), 4, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(0), rs.Read(3), "skipped instruction should not execute when values differ")
	assert.Equal(t, int32(42), rs.Read(4), "instruction after skip should execute")
}

func TestBranchIfGreaterOrEqualTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BGE), 1, 2, 8,
		int(opcodes.ADDI), 3, 0, 99,
		int(opcodes.ADDI), 4, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(0), rs.Read(3), "skipped instruction should not execute when first >= second")
	assert.Equal(t, int32(42), rs.Read(4), "instruction after skip should execute")
}

func TestBranchIfGreaterOrEqualTakenWhenEqual(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BGE), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(0), rs.Read(3), "branch should be taken when values are equal")
}

func TestBranchIfGreaterOrEqualWithNegatives(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, -5,
		int(opcodes.BGE), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(0), rs.Read(3), "branch should be taken when positive >= negative")
}

func TestBranchIfEqualNotTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BEQ), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(99), rs.Read(3), "next instruction should execute when values differ")
}

func TestBranchIfNotEqualNotTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BNE), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(99), rs.Read(3), "next instruction should execute when values are equal")
}

func TestJALSavesReturnAddressAndJumps(t *testing.T) {
	/*	sampleAsm:
		jal x1, 8      # IP=0: save 4 into x1, jump to 8
		addi x2, x0, 1 # IP=4: skipped
		addi x3, x0, 2 # IP=8: executed
	*/
	expectedNextAddr := int32(4)
	expectedValInX2 := int32(0)
	expectedValInX3 := int32(2)

	bytecode := ByteCode{
		int(opcodes.JAL), 1, 0, 8,
		int(opcodes.ADDI), 2, 0, 1,
		int(opcodes.ADDI), 3, 0, 2,
	}

	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)
	vm.Execute(bytecode)

	actualNextAddr := rs.Read(1)
	actualValInX2 := rs.Read(2)
	actualValInX3 := rs.Read(3)

	assert.Equal(t, expectedNextAddr, actualNextAddr, "the next address should be 4 (from x1)")
	assert.Equal(t, expectedValInX3, actualValInX3, "the val should be 2 in x3")
	assert.Equal(t, expectedValInX2, actualValInX2, "the val should be 0 in x2")
}

func TestJALRJumpsToOffestPlusAddr(t *testing.T) {
	expectedVal := int32(1)

	bytecode := ByteCode{
		int(opcodes.ADDI), 5, 0, 1,
		int(opcodes.JAL), 1, 0, 8,
		int(opcodes.ADDI), 2, 0, 1,
		int(opcodes.BEQ), 5, 2, 12,
		int(opcodes.ADDI), 3, 0, 2,
		int(opcodes.JALR), 0, 1, 0,
	}

	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)
	vm.Execute(bytecode)

	actualVal := rs.Read(2)

	assert.Equal(t, expectedVal, actualVal, "should have value 1 after addi sends 0+1 to x2")
}

func TestJALRJumpsToOffestPlusAddrAndSavesNextAddrToX6(t *testing.T) {
	expectedVal := int32(24)

	bytecode := ByteCode{
		int(opcodes.ADDI), 5, 0, 1,
		int(opcodes.JAL), 1, 0, 8,
		int(opcodes.ADDI), 2, 0, 1,
		int(opcodes.BEQ), 5, 2, 12,
		int(opcodes.ADDI), 3, 0, 2,
		int(opcodes.JALR), 6, 1, 0,
	}

	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)
	vm.Execute(bytecode)

	actualVal := rs.Read(6)

	assert.Equal(t, expectedVal, actualVal, "x6 should have jump ip=24 after JALR executes")
}
