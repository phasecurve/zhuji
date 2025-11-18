// Package memory is the conceptualisation of memory the hardware would have
package memory

type Memory struct {
	data []byte
}

func NewMemory(size int) Memory {
	return Memory{
		data: make([]byte, size),
	}
}

func (m *Memory) LoadWord(address int) int32 {
	b0 := int32(m.data[address+0])
	b1 := int32(m.data[address+1])
	b2 := int32(m.data[address+2])
	b3 := int32(m.data[address+3])
	return b0 | (b1 << 8) | (b2 << 16) | (b3 << 24)
}

func (m *Memory) StoreWord(address int, value int32) {
	m.data[address] = byte(value)
	m.data[address+1] = byte(value >> 8)
	m.data[address+2] = byte(value >> 16)
	m.data[address+3] = byte(value >> 24)
}

func (m *Memory) LoadByte(address int) byte {
	return byte(m.data[address])
}

func (m *Memory) StoreByte(address int, value byte) {
	m.data[address] = value
}
