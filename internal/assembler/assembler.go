// Package assembler will take risk-v/gas style assembler and conver to bytecode
package assembler

import (
	"strconv"
	"strings"

	"github.com/phasecurve/zhuji/internal/stack"
)

func Assemble(input string) ([]int, error) {
	byteCode := []int{}

	lines := strings.SplitSeq(input, "\n")

	for line := range lines {
		commentStart := strings.Index(line, "#")
		if commentStart == -1 {
			commentStart = len(line)
		}
		assemblerSplit := strings.FieldsSeq(line[:commentStart])

		for inst := range assemblerSplit {
			switch inst {
			case "push":
				byteCode = append(byteCode, int(stack.PSH))
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
				val, err := strconv.Atoi(inst)
				if err != nil {
					return nil, err
				}
				byteCode = append(byteCode, val)
			}
		}
	}
	return byteCode, nil
}
