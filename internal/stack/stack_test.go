package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	var s Stack[int]

	s.Push(1)
	top, err := s.Peek()

	assert.Nil(t, err)
	assert.Equal(t, 1, top, "The value on the top of the stock is 1")
}

func TestPop(t *testing.T) {
	var s Stack[int]

	s.Push(1)
	top, err := s.Pop()

	assert.Nil(t, err)
	assert.Equal(t, 1, top, "The poped value is 1")
}

func TestEmpty(t *testing.T) {
	var s Stack[int]

	_, err := s.Peek()
	assert.NotNil(t, err)

	_, err = s.Pop()
	assert.NotNil(t, err)
}

func TestMultiplePushPos(t *testing.T) {
	var s Stack[int]

	s.Push(1)
	s.Push(2)
	s.Pop()
	s.Push(3)
	s.Pop()

	top, err := s.Pop()

	assert.Nil(t, err)
	assert.Equal(t, 1, top, "The poped value is 1")
}
