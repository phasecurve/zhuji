// Package vm the vm that executes the stream of byte codes.
package vm

import (
	"fmt"

	"github.com/phasecurve/zhuji/internal/memory"
	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/phasecurve/zhuji/internal/registers"
)

var opToAssemby = map[opcodes.OpCode]string{
	opcodes.ADDI: "addi",
	opcodes.ADD:  "add",
	opcodes.SUB:  "sub",
	opcodes.LW:   "lw",
	opcodes.SW:   "sw",
	opcodes.BNE:  "bne",
	opcodes.BGE:  "bge",
	opcodes.BEQ:  "beq",
	opcodes.BLT:  "blt",
}

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

func (vm *vm) execRegImmOp(opCode opcodes.OpCode, byteCode []int, ip int) int {
	rd := byteCode[ip+1]
	rs := byteCode[ip+2]
	imm := byteCode[ip+3]
	result := vm.registers.Read(rs) + int32(imm)
	vm.registers.Write(rd, result)
	if vm.traceEnabled {
		fmt.Printf("[%d] %s x%d, x%d, %d → x%d = %d\n", ip, opToAssemby[opCode], rd, rs, imm, rd, result)
	}
	return 4
}

func (vm *vm) execRegOp(opCode opcodes.OpCode, byteCode []int, ip int, op func(int32, int32) int32) int {
	rd := byteCode[ip+1]
	rs1 := byteCode[ip+2]
	rs2 := byteCode[ip+3]
	result := op(vm.registers.Read(rs1), vm.registers.Read(rs2))
	vm.registers.Write(rd, result)
	if vm.traceEnabled {
		fmt.Printf("[%d] %s x%d, x%d, x%d → x%d = %d\n", ip, opToAssemby[opCode], rd, rs1, rs2, rd, result)
	}
	return 4
}
func (vm *vm) execBranch(opCode opcodes.OpCode, byteCode []int, ip int, cond func(int32, int32) bool) int {
	rs1Val := vm.registers.Read(byteCode[ip+1])
	rs2Val := vm.registers.Read(byteCode[ip+2])
	target := byteCode[ip+3]
	nextIP := ip + 4
	if cond(rs1Val, rs2Val) {
		nextIP = target
	}
	if vm.traceEnabled {
		fmt.Printf("[%d] %s x%d, x%d, %d → ip = %d\n", ip, opToAssemby[opCode], byteCode[ip+1], byteCode[ip+2], target,
			nextIP)
	}
	return nextIP
}

func (vm *vm) Execute(byteCode ByteCode) {
	for ip := 0; ip < len(byteCode); {
		opCode := opcodes.OpCode(byteCode[ip])

		if vm.traceEnabled {
			fmt.Printf("[Execute] \n\tbyteCode: %v\n\topCode: %v\n\tip: %d", byteCode, opCode, ip)
		}
		switch opCode {
		case opcodes.ADDI:
			ip += vm.execRegImmOp(opCode, byteCode, ip)
		case opcodes.ADD:
			ip += vm.execRegOp(opCode, byteCode, ip, func(v1, v2 int32) int32 {
				return v1 + v2
			})
		case opcodes.SUB:
			ip += vm.execRegOp(opCode, byteCode, ip, func(v1, v2 int32) int32 {
				return v1 - v2
			})
		case opcodes.MUL:
			ip += vm.execRegOp(opCode, byteCode, ip, func(v1, v2 int32) int32 {
				return v1 * v2
			})
		case opcodes.DIV:
			ip += vm.execRegOp(opCode, byteCode, ip, func(v1, v2 int32) int32 {
				return v1 / v2
			})
		case opcodes.MOD:
			ip += vm.execRegOp(opCode, byteCode, ip, func(v1, v2 int32) int32 {
				return v1 % v2
			})
		case opcodes.SW:
			rs2 := byteCode[ip+1]
			offset := byteCode[ip+2]
			rs1 := byteCode[ip+3]
			val := vm.registers.Read(rs2)
			addr := int(vm.registers.Read(rs1)) + offset
			vm.memory.StoreWord(addr, val)
			if vm.traceEnabled {
				fmt.Printf("[%d] sw x%d, %d(x%d) → x%d = %d\n", ip, rs2, offset, rs1, addr, val)
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
				fmt.Printf("[%d] lw x%d, %d(x%d) → x%d = %d\n", ip, rd, offset, rs, addr, val)
			}
			ip += 4
		case opcodes.BLT:
			ip = vm.execBranch(opCode, byteCode, ip, func(v1 int32, v2 int32) bool { return v1 < v2 })
		case opcodes.BEQ:
			ip = vm.execBranch(opCode, byteCode, ip, func(v1 int32, v2 int32) bool { return v1 == v2 })
		case opcodes.BNE:
			ip = vm.execBranch(opCode, byteCode, ip, func(v1 int32, v2 int32) bool { return v1 != v2 })
		case opcodes.BGE:
			ip = vm.execBranch(opCode, byteCode, ip, func(v1 int32, v2 int32) bool { return v1 >= v2 })
		}
	}
}

func (vm *vm) EnableTrace() {
	vm.traceEnabled = true
}

func (vm *vm) DisableTrace() {
	vm.traceEnabled = false
}
