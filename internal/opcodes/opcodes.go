package opcodes

type OpCode int

const (
	MVQ  OpCode = -1
	PSH  OpCode = 0
	ADD  OpCode = 1
	SUB  OpCode = 2
	MUL  OpCode = 3
	DIV  OpCode = 4
	DUP  OpCode = 5
	SWP  OpCode = 6
	DRP  OpCode = 7
	BEQ  OpCode = 8
	BLT  OpCode = 9
	LTE  OpCode = 10
	GT   OpCode = 11
	BGE  OpCode = 12
	JMP  OpCode = 13
	JZ   OpCode = 14
	JNZ  OpCode = 15
	ADDI OpCode = 16
	LW   OpCode = 17
	SW   OpCode = 18
	BNE  OpCode = 19
	MOD  OpCode = 20
	JAL  OpCode = 21
	JALR OpCode = 22
)
