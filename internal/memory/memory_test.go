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

	cases := []struct {
		addr     int
		expected int32
	}{
		{0, 10},
		{4, 20},
		{8, 30},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, m.LoadWord(tc.addr), "word at address %d should be isolated from other addresses", tc.addr)
	}
}

func TestLittleEndianByteOrder(t *testing.T) {
	m := NewMemory(1024)

	m.StoreWord(0, 0x12345678)

	cases := []struct {
		addr     int
		expected byte
	}{
		{0, 0x78},
		{1, 0x56},
		{2, 0x34},
		{3, 0x12},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, m.LoadByte(tc.addr), "byte %d should contain least significant byte first (little-endian)", tc.addr)
	}
}

func TestStoreByteThenLoadIt(t *testing.T) {
	m := NewMemory(1024)

	m.StoreByte(0, 0xFF)
	m.StoreByte(1, 0xAA)

	assert.Equal(t, byte(0xFF), m.LoadByte(0), "stored byte should be retrievable at address 0")
	assert.Equal(t, byte(0xAA), m.LoadByte(1), "stored byte should be retrievable at address 1")
}
