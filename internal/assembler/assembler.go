// Package assembler will take risk-v/gas style assembler and conver to bytecode
package assembler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phasecurve/zhuji/internal/stack"
)

type Assembler struct {
	traceEnabled bool
}

func NewAssembler(traceEnable bool) Assembler {
	assembler := Assembler{
		traceEnabled: traceEnable,
	}
	return assembler
}

func (a *Assembler) buildSymbolTable(lines []string) map[string]int {
	symbolTable := make(map[string]int)
	pos := 0
	for _, inst := range lines {
		inst = strings.TrimSpace(inst)
		token := strings.Split(inst, " ")
		if token[0] == "push" || token[0] == "jmp" || token[0] == "jz" || token[0] == "jnz" {
			pos += 2
		} else if jumpVal, found := strings.CutSuffix(inst, ":"); found {
			if _, err := strconv.Atoi(jumpVal); err != nil {
				if a.traceEnabled {
					fmt.Printf("[buildSymbolTable:34] \n\ttoken: %v\n\tjumpVal: %v\n", token, jumpVal)
				}
				symbolTable[jumpVal] = pos
			}
		} else {
			pos++
		}
	}
	if a.traceEnabled {
		fmt.Printf("[buildSymbolTable:43] %v\n", symbolTable)
	}
	return symbolTable
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

func SplitRemoveEmpty(value, sep string) []string {
	parts := strings.Split(value, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return out
}
func (a *Assembler) resolveJumpLabelsIfPresent(byteCode []int, opCode int, op string, labels map[string]int) []int {
	byteCode = append(byteCode, opCode)
	label := op
	if _, err := strconv.Atoi(label); err != nil {
		byteCode = append(byteCode, int(labels[label]))
		if a.traceEnabled {
			fmt.Printf("[Assemble:115] \n\tbc:%v\n\tval:%v\n\tlabel:%s\n\tlabels:%+v\n\terr: %v\n", byteCode, labels[label], label, labels, err)
		}
	}
	if a.traceEnabled {
		fmt.Printf("[Assemble:119] bc(%v):val/label(%v)\n", byteCode, label)
	}
	return byteCode
}

func (a *Assembler) Assemble(input string) ([]int, error) {
	byteCode := []int{}

	lines := SplitRemoveEmpty(a.removeComments(input), "\n")

	if a.traceEnabled {
		for i, line := range lines {
			fmt.Printf("[Assemble:89] %d %s\n", i, line)
		}
	}

	labels := a.buildSymbolTable(lines)

	for _, line := range lines {

		op := SplitRemoveEmpty(line, " ")
		for i, inst := range op {
			if a.traceEnabled {
				fmt.Printf("[Assemble:101] %d:%v\n", i, inst)
				fmt.Printf("[Assemble:102] bytecode: %+v\n", byteCode)
			}
			switch inst {
			case "push":
				byteCode = append(byteCode, int(stack.PSH))
			case "jmp":
				byteCode = a.resolveJumpLabelsIfPresent(byteCode, int(stack.JMP), op[i+1], labels)
			case "jz":
				byteCode = a.resolveJumpLabelsIfPresent(byteCode, int(stack.JZ), op[i+1], labels)
			case "jnz":
				byteCode = a.resolveJumpLabelsIfPresent(byteCode, int(stack.JNZ), op[i+1], labels)
			case "add":
				byteCode = append(byteCode, int(stack.ADD))
			case "sub":
				byteCode = append(byteCode, int(stack.SUB))
			case "div":
				byteCode = append(byteCode, int(stack.DIV))
			case "mul":
				byteCode = append(byteCode, int(stack.MUL))
			case "dup":
				byteCode = append(byteCode, int(stack.DUP))
			case "swap":
				byteCode = append(byteCode, int(stack.SWP))
			case "drop":
				byteCode = append(byteCode, int(stack.DRP))
			case "eq":
				byteCode = append(byteCode, int(stack.EQ))
			case "lt":
				byteCode = append(byteCode, int(stack.LT))
			case "lte":
				byteCode = append(byteCode, int(stack.LTE))
			case "gt":
				byteCode = append(byteCode, int(stack.GT))
			case "gte":
				byteCode = append(byteCode, int(stack.GTE))
			default:
				if a.traceEnabled {
					fmt.Printf("[Assemble:152] token: %v\n", inst)
				}
				if strings.HasSuffix(inst, ":") {
					// Skip label definitions - they're already in the symbol table
					continue
				}
				val, err := strconv.Atoi(inst)
				if err != nil {
					if a.traceEnabled {
						fmt.Printf("[Assemble:161] %v\n", err)
					}
					if _, found := labels[inst]; found {
						continue
					}
					return nil, err
				}
				byteCode = append(byteCode, val)
			}
		}
	}
	return byteCode, nil
}
