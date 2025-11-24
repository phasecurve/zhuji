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
	opcodes.JAL: "call", opcodes.JALR: "ret",
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

func (c *CodeGen) parseArithOp(op opcodes.OpCode, bytecode []int, ip int) int {
	rd := riscTox86Regs[bytecode[ip+1]]
	rs1 := riscTox86Regs[bytecode[ip+2]]
	rs2 := riscTox86Regs[bytecode[ip+3]]
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

func (c *CodeGen) insertJumpLabel(branches map[int]string, functions map[int]string, ip int) {
	if label, ok := functions[ip]; ok {
		c.emit(fmt.Sprintf("%s:\npushq %%rbp", label))
	} else if label, ok := branches[ip]; ok {
		c.emit(fmt.Sprintf("%s:", label))
	}
}

func (c *CodeGen) Generate(bytecode []int) string {
	branches := c.findBranches(bytecode)
	functions := c.findFunctions(bytecode)
	for ip := 0; ip < len(bytecode); {
		c.insertJumpLabel(branches, functions, ip)
		if functions[ip] != "" {
			c.emit("movq %rsp, %rbp")
		}
		token := bytecode[ip]
		switch token {
		case int(opcodes.ADDI):
			rd := riscTox86Regs[bytecode[ip+1]]
			rs := riscTox86Regs[bytecode[ip+2]]
			imm := bytecode[ip+3]
			if rs == "$0" {
				c.emit(fmt.Sprintf("movq $%d, %s", imm, rd))
			} else {
				c.emit(fmt.Sprintf("movq %s, %s", rs, rd))
				if imm != 0 {
					c.emit(fmt.Sprintf("addq $%d, %s", imm, rd))
				}
			}
			ip += 4
		case int(opcodes.ADD):
			ip = c.parseArithOp(opcodes.ADD, bytecode, ip)
		case int(opcodes.SUB):
			ip = c.parseArithOp(opcodes.SUB, bytecode, ip)
		case int(opcodes.MUL):
			ip = c.parseArithOp(opcodes.MUL, bytecode, ip)
		case int(opcodes.DIV):
			ip = c.parseArithOp(opcodes.DIV, bytecode, ip)
		case int(opcodes.MOD):
			ip = c.parseArithOp(opcodes.MOD, bytecode, ip)
		case int(opcodes.BEQ):
			ip = c.branchOp(opcodes.BEQ, branches, ip, bytecode)
		case int(opcodes.BLT):
			ip = c.branchOp(opcodes.BLT, branches, ip, bytecode)
		case int(opcodes.BNE):
			ip = c.branchOp(opcodes.BNE, branches, ip, bytecode)
		case int(opcodes.BGE):
			ip = c.branchOp(opcodes.BGE, branches, ip, bytecode)
		case int(opcodes.SW):
			rs1 := riscTox86Regs[bytecode[ip+1]]
			offset := bytecode[ip+2]
			c.emit(fmt.Sprintf("%s %s, mem+%d(%s)", opCodeToX86Ops[opcodes.MVQ], rs1, offset, rip))
			ip += 4
		case int(opcodes.LW):
			rd := riscTox86Regs[bytecode[ip+1]]
			offset := bytecode[ip+2]
			c.emit(fmt.Sprintf("%s mem+%d(%s), %s", opCodeToX86Ops[opcodes.MVQ], offset, rip, rd))
			ip += 4
		case int(opcodes.JAL):
			offset := bytecode[ip+3]
			label := fmt.Sprintf("L%d", ip+offset)
			branches[ip+offset] = label
			c.emit(fmt.Sprintf("%s %s", opCodeToX86Ops[opcodes.JAL], label))
			if !strings.Contains(c.assembler.String(), "{{{syscall}}}") {
				c.emit("{{{syscall}}}")
			}
			ip += 4
		case int(opcodes.JALR):
			rd := bytecode[ip+1]
			if rd != 0 {
				panic("JALR with rd != 0 not supported in x86-64 codegen (only return pattern supported)")
			}
			if functions[ip] != "" {
				c.emit("movq %rbp, %rsp")
				c.emit("popq %rbp")
			}
			c.emit(opCodeToX86Ops[opcodes.JALR])
			ip += 4
		}
	}

	c.insertJumpLabel(branches, functions, len(bytecode))
	if !strings.Contains(c.assembler.String(), "{{{syscall}}}") {
		c.emit("{{{syscall}}}")
	}
	asm := c.appendExit()
	c.trace("asm:\n%s\n", asm)
	return asm
}

func (c *CodeGen) branchOp(op opcodes.OpCode, branches map[int]string, ip int, bytecode []int) int {
	rs1 := bytecode[ip+1]
	rs2 := bytecode[ip+2]
	offset := bytecode[ip+3]
	label := branches[ip+offset]
	c.emit(fmt.Sprintf("%s %s, %s", opCodeToX86Ops[op], riscTox86Regs[rs2], riscTox86Regs[rs1]))
	c.emit(fmt.Sprintf("%s %s", branchToJump[op], label))
	return ip + 4
}

func (c *CodeGen) findFunctions(bytecode []int) map[int]string {
	functions := map[int]string{}
	for ip := 0; ip < len(bytecode); {
		opcode := bytecode[ip]
		if opcode != int(opcodes.JAL) {
			ip += 4
			continue
		}

		jmpPos := ip + bytecode[ip+3]
		functions[jmpPos] = fmt.Sprintf("L%d", jmpPos)
		ip += 4
	}
	return functions
}

func (c *CodeGen) findBranches(bytecode []int) map[int]string {
	branches := map[int]string{}
	for ip := 0; ip < len(bytecode); {
		opcode := bytecode[ip]
		if opcode != int(opcodes.BEQ) && opcode != int(opcodes.BLT) && opcode != int(opcodes.BNE) && opcode != int(opcodes.BGE) {
			ip += 4
			continue
		}
		jmpPos := ip + bytecode[ip+3]
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
	asm := c.assembler.String()
	return strings.ReplaceAll(asm, "{{{syscall}}}", "movq %rax, %rdi\nmovq $60, %rax\nsyscall")
}

func (c *CodeGen) toggleTraceOnOff() {
	c.traceEnabled = !c.traceEnabled
}
