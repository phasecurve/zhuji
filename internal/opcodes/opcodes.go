package opcodes

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
	BLT
	LTE
	GT
	GTE
	JMP
	JZ
	JNZ
	ADDI
	LW
	SW
)
