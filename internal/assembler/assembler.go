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
				byteCode = append(byteCode, int(stack.PUSH))
			case "add":
				byteCode = append(byteCode, int(stack.ADD))
			case "sub":
				byteCode = append(byteCode, int(stack.SUB))
			case "div":
				byteCode = append(byteCode, int(stack.DIV))
			case "mul":
				byteCode = append(byteCode, int(stack.MUL))
			default:
				var val int
				var err error
				if val, err = strconv.Atoi(inst); err != nil {
					return nil, err
				}
				byteCode = append(byteCode, val)
			}
		}
	}
	return byteCode, nil
}
