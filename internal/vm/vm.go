// Package vm the vm that executes the stream of byte codes.
package vm

import (
	"fmt"

	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/phasecurve/zhuji/internal/registers"
	"github.com/phasecurve/zhuji/internal/stack"
)

type ByteCode []int

type vm struct {
	stack        *stack.Stack
	registers    *registers.Registers
	traceEnabled bool
}

func NewVM(stack *stack.Stack, registers *registers.Registers) *vm {
	vm := &vm{
		stack:     stack,
		registers: registers,
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
		opCode := opcodes.OpCode(byteCode[ip])

		if vm.traceEnabled {
			fmt.Printf("[Execute:39] \n\tbyteCode: %v\n\topCode: %v\n\tip: %d", byteCode, opCode, ip)
		}
		switch opCode {

		case opcodes.PSH:
			ip++
			bc := byteCode[ip]
			vm.stack.Push(bc)
			if vm.traceEnabled {
				fmt.Printf("[%d] PUSH %d    → %s\n", ip, bc, vm.stack.String())
			}
			ip++
		case opcodes.JMP:
			ip++
			bc := byteCode[ip]
			if vm.traceEnabled {
				fmt.Printf("[%d] JMP %d    → %s\n", ip, bc, vm.stack.String())
			}
			ip = bc
		case opcodes.JZ:
			ip++
			bc := byteCode[ip]
			if vm.traceEnabled {
				fmt.Printf("[%d] JZ %d    → %s\n", ip, bc, vm.stack.String())
			}
			if vm.stack.Pop() == 0 {
				ip = bc
			} else {
				ip++
			}
		case opcodes.JNZ:
			ip++
			bc := byteCode[ip]
			if vm.traceEnabled {
				fmt.Printf("[%d] JNZ %d    → %s\n", ip, bc, vm.stack.String())
			}
			if vm.stack.Pop() != 0 {
				ip = bc
			} else {
				ip++
			}
		case opcodes.ADD:
			vm.executeBinaryOp(func(a, b int) int { return a + b })
			if vm.traceEnabled {
				fmt.Printf("[%d] ADD        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.SUB:
			vm.executeBinaryOp(func(a, b int) int { return a - b })
			if vm.traceEnabled {
				fmt.Printf("[%d] SUB        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.DIV:
			vm.executeBinaryOp(func(a, b int) int { return a / b })
			if vm.traceEnabled {
				fmt.Printf("[%d] DIV        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.MUL:
			vm.executeBinaryOp(func(a, b int) int { return a * b })
			if vm.traceEnabled {
				fmt.Printf("[%d] MUL        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.EQ:
			vm.executeBinaryOp(func(a, b int) int {
				if a == b {
					return 1
				}
				return 0
			})
			if vm.traceEnabled {
				fmt.Printf("[%d] EQ        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.LT:
			vm.executeBinaryOp(func(a, b int) int {
				if a < b {
					return 1
				}
				return 0
			})
			if vm.traceEnabled {
				fmt.Printf("[%d] LT        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.LTE:
			vm.executeBinaryOp(func(a, b int) int {
				if a <= b {
					return 1
				}
				return 0
			})
			if vm.traceEnabled {
				fmt.Printf("[%d] LTE        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.GT:
			vm.executeBinaryOp(func(a, b int) int {
				if a > b {
					return 1
				}
				return 0
			})
			if vm.traceEnabled {
				fmt.Printf("[%d] GT        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.GTE:
			vm.executeBinaryOp(func(a, b int) int {
				if a >= b {
					return 1
				}
				return 0
			})
			if vm.traceEnabled {
				fmt.Printf("[%d] GTE        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.DUP:
			vm.stack.Dup()
			if vm.traceEnabled {
				fmt.Printf("[%d] DUP        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.SWP:
			vm.stack.Swap()
			if vm.traceEnabled {
				fmt.Printf("[%d] SWP        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.DRP:
			vm.stack.Drop()
			if vm.traceEnabled {
				fmt.Printf("[%d] DRP        → %s\n", ip, vm.stack.String())
			}
			ip++
		case opcodes.STREG:
			top := vm.stack.Pop()
			vm.registers.Write(1, int32(top))
			if vm.traceEnabled {
				fmt.Printf("[%d] STREG      → %s\n", ip, vm.stack.String())
			}
		}
	}
}

func (vm *vm) EnableTrace() {
	vm.traceEnabled = true
}

func (vm *vm) DisableTrace() {
	vm.traceEnabled = false
}
