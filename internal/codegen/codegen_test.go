package codegen

import (
	"bytes"
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

func TestCodegenAddWithX0Source(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.ADD), 2, 1, 0,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq %rax, %rbx")
	assert.Contains(t, asm, "addq $0, %rbx")
}

func TestCodegenAddDestinationEqualsSource2EndToEnd(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.ADD), 1, 2, 1,
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
	assert.Equal(t, 15, exitCode)
}

func TestCodegenNoOutputByDefault(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cg.Generate(bytecode)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	assert.Empty(t, output)
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

func TestCodegenDivDifferentDestinationEndToEnd(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.ADDI), 2, 0, 6,
		int(opcodes.DIV), 3, 1, 2,
		int(opcodes.ADD), 1, 3, 0,
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
	assert.Equal(t, 7, exitCode)
}

func TestCodegenDivDividendNotInRaxEndToEnd(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 6,
		int(opcodes.ADDI), 2, 0, 42,
		int(opcodes.DIV), 3, 2, 1,
		int(opcodes.ADD), 1, 3, 0,
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
	assert.Equal(t, 7, exitCode)
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

func TestCodegenBEQTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BEQ), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "cmpq %rbx, %rax")
	assert.Contains(t, asm, "je L16")
	assert.Contains(t, asm, "L16:")
}

func TestCodegenBEQMidProgramLabel(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BEQ), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
		int(opcodes.ADDI), 3, 0, 42,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "je L16")
	assert.Contains(t, asm, "L16:")
	assert.Contains(t, asm, "movq $42, %rcx")
}

func TestCodegenBEQNotTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BEQ), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "cmpq %rbx, %rax")
	assert.Contains(t, asm, "je L16")
	assert.Contains(t, asm, "movq $99, %rax")
}

func TestCodegenBEQBackwardJump(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BEQ), 1, 2, -8,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "je L0")
	assert.Contains(t, asm, "L0:")
}

func TestCodegenBEQEndToEndTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BEQ), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
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
	assert.Equal(t, 5, exitCode)
}

func TestCodegenBLTTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BLT), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "cmpq %rbx, %rax")
	assert.Contains(t, asm, "jl L16")
	assert.Contains(t, asm, "L16:")
}

func TestCodegenBLTNotTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BLT), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "cmpq %rbx, %rax")
	assert.Contains(t, asm, "jl L16")
	assert.Contains(t, asm, "movq $99, %rax")
}

func TestCodegenBLTEndToEndTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BLT), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
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
	assert.Equal(t, 5, exitCode)
}

func TestCodegenBNETaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BNE), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "cmpq %rbx, %rax")
	assert.Contains(t, asm, "jne L16")
	assert.Contains(t, asm, "L16:")
}

func TestCodegenBNENotTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BNE), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "cmpq %rbx, %rax")
	assert.Contains(t, asm, "jne L16")
	assert.Contains(t, asm, "movq $99, %rax")
}

func TestCodegenBNEEndToEndTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BNE), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
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
	assert.Equal(t, 5, exitCode)
}

func TestCodegenBGETaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BGE), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "cmpq %rbx, %rax")
	assert.Contains(t, asm, "jge L16")
	assert.Contains(t, asm, "L16:")
}

func TestCodegenBGENotTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BGE), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "cmpq %rbx, %rax")
	assert.Contains(t, asm, "jge L16")
	assert.Contains(t, asm, "movq $99, %rax")
}

func TestCodegenBGEEndToEndTaken(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BGE), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
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
	assert.Equal(t, 10, exitCode)
}
