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

	assert.Equal(t, int32(42), rs.Read(1), "immediate value should be loaded into destination register")
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

	assert.Equal(t, int32(15), rs.Read(1), "immediate should be added to existing register value")
}

func TestAdd(t *testing.T) {
	cases := []struct {
		name     string
		a        int
		b        int
		expected int32
		message  string
	}{
		{"basic", 10, 5, 15, "should add two positive numbers"},
		{"overflow", 2147483647, 1, -2147483648, "should wrap around on overflow"},
		{"negatives", -5, 3, -2, "should handle negative numbers"},
		{"identity", 42, 0, 42, "adding zero should return same value"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rs := registers.NewRegisters()
			mem := memory.NewMemory(1024)
			vm := NewVM(rs, mem)

			bytecode := ByteCode{
				int(opcodes.ADDI), 1, 0, tc.a,
				int(opcodes.ADDI), 2, 0, tc.b,
				int(opcodes.ADD), 3, 1, 2,
			}

			vm.Execute(bytecode)

			assert.Equal(t, tc.expected, rs.Read(3), tc.message)
		})
	}
}

func TestX0AsSourceReadsZero(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.ADD), 2, 1, 0,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(42), rs.Read(2), "x0 should always read as zero")
}

func TestSub(t *testing.T) {
	cases := []struct {
		name     string
		a        int
		b        int
		expected int32
		message  string
	}{
		{"basic", 10, 3, 7, "should subtract second from first"},
		{"underflow", -2147483648, 1, 2147483647, "should wrap around on underflow"},
		{"zero result", 42, 42, 0, "subtracting equal values should give zero"},
		{"negative result", 5, 10, -5, "should produce negative result when b > a"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rs := registers.NewRegisters()
			mem := memory.NewMemory(1024)
			vm := NewVM(rs, mem)

			bytecode := ByteCode{
				int(opcodes.ADDI), 1, 0, tc.a,
				int(opcodes.ADDI), 2, 0, tc.b,
				int(opcodes.SUB), 3, 1, 2,
			}

			vm.Execute(bytecode)

			assert.Equal(t, tc.expected, rs.Read(3), tc.message)
		})
	}
}

func TestMul(t *testing.T) {
	cases := []struct {
		name     string
		a        int
		b        int
		expected int32
		message  string
	}{
		{"basic", 6, 7, 42, "should multiply two numbers"},
		{"by zero", 42, 0, 0, "multiplying by zero should give zero"},
		{"by one", 42, 1, 42, "multiplying by one should return same value"},
		{"negatives", -3, 4, -12, "should handle negative numbers"},
		{"overflow", 2147483647, 2, -2, "should wrap around on overflow"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rs := registers.NewRegisters()
			mem := memory.NewMemory(1024)
			vm := NewVM(rs, mem)

			bytecode := ByteCode{
				int(opcodes.ADDI), 1, 0, tc.a,
				int(opcodes.ADDI), 2, 0, tc.b,
				int(opcodes.MUL), 3, 1, 2,
			}

			vm.Execute(bytecode)

			assert.Equal(t, tc.expected, rs.Read(3), tc.message)
		})
	}
}

func TestDiv(t *testing.T) {
	cases := []struct {
		name     string
		a        int
		b        int
		expected int32
		message  string
	}{
		{"basic", 20, 4, 5, "should divide first by second"},
		{"truncation", 7, 3, 2, "should truncate towards zero"},
		{"by one", 42, 1, 42, "dividing by one should return same value"},
		{"zero dividend", 0, 5, 0, "zero divided by anything should give zero"},
		{"negatives", -20, 4, -5, "should handle negative dividend"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rs := registers.NewRegisters()
			mem := memory.NewMemory(1024)
			vm := NewVM(rs, mem)

			bytecode := ByteCode{
				int(opcodes.ADDI), 1, 0, tc.a,
				int(opcodes.ADDI), 2, 0, tc.b,
				int(opcodes.DIV), 3, 1, 2,
			}

			vm.Execute(bytecode)

			assert.Equal(t, tc.expected, rs.Read(3), tc.message)
		})
	}
}

func TestMod(t *testing.T) {
	cases := []struct {
		name     string
		a        int
		b        int
		expected int32
		message  string
	}{
		{"basic", 17, 5, 2, "should return remainder after division"},
		{"negatives", -7, 3, -1, "remainder should have same sign as dividend"},
		{"zero dividend", 0, 5, 0, "zero mod anything should give zero"},
		{"divisor larger", 3, 7, 3, "when divisor is larger, result is dividend"},
		{"equal operands", 5, 5, 0, "equal values should give zero remainder"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rs := registers.NewRegisters()
			mem := memory.NewMemory(1024)
			vm := NewVM(rs, mem)

			bytecode := ByteCode{
				int(opcodes.ADDI), 1, 0, tc.a,
				int(opcodes.ADDI), 2, 0, tc.b,
				int(opcodes.MOD), 3, 1, 2,
			}

			vm.Execute(bytecode)

			assert.Equal(t, tc.expected, rs.Read(3), tc.message)
		})
	}
}

func TestAddSubInversion(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.ADDI), 2, 0, 17,
		int(opcodes.ADD), 3, 1, 2,
		int(opcodes.SUB), 4, 3, 2,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(42), rs.Read(1), "original value should be unchanged")
	assert.Equal(t, int32(42), rs.Read(4), "add then sub should return original value")
}

func TestDivModRelationship(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 17,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.DIV), 3, 1, 2,
		int(opcodes.MOD), 4, 1, 2,
		int(opcodes.MUL), 5, 3, 2,
		int(opcodes.ADD), 6, 5, 4,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(17), rs.Read(1), "original dividend should be unchanged")
	assert.Equal(t, int32(17), rs.Read(6), "(a/b)*b + (a%b) should equal a")
}

func TestMulDivInversion(t *testing.T) {
	rs := registers.NewRegisters()
	mem := memory.NewMemory(1024)
	vm := NewVM(rs, mem)

	bytecode := ByteCode{
		int(opcodes.ADDI), 1, 0, 7,
		int(opcodes.ADDI), 2, 0, 3,
		int(opcodes.MUL), 3, 1, 2,
		int(opcodes.DIV), 4, 3, 2,
	}

	vm.Execute(bytecode)

	assert.Equal(t, int32(7), rs.Read(1), "original value should be unchanged")
	assert.Equal(t, int32(7), rs.Read(4), "mul then div should return original value")
}
