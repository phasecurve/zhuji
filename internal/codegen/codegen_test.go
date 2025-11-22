package codegen

import (
	"bytes"
	"os"
	"testing"

	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/stretchr/testify/assert"
)

func TestAddi(t *testing.T) {
	cases := []struct {
		name          string
		bytecode      []int
		shouldContain []string
		message       string
	}{
		{
			"simple immediate",
			[]int{int(opcodes.ADDI), 1, 0, 42},
			[]string{"movq $42, %rax", "movq %rax, %rdi", "movq $60, %rax", "syscall"},
			"immediate to register should generate mov and exit syscall",
		},
		{
			"different immediate",
			[]int{int(opcodes.ADDI), 1, 0, 100},
			[]string{"movq $100, %rax"},
			"different immediate value should appear in generated code",
		},
		{
			"two instructions",
			[]int{
				int(opcodes.ADDI), 1, 0, 10,
				int(opcodes.ADDI), 1, 1, 5,
			},
			[]string{"movq $10, %rax", "addq $5, %rax"},
			"adding to existing register value should use addq instruction",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cg := NewCodeGen()
			asm := cg.Generate(tc.bytecode)
			for _, expected := range tc.shouldContain {
				assert.Contains(t, asm, expected, tc.message)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	cases := []struct {
		name          string
		bytecode      []int
		shouldContain []string
		message       string
	}{
		{
			"registers",
			[]int{
				int(opcodes.ADDI), 1, 0, 10,
				int(opcodes.ADDI), 2, 0, 5,
				int(opcodes.ADD), 1, 1, 2,
			},
			[]string{"movq $10, %rax", "movq $5, %rbx", "addq %rbx, %rax"},
			"adding two registers should generate addq with register operands",
		},
		{
			"different destination",
			[]int{
				int(opcodes.ADDI), 1, 0, 10,
				int(opcodes.ADDI), 2, 0, 5,
				int(opcodes.ADD), 3, 1, 2,
			},
			[]string{"movq %rax, %rcx", "addq %rbx, %rcx"},
			"adding to different destination should move result to target register",
		},
		{
			"x0 source",
			[]int{
				int(opcodes.ADDI), 1, 0, 42,
				int(opcodes.ADD), 2, 1, 0,
			},
			[]string{"movq %rax, %rbx", "addq $0, %rbx"},
			"x0 source should generate add with zero immediate",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cg := NewCodeGen()
			asm := cg.Generate(tc.bytecode)
			for _, expected := range tc.shouldContain {
				assert.Contains(t, asm, expected, tc.message)
			}
		})
	}
}

func TestSub(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 3,
		int(opcodes.SUB), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $10, %rax", "first operand should be in rax")
	assert.Contains(t, asm, "movq $3, %rbx", "second operand should be in rbx")
	assert.Contains(t, asm, "subq %rbx, %rax", "subtraction should use subq instruction")
}

func TestMul(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 6,
		int(opcodes.ADDI), 2, 0, 7,
		int(opcodes.MUL), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $6, %rax", "first operand should be in rax")
	assert.Contains(t, asm, "movq $7, %rbx", "second operand should be in rbx")
	assert.Contains(t, asm, "imulq %rbx, %rax", "multiplication should use imulq instruction")
}

func TestDiv(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.ADDI), 2, 0, 6,
		int(opcodes.DIV), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $42, %rax", "dividend should be in rax")
	assert.Contains(t, asm, "movq $6, %rbx", "divisor should be in rbx")
	assert.Contains(t, asm, "cqto", "division should sign-extend rax to rdx:rax")
	assert.Contains(t, asm, "idivq %rbx", "division should use idivq instruction")
}

func TestMod(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 17,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.MOD), 1, 1, 2,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, "movq $17, %rax", "dividend should be in rax")
	assert.Contains(t, asm, "movq $5, %rbx", "divisor should be in rbx")
	assert.Contains(t, asm, "cqto", "modulo should sign-extend rax to rdx:rax")
	assert.Contains(t, asm, "idivq %rbx", "modulo should use idivq instruction")
	assert.Contains(t, asm, "movq %rdx, %rax", "modulo should move remainder from rdx to result register")
}

func TestHasEntryPoint(t *testing.T) {
	cg := NewCodeGen()
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
	}

	asm := cg.Generate(bytecode)

	assert.Contains(t, asm, ".global _start", "generated code should export _start symbol")
	assert.Contains(t, asm, "_start:", "generated code should have _start label")
}

func TestNoOutputByDefault(t *testing.T) {
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

	assert.Empty(t, output, "code generation should not produce stdout output")
}
