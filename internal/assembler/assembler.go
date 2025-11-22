package assembler

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	ours "github.com/phasecurve/zhuji/internal"
	"github.com/phasecurve/zhuji/internal/opcodes"
)

type Assembler struct {
	traceEnabled bool
}

var reg = map[string]int{
	"x0": 0, "x1": 1, "x2": 2, "x3": 3, "x4": 4, "x5": 5, "x6": 6, "x7": 7,
	"x8": 8, "x9": 9, "x10": 10, "x11": 11, "x12": 12, "x13": 13, "x14": 14, "x15": 15,
	"x16": 16, "x17": 17, "x18": 18, "x19": 19, "x20": 20, "x21": 21, "x22": 22, "x23": 23,
	"x24": 24, "x25": 25, "x26": 26, "x27": 27, "x28": 28, "x29": 29, "x30": 30, "x31": 31,
}

func NewAssembler() *Assembler {
	return &Assembler{}
}

func (a *Assembler) findLabels(lines []string) map[string]int {
	labels := map[string]int{}

	ip := 0
	for _, line := range lines {
		if strings.Contains(line, ":") {
			line = ours.TrimSuffix(line, ':')
			labels[line] = ip
		} else {
			ip += 4
		}
	}

	return labels
}

func (a *Assembler) Assemble(assembly string) []int {
	noComments := a.removeComments(assembly)
	lines := ours.SplitRemoveEmpty(noComments, "\n")
	byteCode := []int{}
	labels := a.findLabels(lines)
	ip := 0
	for _, line := range lines {
		tks := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == ',' || r == '\t'
		})
		if strings.Contains(tks[0], ":") {
			continue
		}
		switch tks[0] {
		case "addi":
			byteCode = handleImmediateOp3(opcodes.ADDI, byteCode, tks)
		case "add":
			byteCode = handleRegistersOp3(opcodes.ADD, byteCode, tks)
		case "sub":
			byteCode = handleRegistersOp3(opcodes.SUB, byteCode, tks)
		case "mul":
			byteCode = handleRegistersOp3(opcodes.MUL, byteCode, tks)
		case "div":
			byteCode = handleRegistersOp3(opcodes.DIV, byteCode, tks)
		case "mod":
			byteCode = handleRegistersOp3(opcodes.MOD, byteCode, tks)
		case "lw":
			byteCode = handleLoadOrStore(opcodes.LW, byteCode, tks)
		case "sw":
			byteCode = handleLoadOrStore(opcodes.SW, byteCode, tks)
		case "blt":
			byteCode = handleBranchOp(opcodes.BLT, byteCode, tks, labels, ip)
		case "beq":
			byteCode = handleBranchOp(opcodes.BEQ, byteCode, tks, labels, ip)
		case "bne":
			byteCode = handleBranchOp(opcodes.BNE, byteCode, tks, labels, ip)
		case "bge":
			byteCode = handleBranchOp(opcodes.BGE, byteCode, tks, labels, ip)
		}
		ip += 4
	}
	return byteCode
}

func handleRegistersOp3(op opcodes.OpCode, byteCode []int, tks []string) []int {
	byteCode = append(byteCode, int(op))
	byteCode = append(byteCode, reg[tks[1]])
	byteCode = append(byteCode, reg[tks[2]])
	byteCode = append(byteCode, reg[tks[3]])
	return byteCode
}

func handleImmediateOp3(op opcodes.OpCode, byteCode []int, tks []string) []int {
	byteCode = append(byteCode, int(op))
	byteCode = append(byteCode, reg[tks[1]])
	byteCode = append(byteCode, reg[tks[2]])
	if n, err := strconv.Atoi(tks[3]); err != nil {
		log.Fatalf("error while trying to parse an instruction: %v", err)
	} else {
		byteCode = append(byteCode, n)
	}
	return byteCode
}

func handleBranchOp(op opcodes.OpCode, byteCode []int, tks []string, labels map[string]int, ip int) []int {
	if pos, ok := labels[tks[3]]; ok {
		tks[3] = strconv.Itoa(pos - ip)
	}
	return handleImmediateOp3(op, byteCode, tks)
}

func handleLoadOrStore(op opcodes.OpCode, byteCode []int, tks []string) []int {
	byteCode = append(byteCode, int(op))
	byteCode = append(byteCode, reg[tks[1]])
	offsetAndBase := strings.FieldsFunc(tks[2], func(r rune) bool {
		return r == '(' || r == ')'
	})
	offset, err := strconv.Atoi(offsetAndBase[0])
	if err != nil {
		log.Fatalf("error attempting to parse offset: %v err: %v", offsetAndBase[0], err)
	}
	byteCode = append(byteCode, offset)
	byteCode = append(byteCode, reg[offsetAndBase[1]])
	return byteCode
}

func stripComment(line string) string {
	commentStart := strings.Index(line, "#")
	if commentStart == -1 {
		commentStart = len(line)
	}
	assemblerSplit := strings.Fields(line[:commentStart])

	return strings.Join(assemblerSplit, " ")
}

func (a *Assembler) removeComments(input string) string {
	lines := strings.Split(input, "\n")
	linesNoComments := strings.Builder{}
	for _, line := range lines {
		linesNoComments.WriteString(fmt.Sprintf("%s\n", stripComment(line)))
	}
	result := linesNoComments.String()
	if a.traceEnabled {
		fmt.Printf("[removeComments:66] %s\n", result)
	}
	return result
}
