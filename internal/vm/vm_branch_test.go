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
