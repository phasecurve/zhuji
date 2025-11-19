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
		case opcodes.ADD:
			rd := byteCode[ip+1]
			rs1 := byteCode[ip+2]
			rs2 := byteCode[ip+3]
			result := vm.registers.Read(rs1) + vm.registers.Read(rs2)
			vm.registers.Write(rd, result)
			if vm.traceEnabled {
				fmt.Printf("[%d] ADD x%d, x%d, x%d → x%d = %d\n", ip, rd, rs1, rs2, rd, result)
			}
			ip += 4
		case opcodes.SUB:
			rd := byteCode[ip+1]
			rs1 := byteCode[ip+2]
			rs2 := byteCode[ip+3]
			result := vm.registers.Read(rs1) - vm.registers.Read(rs2)
			vm.registers.Write(rd, result)
			if vm.traceEnabled {
				fmt.Printf("[%d] SUB x%d, x%d, x%d → x%d = %d\n", ip, rd, rs1, rs2, rd, result)
			}
			ip += 4
		case opcodes.SW:
			rs2 := byteCode[ip+1]
			offset := byteCode[ip+2]
			rs1 := byteCode[ip+3]
			val := vm.registers.Read(rs2)
			addr := int(vm.registers.Read(rs1)) + offset
			vm.memory.StoreWord(addr, val)
			if vm.traceEnabled {
				fmt.Printf("[%d] SW x%d, %d(x%d) → x%d = %d\n", ip, rs2, offset, rs1, addr, val)
			}
			ip += 4
		case opcodes.LW:
			rd := byteCode[ip+1]
			offset := byteCode[ip+2]
			rs := byteCode[ip+3]
			addr := int(vm.registers.Read(rs)) + offset
			val := vm.memory.LoadWord(addr)
			vm.registers.Write(rd, val)
			if vm.traceEnabled {
				fmt.Printf("[%d] LW x%d, %d(x%d) → x%d = %d\n", ip, rd, offset, rs, addr, val)
			}
			ip += 4
		case opcodes.BLT:
			rs1 := vm.registers.Read(byteCode[ip+1])
			rs2 := vm.registers.Read(byteCode[ip+2])
			target := byteCode[ip+3]
			addr := ip + 4
			if rs1 < rs2 {
				addr = target
			}
			if vm.traceEnabled {
				fmt.Printf("[%d] BLT x%d, x%d, %d → ip = %d\n", ip, rs1, rs2, target, addr)
			}
			ip = addr
		}
	}
}

func (vm *vm) EnableTrace() {
	vm.traceEnabled = true
}

func (vm *vm) DisableTrace() {
	vm.traceEnabled = false
}
