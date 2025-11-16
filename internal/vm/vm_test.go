package vm

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/stack"
	"github.com/stretchr/testify/assert"
)

func TestExecuteSinglePushInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	pushOp := []int{int(stack.PSH), 42}

	vm.Execute(pushOp)

	assert.Equal(t, st.Peek(), 42, "should have executed a push instruction")
}

func TestExecuteMultiplePushInstructions(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	bytecode := []int{int(stack.PSH), 10, int(stack.PSH), 20}

	vm.Execute(bytecode)

	assert.Equal(t, 20, st.Pop())
	assert.Equal(t, 10, st.Pop())
}

func TestExecuteAddInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// PSH 10, PSH 20, ADD
	bytecode := []int{int(stack.PSH), 10, int(stack.PSH), 20, int(stack.ADD)}

	vm.Execute(bytecode)

	assert.Equal(t, 30, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteSubInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// PSH 20, PSH 10, SUB → (20 - 10 = 10)
	bytecode := []int{int(stack.PSH), 20, int(stack.PSH), 10, int(stack.SUB)}

	vm.Execute(bytecode)

	assert.Equal(t, 10, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteDivInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// PSH 20, PSH 4, DIV → (20 / 4 = 5)
	bytecode := []int{int(stack.PSH), 20, int(stack.PSH), 4, int(stack.DIV)}

	vm.Execute(bytecode)

	assert.Equal(t, 5, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteMulInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// PSH 20, PSH 4, MUL → (20 * 4 = 80)
	bytecode := []int{int(stack.PSH), 20, int(stack.PSH), 4, int(stack.MUL)}

	vm.Execute(bytecode)

	assert.Equal(t, 80, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteComplexExpression(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// (5 + 3) * 2 = 16
	// RPN: 5 3 + 2 *
	bytecode := []int{
		int(stack.PSH), 5,
		int(stack.PSH), 3,
		int(stack.ADD),
		int(stack.PSH), 2,
		int(stack.MUL),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 16, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestVMTraceMode(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	vm.EnableTrace() // Enable trace mode

	bytecode := []int{int(stack.PSH), 10, int(stack.PSH), 20, int(stack.ADD)}

	// In trace mode, VM prints to stdout (we won't assert output, just verify it doesn't crash)
	vm.Execute(bytecode)

	assert.Equal(t, 30, st.Pop())
}

func TestEqualityTrue(t *testing.T) {
	s := stack.NewStack()
	vm := NewVM(s)

	bytecode := ByteCode{
		int(stack.PSH), 5,
		int(stack.PSH), 5,
		int(stack.EQ),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 1, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestEqualityFalse(t *testing.T) {
	s := stack.NewStack()
	vm := NewVM(s)

	bytecode := ByteCode{
		int(stack.PSH), 5,
		int(stack.PSH), 3,
		int(stack.EQ),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 0, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestLessThanTrue(t *testing.T) {
	s := stack.NewStack()
	vm := NewVM(s)

	bytecode := ByteCode{
		int(stack.PSH), 3,
		int(stack.PSH), 5,
		int(stack.LT),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 1, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestLessThanFalse(t *testing.T) {
	s := stack.NewStack()
	vm := NewVM(s)

	bytecode := ByteCode{
		int(stack.PSH), 5,
		int(stack.PSH), 3,
		int(stack.LT),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 0, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestGreaterThanTrue(t *testing.T) {
	s := stack.NewStack()
	vm := NewVM(s)

	bytecode := ByteCode{
		int(stack.PSH), 5,
		int(stack.PSH), 3,
		int(stack.GT),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 1, s.Pop())
	assert.True(t, s.IsEmpty())
}

func TestGreaterThanFalse(t *testing.T) {
	s := stack.NewStack()
	vm := NewVM(s)

	bytecode := ByteCode{
		int(stack.PSH), 3,
		int(stack.PSH), 5,
		int(stack.GT),
	}

	vm.Execute(bytecode)

	assert.Equal(t, 0, s.Pop())
	assert.True(t, s.IsEmpty())
}
