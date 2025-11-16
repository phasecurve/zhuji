package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackIsEmpty(t *testing.T) {
	s := NewStack()
	assert.True(t, s.IsEmpty())
}

func TestStackNotEmptyAfterPush(t *testing.T) {
	s := NewStack()
	s.Push(3025)
	assert.False(t, s.IsEmpty())
}

func TestStackIsEmptyAfterLastItemPopped(t *testing.T) {
	s := NewStack()
	s.Push(3025)
	s.Pop()
	assert.True(t, s.IsEmpty())
}

func TestStackIsLIFO(t *testing.T) {
	s := NewStack()
	s.Push(3025)
	s.Push(4025)
	s.Push(5025)
	assert.Equal(t, 5025, s.Pop(), "popped value should be 5025")
	assert.Equal(t, 4025, s.Pop(), "popped value should be 4025")
	assert.Equal(t, 3025, s.Pop(), "popped value should be 3025")
}

func TestStackCanPeek(t *testing.T) {
	s := NewStack()
	s.Push(3025)
	s.Push(4025)
	s.Push(5025)
	assert.Equal(t, 5025, s.Peek(), "peeked value should be 5025")
	assert.Equal(t, 5025, s.Pop(), "popped value should be 5025")
}

func TestStackSize(t *testing.T) {
	s := NewStack()
	assert.Equal(t, 0, s.Size())

	s.Push(10)
	assert.Equal(t, 1, s.Size())

	s.Push(20)
	assert.Equal(t, 2, s.Size())

	s.Pop()
	assert.Equal(t, 1, s.Size())

	s.Pop()
	assert.Equal(t, 0, s.Size())
}

func TestStackString(t *testing.T) {
	s := NewStack()
	assert.Equal(t, "Stack: []", s.String())

	s.Push(10)
	assert.Equal(t, "Stack: [10]", s.String())

	s.Push(20)
	s.Push(30)
	assert.Equal(t, "Stack: [10, 20, 30] (top: 30)", s.String())
}
