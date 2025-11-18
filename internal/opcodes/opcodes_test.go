package opcodes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpCodesAreDefined(t *testing.T) {
	tests := []struct {
		name     string
		opCode   OpCode
		expected int
	}{
		{
			name:     "OpCodePushIs_0",
			opCode:   PSH,
			expected: 0,
		}, {
			name:     "OpCodeAddIs_1",
			opCode:   ADD,
			expected: 1,
		}, {
			name:     "OpCodeSubIs_2",
			opCode:   SUB,
			expected: 2,
		}, {
			name:     "OpCodeMulIs_3",
			opCode:   MUL,
			expected: 3,
		}, {
			name:     "OpCodeDivIs_4",
			opCode:   DIV,
			expected: 4,
		}, {
			name:     "OpCodeDupIs_5",
			opCode:   DUP,
			expected: 5,
		}, {
			name:     "OpCodeSwpIs_6",
			opCode:   SWP,
			expected: 6,
		}, {
			name:     "OpCodeDrpIs_7",
			opCode:   DRP,
			expected: 7,
		}, {
			name:     "OpCodeEqIs_8",
			opCode:   EQ,
			expected: 8,
		}, {
			name:     "OpCodeLtIs_9",
			opCode:   LT,
			expected: 9,
		}, {
			name:     "OpCodeLteIs_10",
			opCode:   LTE,
			expected: 10,
		}, {
			name:     "OpCodeGtIs_11",
			opCode:   GT,
			expected: 11,
		}, {
			name:     "OpCodeGteIs_12",
			opCode:   GTE,
			expected: 12,
		}, {
			name:     "OpCodeJmpIs_13",
			opCode:   JMP,
			expected: 13,
		}, {
			name:     "OpCodeJzIs_14",
			opCode:   JZ,
			expected: 14,
		}, {
			name:     "OpCodeJnzIs_15",
			opCode:   JNZ,
			expected: 15,
		}, {
			name:     "OpCodeAddiIs_16",
			opCode:   ADDI,
			expected: 16,
		}, {
			name:     "OpCodeLwIs_17",
			opCode:   LW,
			expected: 17,
		}, {
			name:     "OpCodeSwIs_18",
			opCode:   SW,
			expected: 18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, OpCode(tt.expected), tt.opCode)
		})
	}
}
