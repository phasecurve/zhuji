package codegen

import (
	"testing"

	"github.com/phasecurve/zhuji/internal/opcodes"
)

func TestEndToEndSimple(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
	}
	runEndToEnd(t, bytecode, 42, "simple immediate value should pass through to exit code")
}

func TestEndToEndAddDestinationEqualsSource2(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.ADD), 1, 2, 1,
	}
	runEndToEnd(t, bytecode, 15, "add with swapped operand order should compute correctly")
}

func TestEndToEndMul(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 6,
		int(opcodes.ADDI), 2, 0, 7,
		int(opcodes.MUL), 1, 1, 2,
	}
	runEndToEnd(t, bytecode, 42, "multiplication should produce correct result")
}

func TestEndToEndDivDifferentDestination(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.ADDI), 2, 0, 6,
		int(opcodes.DIV), 3, 1, 2,
		int(opcodes.ADD), 1, 3, 0,
	}
	runEndToEnd(t, bytecode, 7, "division result should be movable to different register")
}

func TestEndToEndDivDividendNotInRax(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 6,
		int(opcodes.ADDI), 2, 0, 42,
		int(opcodes.DIV), 3, 2, 1,
		int(opcodes.ADD), 1, 3, 0,
	}
	runEndToEnd(t, bytecode, 7, "division should work when dividend is not in rax")
}

func TestEndToEndBEQTaken(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BEQ), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}
	runEndToEnd(t, bytecode, 5, "branch on equal should skip instruction when values match")
}

func TestEndToEndBLTTaken(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BLT), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}
	runEndToEnd(t, bytecode, 5, "branch on less than should skip instruction when first < second")
}

func TestEndToEndBNETaken(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 5,
		int(opcodes.ADDI), 2, 0, 10,
		int(opcodes.BNE), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}
	runEndToEnd(t, bytecode, 5, "branch on not equal should skip instruction when values differ")
}

func TestEndToEndBGETaken(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 10,
		int(opcodes.ADDI), 2, 0, 5,
		int(opcodes.BGE), 1, 2, 8,
		int(opcodes.ADDI), 1, 0, 99,
	}
	runEndToEnd(t, bytecode, 10, "branch on greater or equal should skip instruction when first >= second")
}

func TestEndToEndStoreLoad(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 42,
		int(opcodes.SW), 1, 0, 0,
		int(opcodes.LW), 2, 0, 0,
		int(opcodes.ADD), 1, 2, 0,
	}
	runEndToEnd(t, bytecode, 42, "value stored to memory should be loadable into different register")
}

func TestEndToEndFibonacci(t *testing.T) {
	bytecode := []int{
		int(opcodes.ADDI), 1, 0, 0,
		int(opcodes.ADDI), 2, 0, 1,
		int(opcodes.ADDI), 3, 0, 9,
		int(opcodes.ADDI), 4, 0, 0,
		int(opcodes.ADD), 5, 1, 2,
		int(opcodes.ADD), 1, 2, 0,
		int(opcodes.ADD), 2, 5, 0,
		int(opcodes.ADDI), 4, 4, 1,
		int(opcodes.BLT), 4, 3, -16,
		int(opcodes.ADD), 1, 2, 0,
	}
	runEndToEnd(t, bytecode, 55, "fibonacci loop should compute 10th fibonacci number")
}
