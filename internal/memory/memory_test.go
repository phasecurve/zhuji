package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryInitialisesToZero(t *testing.T) {
	m := NewMemory(1024)

	assert.Equal(t, int32(0), m.LoadWord(0))
}

func TestStoreWordThenLoadIt(t *testing.T) {
	m := NewMemory(1024)

	m.StoreWord(0, 42)

	assert.Equal(t, int32(42), m.LoadWord(0))
}

func TestStoreWordsAtDifferentAddresses(t *testing.T) {
	m := NewMemory(1024)

	m.StoreWord(0, 10)
	m.StoreWord(4, 20)
	m.StoreWord(8, 30)

	assert.Equal(t, int32(10), m.LoadWord(0))
	assert.Equal(t, int32(20), m.LoadWord(4))
	assert.Equal(t, int32(30), m.LoadWord(8))
}

func TestLittleEndianByteOrder(t *testing.T) {
	m := NewMemory(1024)

	m.StoreWord(0, 0x12345678)

	assert.Equal(t, byte(0x78), m.LoadByte(0))
	assert.Equal(t, byte(0x56), m.LoadByte(1))
	assert.Equal(t, byte(0x34), m.LoadByte(2))
	assert.Equal(t, byte(0x12), m.LoadByte(3))
}

func TestStoreByteThenLoadIt(t *testing.T) {
	m := NewMemory(1024)

	m.StoreByte(0, 0xFF)
	m.StoreByte(1, 0xAA)

	assert.Equal(t, byte(0xFF), m.LoadByte(0))
	assert.Equal(t, byte(0xAA), m.LoadByte(1))
}
