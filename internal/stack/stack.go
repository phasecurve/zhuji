// Package stack ...
package stack

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
