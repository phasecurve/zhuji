package codegen

import (
	"fmt"
	"strings"

	"github.com/phasecurve/zhuji/internal/opcodes"
)

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
			rd := riscTox86Regs[byteCode[ip+1]]
			rs1 := riscTox86Regs[byteCode[ip+2]]
			rs2 := riscTox86Regs[byteCode[ip+3]]
			c.emit(fmt.Sprintf("movq %s, %s\n", rs1, rd))
			c.emit(fmt.Sprintf("addq %s, %s\n", rs2, rd))
			ip += 4
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
