package codegen

import (
	"fmt"
	"strings"

	"github.com/phasecurve/zhuji/internal/opcodes"
)

var opCodeToX86Ops = map[opcodes.OpCode]string{
	opcodes.ADD: "addq", opcodes.SUB: "subq",
	opcodes.MUL: "imulq", opcodes.DIV: "idivq",
	opcodes.MVQ: "movq",
}

var riscTox86Regs = map[int]string{
	1: "%rax", 2: "%rbx", 3: "%rcx", 4: "%rdx", 5: "%rsi",
	6: "%rdi", 7: "%r8", 8: "%r9", 9: "%r10", 10: "%r11",
	11: "%r12", 12: "%r13", 13: "%r14", 14: "%r15", 15: "%rbp",
}

type CodeGen struct {
	assembler strings.Builder
}

func NewCodeGen() *CodeGen {
	cg := &CodeGen{}
	cg.prependStart()
	return cg
}
func (c *CodeGen) emit(s string) {
	c.assembler.WriteString(s)
}

func (c *CodeGen) parseArithOp(op opcodes.OpCode, byteCode []int, ip int) int {
	rd := riscTox86Regs[byteCode[ip+1]]
	rs1 := riscTox86Regs[byteCode[ip+2]]
	rs2 := riscTox86Regs[byteCode[ip+3]]
	c.emit(fmt.Sprintf("%s %s, %s\n", opCodeToX86Ops[opcodes.MVQ], rs1, rd))
	if op == opcodes.DIV || op == opcodes.MOD {
		c.emit("cqto\n")
		if op == opcodes.MOD {
			c.emit(fmt.Sprintf("idivq %s\n", rs2))
			c.emit(fmt.Sprintf("movq %%rdx, %s\n", rd))
		} else {
			c.emit(fmt.Sprintf("%s %s\n", opCodeToX86Ops[op], rs2))
		}
	} else {
		c.emit(fmt.Sprintf("%s %s, %s\n", opCodeToX86Ops[op], rs2, rd))
	}
	return ip + 4
}

func (c *CodeGen) Generate(byteCode []int) string {
	for ip := 0; ip < len(byteCode); {
		token := byteCode[ip]
		switch token {
		case int(opcodes.ADDI):
			dest := riscTox86Regs[byteCode[ip+1]]
			imm := byteCode[ip+3]
			op := byteCode[ip+2]
			if op == 0 {
				c.emit(fmt.Sprintf("movq $%d, %s\n", imm, dest))
			} else {
				c.emit(fmt.Sprintf("addq $%d, %s\n", imm, dest))
			}
			ip += 4
		case int(opcodes.ADD):
			ip = c.parseArithOp(opcodes.ADD, byteCode, ip)
		case int(opcodes.SUB):
			ip = c.parseArithOp(opcodes.SUB, byteCode, ip)
		case int(opcodes.MUL):
			ip = c.parseArithOp(opcodes.MUL, byteCode, ip)
		case int(opcodes.DIV):
			ip = c.parseArithOp(opcodes.DIV, byteCode, ip)
		case int(opcodes.MOD):
			ip = c.parseArithOp(opcodes.MOD, byteCode, ip)
		}
	}
	c.emit("movq %rax, %rdi\n")

	return c.appendExit()
}
func (c *CodeGen) prependStart() {
	c.emit(".global _start\n")
	c.emit("_start:\n")
}

func (c *CodeGen) appendExit() string {
	c.emit("movq $60, %rax\n")
	c.emit("syscall\n")
	return c.assembler.String()
}
