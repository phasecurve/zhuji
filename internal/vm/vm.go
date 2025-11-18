// Package vm the vm that executes the stream of byte codes.
package vm

import (
	"fmt"

	"github.com/phasecurve/zhuji/internal/memory"
	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/phasecurve/zhuji/internal/registers"
)

type ByteCode []int

type vm struct {
	registers    *registers.Registers
	memory       *memory.Memory
	traceEnabled bool
}

func NewVM(registers *registers.Registers, memory *memory.Memory) *vm {
	vm := &vm{
		registers: registers,
		memory:    memory,
	}
	return vm
}

func (vm *vm) Execute(byteCode ByteCode) {
	for ip := 0; ip < len(byteCode); {
		opCode := opcodes.OpCode(byteCode[ip])

		if vm.traceEnabled {
			fmt.Printf("[Execute] \n\tbyteCode: %v\n\topCode: %v\n\tip: %d", byteCode, opCode, ip)
		}
		switch opCode {
		case opcodes.ADDI:
			rd := byteCode[ip+1]
			rs := byteCode[ip+2]
			imm := byteCode[ip+3]
			result := vm.registers.Read(rs) + int32(imm)
			vm.registers.Write(rd, result)
			if vm.traceEnabled {
				fmt.Printf("[%d] ADDI x%d, x%d, %d → x%d = %d\n", ip, rd, rs, imm, rd, result)
			}
			ip += 4
		}
	}
}

func (vm *vm) EnableTrace() {
	vm.traceEnabled = true
}

func (vm *vm) DisableTrace() {
	vm.traceEnabled = false
}
