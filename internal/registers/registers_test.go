package registers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX0AlwaysReadsZero(t *testing.T) {
	r := NewRegisters()

	assert.Equal(t, int32(0), r.Read(0))
}

func TestWriteToX1ThenReadIt(t *testing.T) {
	r := NewRegisters()

	r.Write(1, 42)

	assert.Equal(t, int32(42), r.Read(1))
}

func TestAllRegistersInitialiseToZero(t *testing.T) {
	r := NewRegisters()

	for i := range 32 {
		assert.Equal(t, int32(0), r.Read(i), "register x%d", i)
	}
}

func TestWritesToX0AreIgnored(t *testing.T) {
	r := NewRegisters()

	r.Write(0, 42)

	assert.Equal(t, int32(0), r.Read(0))
}

func TestReadIsNonDestructive(t *testing.T) {
	r := NewRegisters()

	r.Write(1, 42)

	firstRead := r.Read(1)
	secondRead := r.Read(1)

	assert.Equal(t, int32(42), firstRead, "value should persist after first read")
	assert.Equal(t, int32(42), secondRead, "value should persist after second read")
}
