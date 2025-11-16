// Package stack ...
package stack

import (
	"fmt"
)

type Stack struct {
	stack []int
}

func NewStack() *Stack {
	return &Stack{
		stack: make([]int, 0),
	}
}

func (s *Stack) IsEmpty() bool {
	return len(s.stack) == 0
}

func (s *Stack) Push(val int) {
	s.stack = append(s.stack, val)
}

func (s *Stack) Pop() int {
	popped := s.stack[len(s.stack)-1]
	s.stack = s.stack[0 : len(s.stack)-1]
	return popped
}

func (s *Stack) Peek() int {
	return s.stack[len(s.stack)-1]
}

func (s *Stack) Swap() {
	first := s.Pop()
	second := s.Pop()
	s.Push(first)
	s.Push(second)
}

func (s *Stack) Drop() {
	s.Pop()
}

func (s *Stack) Dup() {
	s.Push(s.Peek())
}

func (s *Stack) Size() int {
	return len(s.stack)
}

func (s *Stack) String() string {
	if len(s.stack) == 0 {
		return "Stack: []"
	}
	if len(s.stack) == 1 {
		return fmt.Sprintf("Stack: [%d]", s.stack[0])
	}

	top := s.Peek()
	stackStr := "Stack: ["
	for i := 0; i < len(s.stack); {
		if i+1 == len(s.stack) {
			stackStr = fmt.Sprintf("%s%d]", stackStr, s.stack[i])
		} else {
			stackStr = fmt.Sprintf("%s%d, ", stackStr, s.stack[i])
		}
		i++
	}
	return fmt.Sprintf("%s (top: %d)", stackStr, top)
}
