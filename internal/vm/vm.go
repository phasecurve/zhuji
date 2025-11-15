package vm

import (
	"github.com/phasecurve/zhuji/internal/stack"
)

type ByteCode []int

type vm struct {
	stack *stack.Stack
}

func NewVM(stack *stack.Stack) *vm {
	vm := &vm{
		stack: stack,
	}
	return vm
}

func (vm *vm) executeBinaryOp(op func(int, int) int) {
	operand2 := vm.stack.Pop()
	operand1 := vm.stack.Pop()
	vm.stack.Push(op(operand1, operand2))
}

func (vm *vm) Execute(byteCode ByteCode) {
	for ip := 0; ip < len(byteCode); {
		switch stack.OpCode(byteCode[ip]) {
		case stack.PUSH:
			ip++
			vm.stack.Push(byteCode[ip])
			ip++
		case stack.ADD:
			vm.executeBinaryOp(func(a, b int) int { return a + b })
			ip++
		case stack.SUB:
			vm.executeBinaryOp(func(a, b int) int { return a - b })
			ip++
		case stack.DIV:
			vm.executeBinaryOp(func(a, b int) int { return a / b })
			ip++
		case stack.MUL:
			vm.executeBinaryOp(func(a, b int) int { return a * b })
			ip++
		}
	}
}
