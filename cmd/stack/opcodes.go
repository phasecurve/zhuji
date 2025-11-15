package stack

type OpCode int

const (
	PUSH OpCode = iota
	ADD
	SUB
	MUL
	DIV
)
