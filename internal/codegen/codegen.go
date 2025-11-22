package codegen

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/phasecurve/zhuji/internal/opcodes"
)

const (
	rax = "%rax"
	rbx = "%rbx"
	rcx = "%rcx"
	rdx = "%rdx"
	rsi = "%rsi"
	rdi = "%rdi"
	r8  = "%r8"
	r9  = "%r9"
	r10 = "%r10"
	r11 = "%r11"
	r12 = "%r12"
	r13 = "%r13"
	r14 = "%r14"
	r15 = "%r15"
	rbp = "%rbp"
	rip = "%rip"
)

var opCodeToX86Ops = map[opcodes.OpCode]string{
	opcodes.ADD: "addq", opcodes.SUB: "subq",
	opcodes.MUL: "imulq", opcodes.DIV: "idivq",
	opcodes.BEQ: "cmpq", opcodes.BLT: "cmpq", opcodes.BNE: "cmpq", opcodes.BGE: "cmpq",
	opcodes.MVQ: "movq",
}

var branchToJump = map[opcodes.OpCode]string{
	opcodes.BEQ: "je",
	opcodes.BLT: "jl",
	opcodes.BNE: "jne",
	opcodes.BGE: "jge",
}

var riscTox86Regs = map[int]string{
	0: "$0",
	1: rax, 2: rbx, 3: rcx, 4: rdx, 5: rsi,
	6: rdi, 7: r8, 8: r9, 9: r10, 10: r11,
	11: r12, 12: r13, 13: r14, 14: r15, 15: rbp,
}

type CodeGen struct {
	assembler    strings.Builder
	traceEnabled bool
}

func NewCodeGen() *CodeGen {
	cg := &CodeGen{
		traceEnabled: false,
	}
	cg.prependStart()
	return cg
}

func (c *CodeGen) emit(s string) {
	c.assembler.WriteString(s + "\n")
}

func (c *CodeGen) trace(format string, args ...any) {
	if c.traceEnabled {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			filePath := strings.Split(file, "/")
			file = filePath[len(filePath)-1]
			prefix := fmt.Sprintf("%s:%d ", file, line)
			fmt.Fprintf(os.Stderr, prefix+format, args...)
		}
	}
}

func (c *CodeGen) parseArithOp(op opcodes.OpCode, byteCode []int, ip int) int {
	rd := riscTox86Regs[byteCode[ip+1]]
	rs1 := riscTox86Regs[byteCode[ip+2]]
	rs2 := riscTox86Regs[byteCode[ip+3]]
	if op == opcodes.DIV || op == opcodes.MOD {
		resultReg := rax
		if op == opcodes.MOD {
			resultReg = rdx
		}
		if rs1 == rax {
			c.emit("cqto")
			c.emit(fmt.Sprintf("idivq %s", rs2))
			c.emit(fmt.Sprintf("movq %s, %s", resultReg, rd))
		} else if rs1 != rax && rs2 != rax {
			c.emit(fmt.Sprintf("movq %s, %s", rs1, rax))
			c.emit("cqto")
			c.emit(fmt.Sprintf("idivq %s", rs2))
			c.emit(fmt.Sprintf("movq %s, %s", resultReg, rd))
		} else if rs1 != rax && rs2 == rax {
			c.emit(fmt.Sprintf("xchgq %s, %s", rs1, rs2))
			c.emit("cqto")
			c.emit(fmt.Sprintf("idivq %s", rs1))
			c.emit(fmt.Sprintf("movq %s, %s", resultReg, rd))
		}
	} else {
		switch rd {
		case rs1:
			c.emit(fmt.Sprintf("%s %s, %s", opCodeToX86Ops[op], rs2, rd))
		case rs2:
			c.emit(fmt.Sprintf("%s %s, %s", opCodeToX86Ops[op], rs1, rd))
		default:
			c.emit(fmt.Sprintf("%s %s, %s", opCodeToX86Ops[opcodes.MVQ], rs1, rd))
			c.emit(fmt.Sprintf("%s %s, %s", opCodeToX86Ops[op], rs2, rd))
		}
	}
	return ip + 4
}

func (c *CodeGen) insertJumpLabel(branches map[int]string, ip int) {
	if label, ok := branches[ip]; ok {
		c.emit(fmt.Sprintf("%s:", label))
	}
}

func (c *CodeGen) Generate(byteCode []int) string {
	branches := c.findBranches(byteCode)
	for ip := 0; ip < len(byteCode); {
		c.insertJumpLabel(branches, ip)
		token := byteCode[ip]
		switch token {
		case int(opcodes.ADDI):
			dest := riscTox86Regs[byteCode[ip+1]]
			imm := byteCode[ip+3]
			op := byteCode[ip+2]
			if op == 0 {
				c.emit(fmt.Sprintf("movq $%d, %s", imm, dest))
			} else {
				c.emit(fmt.Sprintf("addq $%d, %s", imm, dest))
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
		case int(opcodes.BEQ):
			ip = c.branchOp(opcodes.BEQ, branches, ip, byteCode)
		case int(opcodes.BLT):
			ip = c.branchOp(opcodes.BLT, branches, ip, byteCode)
		case int(opcodes.BNE):
			ip = c.branchOp(opcodes.BNE, branches, ip, byteCode)
		case int(opcodes.BGE):
			ip = c.branchOp(opcodes.BGE, branches, ip, byteCode)
		case int(opcodes.SW):
			rs1 := riscTox86Regs[byteCode[ip+1]]
			offset := byteCode[ip+2]
			c.emit(fmt.Sprintf("%s %s, mem+%d(%s)", opCodeToX86Ops[opcodes.MVQ], rs1, offset, rip))
			ip += 4
		case int(opcodes.LW):
			rd := riscTox86Regs[byteCode[ip+1]]
			offset := byteCode[ip+2]
			c.emit(fmt.Sprintf("%s mem+%d(%s), %s", opCodeToX86Ops[opcodes.MVQ], offset, rip, rd))
			ip += 4
		}
	}

	c.insertJumpLabel(branches, len(byteCode))
	c.emit(fmt.Sprintf("movq %s, %s", rax, rdi))
	asm := c.appendExit()
	c.trace("asm:\n%s\n", asm)
	return asm
}

func (c *CodeGen) branchOp(op opcodes.OpCode, branches map[int]string, ip int, byteCode []int) int {
	rs1 := byteCode[ip+1]
	rs2 := byteCode[ip+2]
	offset := byteCode[ip+3]
	label := branches[ip+offset]
	c.emit(fmt.Sprintf("%s %s, %s", opCodeToX86Ops[op], riscTox86Regs[rs2], riscTox86Regs[rs1]))
	c.emit(fmt.Sprintf("%s %s", branchToJump[op], label))
	return ip + 4
}

func (c *CodeGen) findBranches(byteCode []int) map[int]string {
	branches := map[int]string{}
	for ip := 0; ip < len(byteCode); {
		opcode := byteCode[ip]
		if opcode != int(opcodes.BEQ) && opcode != int(opcodes.BLT) && opcode != int(opcodes.BNE) && opcode != int(opcodes.BGE) {
			ip += 4
			continue
		}
		jmpPos := ip + byteCode[ip+3]
		branches[jmpPos] = fmt.Sprintf("L%d", jmpPos)
		ip += 4
	}
	return branches
}

func (c *CodeGen) prependStart() {
	c.emit(".bss")
	c.emit("mem: .space 1024")
	c.emit(".text")
	c.emit(".global _start")
	c.emit("_start:")
}

func (c *CodeGen) appendExit() string {
	c.emit("movq $60, %rax")
	c.emit("syscall")
	return c.assembler.String()
}

func (c *CodeGen) toggleTraceOnOff() {
	c.traceEnabled = !c.traceEnabled
}
