package stack

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
			opCode:   PUSH,
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, OpCode(tt.expected), tt.opCode)
		})
	}
}
