package compiler

import (
	"github.com/phasecurve/zhuji/internal/assembler"
	"github.com/phasecurve/zhuji/internal/codegen"
)

func Compile(riscvAsm string) string {
	asm := assembler.NewAssembler()
	bytecode := asm.Assemble(riscvAsm)

	gen := codegen.NewCodeGen()
	return gen.Generate(bytecode)
}
