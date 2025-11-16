// Package vm the vm that executes the stream of byte codes.
package vm

import (
	"fmt"

	"github.com/phasecurve/zhuji/internal/stack"
)

type ByteCode []int

type vm struct {
	stack        *stack.Stack
	traceEnabled bool
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
		opCode := stack.OpCode(byteCode[ip])

		switch opCode {

		case stack.PSH:
			ip++
			bc := byteCode[ip]
			vm.stack.Push(bc)
			if vm.traceEnabled {
				fmt.Printf("[%d] PUSH %d    → %s\n", ip, bc, vm.stack.String())
			}
			ip++
		case stack.ADD:
			vm.executeBinaryOp(func(a, b int) int { return a + b })
			if vm.traceEnabled {
				fmt.Printf("[%d] ADD        → %s\n", ip, vm.stack.String())
			}
			ip++
		case stack.SUB:
			vm.executeBinaryOp(func(a, b int) int { return a - b })
			if vm.traceEnabled {
				fmt.Printf("[%d] SUB        → %s\n", ip, vm.stack.String())
			}
			ip++
		case stack.DIV:
			vm.executeBinaryOp(func(a, b int) int { return a / b })
			if vm.traceEnabled {
				fmt.Printf("[%d] DIV        → %s\n", ip, vm.stack.String())
			}
			ip++
		case stack.MUL:
			vm.executeBinaryOp(func(a, b int) int { return a * b })
			if vm.traceEnabled {
				fmt.Printf("[%d] MUL        → %s\n", ip, vm.stack.String())
			}
			ip++
		case stack.DUP:
			vm.stack.Dup()
			ip++
		case stack.SWP:
			vm.stack.Swap()
			ip++
		case stack.DRP:
			vm.stack.Drop()
			ip++
		}
	}
}

func (vm *vm) EnableTrace() {
	vm.traceEnabled = true
}

func (vm *vm) DisableTrace() {
	vm.traceEnabled = false
}
