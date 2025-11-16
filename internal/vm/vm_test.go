package vm

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/stack"
	"github.com/stretchr/testify/assert"
)

func TestExecuteSinglePushInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	pushOp := []int{int(stack.PUSH), 42}

	vm.Execute(pushOp)

	assert.Equal(t, st.Peek(), 42, "should have executed a push instruction")
}

func TestExecuteMultiplePushInstructions(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	bytecode := []int{int(stack.PUSH), 10, int(stack.PUSH), 20}

	vm.Execute(bytecode)

	assert.Equal(t, 20, st.Pop())
	assert.Equal(t, 10, st.Pop())
}

func TestExecuteAddInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// PUSH 10, PUSH 20, ADD
	bytecode := []int{int(stack.PUSH), 10, int(stack.PUSH), 20, int(stack.ADD)}

	vm.Execute(bytecode)

	assert.Equal(t, 30, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteSubInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// PUSH 20, PUSH 10, SUB → (20 - 10 = 10)
	bytecode := []int{int(stack.PUSH), 20, int(stack.PUSH), 10, int(stack.SUB)}

	vm.Execute(bytecode)

	assert.Equal(t, 10, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteDivInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// PUSH 20, PUSH 4, DIV → (20 / 4 = 5)
	bytecode := []int{int(stack.PUSH), 20, int(stack.PUSH), 4, int(stack.DIV)}

	vm.Execute(bytecode)

	assert.Equal(t, 5, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteMulInstruction(t *testing.T) {
	st := stack.NewStack()
	vm := NewVM(st)
	// PUSH 20, PUSH 4, MUL → (20 * 4 = 80)
	bytecode := []int{int(stack.PUSH), 20, int(stack.PUSH), 4, int(stack.MUL)}

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
		int(stack.PUSH), 5,
		int(stack.PUSH), 3,
		int(stack.ADD),
		int(stack.PUSH), 2,
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

	bytecode := []int{int(stack.PUSH), 10, int(stack.PUSH), 20, int(stack.ADD)}

	// In trace mode, VM prints to stdout (we won't assert output, just verify it doesn't crash)
	vm.Execute(bytecode)

	assert.Equal(t, 30, st.Pop())
}
