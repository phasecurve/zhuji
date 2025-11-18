// Package registers conceptualises haw registers will work
package registers

type Registers struct {
	values [2]int32
}

func NewRegisters() *Registers {
	return &Registers{}
}

func (r *Registers) Read(register int) int32 {
	return r.values[register]
}

func (r *Registers) Write(register int, val int32) {
	if register != 0 {
		r.values[register] = val
	}
}
