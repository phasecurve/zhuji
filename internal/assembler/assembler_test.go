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
	assert.Equal(t, []int{int(stack.PSH), 42}, bytecode)
}

func TestAssembleMultiplePushes(t *testing.T) {
	input := `push 10
	push 20`
	bytecode, err := Assemble(input)

	assert.NoError(t, err)
	assert.Equal(t, []int{int(stack.PSH), 10, int(stack.PSH), 20}, bytecode)
}
func TestAssembleAddInstruction(t *testing.T) {
	input := `push 10
	push 20
	add`
	bytecode, err := Assemble(input)

	assert.NoError(t, err)
	expected := []int{
		int(stack.PSH), 10,
		int(stack.PSH), 20,
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
		int(stack.PSH), 20,
		int(stack.PSH), 10,
		int(stack.SUB),
		int(stack.PSH), 2,
		int(stack.MUL),
		int(stack.PSH), 4,
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
		int(stack.PSH), 42,
		int(stack.PSH), 10,
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

func TestDupInstruction(t *testing.T) {
	input := `push 42
	dup`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	// Execute on VM
	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	// Stack should have [42, 42]
	assert.Equal(t, 42, st.Pop())
	assert.Equal(t, 42, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestSwapInstruction(t *testing.T) {
	input := `push 10
  push 20
  swap`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	// Execute on VM
	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	// Stack should have [20, 10] (10 on top now)
	assert.Equal(t, 10, st.Pop())
	assert.Equal(t, 20, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestDropInstruction(t *testing.T) {
	input := `push 10
  push 20
  push 30
  drop`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	// Execute on VM
	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	// Stack should have [10, 20] (30 was dropped)
	assert.Equal(t, 20, st.Pop())
	assert.Equal(t, 10, st.Pop())
	assert.True(t, st.IsEmpty())
}
