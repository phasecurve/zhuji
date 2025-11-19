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

func (a *Assembler) Assemble(assembly string) []int {
	noComments := a.removeComments(assembly)
	lines := ours.SplitRemoveEmpty(noComments, "\n")
	byteCode := []int{}
	for _, line := range lines {
		tks := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == ',' || r == '\t'
		})
		switch tks[0] {
		case "addi":
			byteCode = append(byteCode, int(opcodes.ADDI))
			byteCode = append(byteCode, reg[tks[1]])
			byteCode = append(byteCode, reg[tks[2]])
			if n, err := strconv.Atoi(tks[3]); err != nil {
				log.Fatalf("error while trying to parse an instruction: %v", err)
			} else {
				byteCode = append(byteCode, n)
			}
		case "add":
			byteCode = append(byteCode, int(opcodes.ADD))
			byteCode = append(byteCode, reg[tks[1]])
			byteCode = append(byteCode, reg[tks[2]])
			byteCode = append(byteCode, reg[tks[3]])
		case "sub":
			byteCode = append(byteCode, int(opcodes.SUB))
			byteCode = append(byteCode, reg[tks[1]])
			byteCode = append(byteCode, reg[tks[2]])
			byteCode = append(byteCode, reg[tks[3]])
		case "lw":
			byteCode = append(byteCode, int(opcodes.LW))
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
		case "sw":
			byteCode = append(byteCode, int(opcodes.SW))
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
		case "blt":
			byteCode = append(byteCode, int(opcodes.BLT))
			byteCode = append(byteCode, reg[tks[1]])
			byteCode = append(byteCode, reg[tks[2]])
			if n, err := strconv.Atoi(tks[3]); err != nil {
				log.Fatalf("error while trying to parse an instruction: %v", err)
			} else {
				byteCode = append(byteCode, n)
			}
		}
	}
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
