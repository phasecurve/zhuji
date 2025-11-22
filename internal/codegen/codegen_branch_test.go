package codegen

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/opcodes"
	"github.com/stretchr/testify/assert"
)

func TestBEQ(t *testing.T) {
	cases := []struct {
		name          string
		bytecode      []int
		shouldContain []string
		message       string
	}{
		{
			"taken",
			[]int{
				int(opcodes.ADDI), 1, 0, 5,
				int(opcodes.ADDI), 2, 0, 5,
				int(opcodes.BEQ), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
			},
			[]string{"cmpq %rbx, %rax", "je L16", "L16:"},
			"equal values should generate jump to target label",
		},
		{
			"not taken",
			[]int{
				int(opcodes.ADDI), 1, 0, 5,
				int(opcodes.ADDI), 2, 0, 10,
				int(opcodes.BEQ), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
			},
			[]string{"cmpq %rbx, %rax", "je L16", "movq $99, %rax"},
			"different values should generate code for both paths",
		},
		{
			"mid program label",
			[]int{
				int(opcodes.ADDI), 1, 0, 5,
				int(opcodes.ADDI), 2, 0, 5,
				int(opcodes.BEQ), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
				int(opcodes.ADDI), 3, 0, 42,
			},
			[]string{"je L16", "L16:", "movq $42, %rcx"},
			"branch target in middle should place label before subsequent instruction",
		},
		{
			"backward jump",
			[]int{
				int(opcodes.ADDI), 1, 0, 5,
				int(opcodes.ADDI), 2, 0, 10,
				int(opcodes.BEQ), 1, 2, -8,
			},
			[]string{"je L0", "L0:"},
			"negative offset should generate backward jump to earlier label",
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

func TestBLT(t *testing.T) {
	cases := []struct {
		name          string
		bytecode      []int
		shouldContain []string
		message       string
	}{
		{
			"taken",
			[]int{
				int(opcodes.ADDI), 1, 0, 5,
				int(opcodes.ADDI), 2, 0, 10,
				int(opcodes.BLT), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
			},
			[]string{"cmpq %rbx, %rax", "jl L16", "L16:"},
			"smaller first operand should generate jump to target label",
		},
		{
			"not taken",
			[]int{
				int(opcodes.ADDI), 1, 0, 10,
				int(opcodes.ADDI), 2, 0, 5,
				int(opcodes.BLT), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
			},
			[]string{"cmpq %rbx, %rax", "jl L16", "movq $99, %rax"},
			"larger first operand should generate code for both paths",
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

func TestBNE(t *testing.T) {
	cases := []struct {
		name          string
		bytecode      []int
		shouldContain []string
		message       string
	}{
		{
			"taken",
			[]int{
				int(opcodes.ADDI), 1, 0, 5,
				int(opcodes.ADDI), 2, 0, 10,
				int(opcodes.BNE), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
			},
			[]string{"cmpq %rbx, %rax", "jne L16", "L16:"},
			"different values should generate jump to target label",
		},
		{
			"not taken",
			[]int{
				int(opcodes.ADDI), 1, 0, 5,
				int(opcodes.ADDI), 2, 0, 5,
				int(opcodes.BNE), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
			},
			[]string{"cmpq %rbx, %rax", "jne L16", "movq $99, %rax"},
			"equal values should generate code for both paths",
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

func TestBGE(t *testing.T) {
	cases := []struct {
		name          string
		bytecode      []int
		shouldContain []string
		message       string
	}{
		{
			"taken",
			[]int{
				int(opcodes.ADDI), 1, 0, 10,
				int(opcodes.ADDI), 2, 0, 5,
				int(opcodes.BGE), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
			},
			[]string{"cmpq %rbx, %rax", "jge L16", "L16:"},
			"larger first operand should generate jump to target label",
		},
		{
			"not taken",
			[]int{
				int(opcodes.ADDI), 1, 0, 5,
				int(opcodes.ADDI), 2, 0, 10,
				int(opcodes.BGE), 1, 2, 8,
				int(opcodes.ADDI), 1, 0, 99,
			},
			[]string{"cmpq %rbx, %rax", "jge L16", "movq $99, %rax"},
			"smaller first operand should generate code for both paths",
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
