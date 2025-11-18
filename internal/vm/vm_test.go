package vm

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/phasecurve/zhuji/internal/registers"
	"github.com/phasecurve/zhuji/internal/stack"
	"github.com/stretchr/testify/assert"
)

func TestExecuteSinglePushInstruction(t *testing.T) {
	st := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(st, rs)
	pushOp := []int{int(opcodes.PSH), 42}

	vm.Execute(pushOp)

	assert.Equal(t, st.Peek(), 42, "should have executed a push instruction")
}

func TestExecuteMultiplePushInstructions(t *testing.T) {
	st := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(st, rs)
	bytecode := []int{int(opcodes.PSH), 10, int(opcodes.PSH), 20}

	vm.Execute(bytecode)

	assert.Equal(t, 20, st.Pop())
	assert.Equal(t, 10, st.Pop())
}

func TestExecuteAddInstruction(t *testing.T) {
	st := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(st, rs)
	// PSH 10, PSH 20, ADD
	bytecode := []int{int(opcodes.PSH), 10, int(opcodes.PSH), 20, int(opcodes.ADD)}

	vm.Execute(bytecode)

	assert.Equal(t, 30, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteSubInstruction(t *testing.T) {
	st := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(st, rs)
	// PSH 20, PSH 10, SUB → (20 - 10 = 10)
	bytecode := []int{int(opcodes.PSH), 20, int(opcodes.PSH), 10, int(opcodes.SUB)}

	vm.Execute(bytecode)

	assert.Equal(t, 10, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteDivInstruction(t *testing.T) {
	st := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(st, rs)
	// PSH 20, PSH 4, DIV → (20 / 4 = 5)
	bytecode := []int{int(opcodes.PSH), 20, int(opcodes.PSH), 4, int(opcodes.DIV)}

	vm.Execute(bytecode)

	assert.Equal(t, 5, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteMulInstruction(t *testing.T) {
	st := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(st, rs)
	// PSH 20, PSH 4, MUL → (20 * 4 = 80)
	bytecode := []int{int(opcodes.PSH), 20, int(opcodes.PSH), 4, int(opcodes.MUL)}

	vm.Execute(bytecode)

	assert.Equal(t, 80, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteComplexExpression(t *testing.T) {
	st := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(st, rs)
	// (5 + 3) * 2 = 16
	// RPN: 5 3 + 2 *
	bytecode := []int{
		int(opcodes.PSH), 5,
		int(opcodes.PSH), 3,
		int(opcodes.ADD),
		int(opcodes.PSH), 2,
		int(opcodes.MUL),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 16, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestVMTraceMode(t *testing.T) {
	st := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(st, rs)
	vm.EnableTrace() // Enable trace mode

	bytecode := []int{int(opcodes.PSH), 10, int(opcodes.PSH), 20, int(opcodes.ADD)}

	// In trace mode, VM prints to stdout (we won't assert output, just verify it doesn't crash)
	vm.Execute(bytecode)

	assert.Equal(t, 30, st.Pop())
}

func TestEqualityTrue(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 5,
		int(opcodes.PSH), 5,
		int(opcodes.EQ),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 1, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestEqualityFalse(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 5,
		int(opcodes.PSH), 3,
		int(opcodes.EQ),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 0, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestLessThanTrue(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 3,
		int(opcodes.PSH), 5,
		int(opcodes.LT),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 1, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestLessThanFalse(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 5,
		int(opcodes.PSH), 3,
		int(opcodes.LT),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 0, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestGreaterThanTrue(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 5,
		int(opcodes.PSH), 3,
		int(opcodes.GT),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 1, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestGreaterThanFalse(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 3,
		int(opcodes.PSH), 5,
		int(opcodes.GT),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 0, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestUnconditionalJump(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 1,
		int(opcodes.JMP), 6, // Jump to position 6 (skips next PUSH)
		int(opcodes.PSH), 99, // This gets skipped
		int(opcodes.PSH), 2, // Position 6: execution resumes here
	}

	vm.Execute(bytecode)

	assert.Equal(t, 2, s.Pop())
	assert.Equal(t, 1, s.Pop())
	assert.True(t, s.IsEmpty()) // 99 never pushed
}

func TestJumpIfZeroWhenZero(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 0, // Push 0 (false)
		int(opcodes.JZ), 6, // Jump to position 6 because top is 0
		int(opcodes.PSH), 99, // This gets skipped
		int(opcodes.PSH), 2, // Position 6: execution resumes here
	}

	vm.Execute(bytecode)

	assert.Equal(t, 2, s.Pop())
	assert.True(t, s.IsEmpty()) // 99 never pushed, 0 was popped by jz
}

func TestJumpIfZeroWhenNotZero(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 1, // Push 1 (true/non-zero)
		int(opcodes.JZ), 6, // Don't jump because top is not 0
		int(opcodes.PSH), 99, // This executes
		int(opcodes.PSH), 2, // Position 6: this also executes
	}

	vm.Execute(bytecode)

	assert.Equal(t, 2, s.Pop())
	assert.Equal(t, 99, s.Pop())
	assert.True(t, s.IsEmpty()) // 1 was popped by jz
}

func TestJumpIfNotZeroWhenNotZero(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 1, // Push 1 (non-zero)
		int(opcodes.JNZ), 6, // Jump to position 6 because top is not 0
		int(opcodes.PSH), 99, // This gets skipped
		int(opcodes.PSH), 2, // Position 6: execution resumes here
	}

	vm.Execute(bytecode)

	assert.Equal(t, 2, s.Pop())
	assert.True(t, s.IsEmpty()) // 99 never pushed, 1 was popped by jnz
}

func TestJumpIfNotZeroWhenZero(t *testing.T) {
	s := stack.NewStack()
	rs := registers.NewRegisters()
	vm := NewVM(s, rs)

	bytecode := ByteCode{
		int(opcodes.PSH), 0, // Push 0 (zero)
		int(opcodes.JNZ), 6, // Don't jump because top is 0
		int(opcodes.PSH), 99, // This executes
		int(opcodes.PSH), 2, // Position 6: this also executes
	}

	vm.Execute(bytecode)

	assert.Equal(t, 2, s.Pop())
	assert.Equal(t, 99, s.Pop())
	assert.True(t, s.IsEmpty()) // 0 was popped by jnz
}
