package compiler

import (
	"strings"
	"testing"
)

func TestCompileProducesX86Assembly(t *testing.T) {
	riscvAsm := "addi x1, x0, 42"

	result := Compile(riscvAsm)

	if !strings.Contains(result, ".global _start") {
		t.Error("expected x86-64 assembly to contain .global _start")
	}
	if !strings.Contains(result, "movq $42, %rax") {
		t.Errorf("expected x86-64 assembly to contain movq $42, %%rax, got:\n%s", result)
	}
}
