package assembler

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/stack"
	"github.com/phasecurve/zhuji/internal/vm"
	"github.com/stretchr/testify/assert"
)

func TestAssembleSinglePush(t *testing.T) {
	input := "push 42"
	bytecode, err := Assemble(input)

	assert.NoError(t, err)
	assert.Equal(t, []int{int(stack.PUSH), 42}, bytecode)
}

func TestAssembleMultiplePushes(t *testing.T) {
	input := `push 10
	push 20`
	bytecode, err := Assemble(input)

	assert.NoError(t, err)
	assert.Equal(t, []int{int(stack.PUSH), 10, int(stack.PUSH), 20}, bytecode)
}
func TestAssembleAddInstruction(t *testing.T) {
	input := `push 10
	push 20
	add`
	bytecode, err := Assemble(input)

	assert.NoError(t, err)
	expected := []int{
		int(stack.PUSH), 10,
		int(stack.PUSH), 20,
		int(stack.ADD),
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleAllArithmeticOperations(t *testing.T) {
	input := `push 20
	push 10
	sub
	push 2
	mul
	push 4
	div`
	bytecode, err := Assemble(input)

	assert.NoError(t, err)
	expected := []int{
		int(stack.PUSH), 20,
		int(stack.PUSH), 10,
		int(stack.SUB),
		int(stack.PUSH), 2,
		int(stack.MUL),
		int(stack.PUSH), 4,
		int(stack.DIV),
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleWithComments(t *testing.T) {
	input := `# This is a comment
	push 42  # Push the answer
	push 10
	add      # Add them together`

	bytecode, err := Assemble(input)

	assert.NoError(t, err)
	expected := []int{
		int(stack.PUSH), 42,
		int(stack.PUSH), 10,
		int(stack.ADD),
	}
	assert.Equal(t, expected, bytecode)
}

func TestAssembleAndExecute(t *testing.T) {
	// Calculate (5 + 3) * 2
	input := `# Calculate (5 + 3) * 2
	push 5
	push 3
	add
	push 2
	mul`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	// Execute on VM
	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 16, st.Pop())
	assert.True(t, st.IsEmpty())
}
