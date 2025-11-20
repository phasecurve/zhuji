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
	BEQ
	BLT
	LTE
	GT
	BGE
	JMP
	JZ
	JNZ
	ADDI
	LW
	SW
	BNE
	MOD
)
const MVQ = -1
