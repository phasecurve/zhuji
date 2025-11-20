package vm

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/memory"
	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/phasecurve/zhuji/internal/registers"
	"github.com/stretchr/testify/assert"
)

func TestAddiLoadsImmediateIntoRegister(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(42), rs.Read(1))
}

func TestAddiAddsToNonZeroRegister(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 1, 1, 5,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(15), rs.Read(1))
}

func TestAddRegisters(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.ADD), 3, 1, 2,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(10), rs.Read(1))
	assert.Equal(t, int32(5), rs.Read(2))
	assert.Equal(t, int32(15), rs.Read(3))
}

func TestSubRegisters(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 3,
		int(opcodes.SUB), 3, 1, 2,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(10), rs.Read(1))
	assert.Equal(t, int32(3), rs.Read(2))
	assert.Equal(t, int32(7), rs.Read(3))
}

func TestMulRegisters(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 6,
		int(opcodes.ADDI), 2, 0, 7,
		int(opcodes.MUL), 3, 1, 2,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(6), rs.Read(1))
	assert.Equal(t, int32(7), rs.Read(2))
	assert.Equal(t, int32(42), rs.Read(3))
}

func TestDivRegisters(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 20,
		int(opcodes.ADDI), 2, 0, 4,
		int(opcodes.DIV), 3, 1, 2,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(20), rs.Read(1))
	assert.Equal(t, int32(4), rs.Read(2))
	assert.Equal(t, int32(5), rs.Read(3))
}

func TestModRegisters(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 17,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.MOD), 3, 1, 2,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(17), rs.Read(1))
	assert.Equal(t, int32(5), rs.Read(2))
	assert.Equal(t, int32(2), rs.Read(3))
}

func TestStoreWord(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.SW), 1, 0, 0,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(42), mem.LoadWord(0))
}

func TestLoadWord(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	mem.StoreWord(0, 99)

	bytecode := ByteCode{
		int(opcodes.LW), 1, 0, 0,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(99), rs.Read(1))
}

func TestBranchIfLessThanTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BLT), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
		int(opcodes.ADDI), 4, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(5), rs.Read(1))
	assert.Equal(t, int32(10), rs.Read(2))
	assert.Equal(t, int32(0), rs.Read(3))
	assert.Equal(t, int32(42), rs.Read(4))
}

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
		int(opcodes.BLT), 2, 3, 12,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(15), rs.Read(1))
	assert.Equal(t, int32(6), rs.Read(2))
}

func TestBranchIfEqualTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BEQ), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
		int(opcodes.ADDI), 4, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(10), rs.Read(1))
	assert.Equal(t, int32(10), rs.Read(2))
	assert.Equal(t, int32(0), rs.Read(3))
	assert.Equal(t, int32(42), rs.Read(4))
}

func TestBranchIfNotEqualTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BNE), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
		int(opcodes.ADDI), 4, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(5), rs.Read(1))
	assert.Equal(t, int32(10), rs.Read(2))
	assert.Equal(t, int32(0), rs.Read(3))
	assert.Equal(t, int32(42), rs.Read(4))
}

func TestBranchIfGreaterOrEqualTaken(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BGE), 1, 2, 16,
		int(opcodes.ADDI), 3, 0, 99,
		int(opcodes.ADDI), 4, 0, 42,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(10), rs.Read(1))
	assert.Equal(t, int32(5), rs.Read(2))
	assert.Equal(t, int32(0), rs.Read(3))
	assert.Equal(t, int32(42), rs.Read(4))
}

func TestFibonacci(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
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

	vm.Execute(bytecode)

	assert.Equal(t, int32(55), rs.Read(2))
}
