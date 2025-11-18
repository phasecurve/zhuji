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
