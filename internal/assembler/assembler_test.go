package assembler

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/stack"
	"github.com/phasecurve/zhuji/internal/vm"
	"github.com/stretchr/testify/assert"
)

var a = &Assembler{
	traceEnabled: false,
}
var Assemble = a.Assemble

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

func TestAssembleAndExecuteGT(t *testing.T) {
	input := `push 5
      push 3
      gt`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 1, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestAssembleAndExecuteGTE(t *testing.T) {
	input := `push 5
      push 5
      gte`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 1, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestAssembleAndExecuteLTE(t *testing.T) {
	input := `push 5
      push 5
      lte`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 1, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestAssembleAndExecuteLT(t *testing.T) {
	input := `push 3
      push 5
      lt`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 1, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestAssembleAndExecuteEQ(t *testing.T) {
	input := `push 5
      push 5
      eq`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 1, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestAssembleJump(t *testing.T) {
	input := `push 1
      jmp 6
      push 99
      push 2`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	expected := []int{
		int(stack.PSH), 1,
		int(stack.JMP), 6,
		int(stack.PSH), 99,
		int(stack.PSH), 2,
	}

	assert.Equal(t, expected, bytecode)
}

func TestAssembleLabelWithJump(t *testing.T) {
	input := `push 1
      jmp skip
      push 99
  skip:
      push 2`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	expected := []int{
		int(stack.PSH), 1,
		int(stack.JMP), 6, // skip is at position 6
		int(stack.PSH), 99,
		int(stack.PSH), 2,
	}

	assert.Equal(t, expected, bytecode)
}

func TestAssembleLabelWithJumpIfZero(t *testing.T) {
	input := `push 0
      jz skip
      push 99
  skip:
      push 2`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	expected := []int{
		int(stack.PSH), 0,
		int(stack.JZ), 6, // jump because top 0 skip is at position 6
		int(stack.PSH), 99,
		int(stack.PSH), 2,
	}

	assert.Equal(t, expected, bytecode)
}

func TestAssembleLabelWithJumpIfNotZero(t *testing.T) {
	input := `push 1
      jnz skip
      push 99
  skip:
      push 2`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	expected := []int{
		int(stack.PSH), 1,
		int(stack.JNZ), 6, // jump because top 0 skip is at position 6
		int(stack.PSH), 99,
		int(stack.PSH), 2,
	}

	assert.Equal(t, expected, bytecode)
}

func TestExecuteJumpIfZeroWithLabelWhenZero(t *testing.T) {
	input := `push 0
      jz skip
      push 99
  skip:
      push 2`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 2, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteJumpIfZeroWithLabelWhenNotZero(t *testing.T) {
	input := `push 1
      jz skip
      push 99
  skip:
      push 2`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 2, st.Pop())
	assert.Equal(t, 99, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteJumpIfNotZeroWithLabelWhenNotZero(t *testing.T) {
	input := `push 1
      jnz skip
      push 99
  skip:
      push 2`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 2, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestExecuteJumpIfNotZeroWithLabelWhenZero(t *testing.T) {
	input := `push 0
      jnz skip
      push 99
  skip:
      push 2`

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 2, st.Pop())
	assert.Equal(t, 99, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestProgramCountToTen(t *testing.T) {
	input := `push 0
  loop:
      push 1
      add
      dup
      push 10
      gt
      jnz done
      jmp loop
  done:
      `

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 11, st.Pop())
	assert.True(t, st.IsEmpty())
}

func TestProgramSumToN(t *testing.T) {
	input := `push 0
      push 0
  loop:
      push 1
      add
      swap
      dup
      add
      swap
      dup
      push 5
      gt
      jz loop
      swap
      drop
      `

	bytecode, err := Assemble(input)
	assert.NoError(t, err)

	st := stack.NewStack()
	vm := vm.NewVM(st)
	vm.Execute(bytecode)

	assert.Equal(t, 15, st.Pop())
	assert.True(t, st.IsEmpty())
}
