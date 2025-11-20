package codegen

import (
	"os"
	"os/exec"
	"testing"

	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/stretchr/testify/assert"
)

func TestCodegenSimpleAddi(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $42, %rax")
	assert.Contains(t, asm, "movq %rax, %rdi")
	assert.Contains(t, asm, "movq $60, %rax")
	assert.Contains(t, asm, "syscall")
}

func TestCodegenAddiDifferentImmediate(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 100,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $100, %rax")
}

func TestCodegenHasEntryPoint(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, ".global _start")
	assert.Contains(t, asm, "_start:")
}

func TestCodegenTwoAddiInstructions(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 1, 1, 5,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $10, %rax")
	assert.Contains(t, asm, "addq $5, %rax")
}

func TestCodegenAddRegisters(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.ADD), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $10, %rax")
	assert.Contains(t, asm, "movq $5, %rbx")
	assert.Contains(t, asm, "addq %rbx, %rax")
}

func TestCodegenAddDifferentDestination(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.ADD), 3, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq %rax, %rcx")
	assert.Contains(t, asm, "addq %rbx, %rcx")
}

func TestCodegenEndToEndSimple(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
	}

	asm := cg.Generate(bytecode)

	tmpDir := t.TempDir()
	asmFile := tmpDir + "/test.s"
	objFile := tmpDir + "/test.o"
	exeFile := tmpDir + "/test"

	err := os.WriteFile(asmFile, []byte(asm), 0644)
	assert.NoError(t, err)

	cmd := exec.Command("as", "-o", objFile, asmFile)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "assembler failed: %s", string(output))

	cmd = exec.Command("ld", "-o", exeFile, objFile)
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "linker failed: %s", string(output))

	cmd = exec.Command(exeFile)
	cmd.Run()

	exitCode := cmd.ProcessState.ExitCode()
	assert.Equal(t, 42, exitCode)
}

func TestCodegenSubRegisters(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 3,
		int(opcodes.SUB), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $10, %rax")
	assert.Contains(t, asm, "movq $3, %rbx")
	assert.Contains(t, asm, "subq %rbx, %rax")
}

func TestCodegenMulRegisters(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 6,
		int(opcodes.ADDI), 2, 0, 7,
		int(opcodes.MUL), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $6, %rax")
	assert.Contains(t, asm, "movq $7, %rbx")
	assert.Contains(t, asm, "imulq %rbx, %rax")
}

func TestCodegenDivRegisters(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.ADDI), 2, 0, 6,
		int(opcodes.DIV), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $42, %rax")
	assert.Contains(t, asm, "movq $6, %rbx")
	assert.Contains(t, asm, "cqto")
	assert.Contains(t, asm, "idivq %rbx")
}

func TestCodegenModRegisters(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 17,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.MOD), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $17, %rax")
	assert.Contains(t, asm, "movq $5, %rbx")
	assert.Contains(t, asm, "cqto")
	assert.Contains(t, asm, "idivq %rbx")
	assert.Contains(t, asm, "movq %rdx, %rax")
}

func TestCodegenEndToEndMul(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 6,
		int(opcodes.ADDI), 2, 0, 7,
		int(opcodes.MUL), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	tmpDir := t.TempDir()
	asmFile := tmpDir + "/test.s"
	objFile := tmpDir + "/test.o"
	exeFile := tmpDir + "/test"

	err := os.WriteFile(asmFile, []byte(asm), 0644)
	assert.NoError(t, err)

	cmd := exec.Command("as", "-o", objFile, asmFile)
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "assembler failed: %s", string(output))

	cmd = exec.Command("ld", "-o", exeFile, objFile)
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "linker failed: %s", string(output))

	cmd = exec.Command(exeFile)
	cmd.Run()

	exitCode := cmd.ProcessState.ExitCode()
	assert.Equal(t, 42, exitCode)
}

// func TestCodegenBEQTaken(t *testing.T) {
// 	cg := NewCodeGen()
// 	bytecode := []int{
// 		int(opcodes.ADDI), 1, 0, 5,
// 		int(opcodes.ADDI), 2, 0, 5,
// 		int(opcodes.BEQ), 1, 2, 16,
// 		int(opcodes.ADDI), 1, 0, 99,
// 	}
//
// 	asm := cg.Generate(bytecode)
//
// 	assert.Contains(t, asm, "cmpq %rbx, %rax")
// 	assert.Contains(t, asm, "je L16")
// 	assert.Contains(t, asm, "L16:")
// }
