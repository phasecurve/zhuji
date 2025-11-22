package vm

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/memory"
	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/phasecurve/zhuji/internal/registers"
	"github.com/stretchr/testify/assert"
)

func TestStoreWord(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.SW), 1, 0, 0,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(42), mem.LoadWord(0), "register value should be stored to memory")
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

	assert.Equal(t, int32(99), rs.Read(1), "memory value should be loaded into register")
}

func TestLoadStoreConsecutiveAddresses(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.ADDI), 2, 0, 99,
		int(opcodes.SW), 1, 0, 0,
		int(opcodes.SW), 2, 4, 0,
		int(opcodes.LW), 3, 0, 0,
		int(opcodes.LW), 4, 4, 0,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(42), rs.Read(3), "value at address 0 should be isolated from address 4")
	assert.Equal(t, int32(99), rs.Read(4), "value at address 4 should be isolated from address 0")
}

func TestLoadStoreOverwrite(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.SW), 1, 0, 0,
		int(opcodes.ADDI), 1, 0, 99,
		int(opcodes.SW), 1, 0, 0,
		int(opcodes.LW), 2, 0, 0,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(99), rs.Read(2), "second store should overwrite first value")
}
