package stack

type OpCode int

const (
	PSH OpCode = iota
	ADD
	SUB
	MUL
	DIV
	DUP
	SWP
	DRP
	EQ
	LT
	LTE
	GT
	GTE
	JMP
	JZ
	JNZ
	SKP
)
