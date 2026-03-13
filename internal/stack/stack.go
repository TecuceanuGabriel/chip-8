package stack

import "fmt"

// Stack is a generic LIFO stack backed by a slice.
type Stack[T any] struct {
	items []T
}

// Push appends item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Peek returns the top item without removing it, or an error if the stack is empty.
func (s *Stack[T]) Peek() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, fmt.Errorf("Peek() failed: stack is empty")
	}

	return s.items[len(s.items)-1], nil
}

// Pop removes and returns the top item, or an error if the stack is empty.
func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, fmt.Errorf("Pop() failed: stack is empty")
	}

	top, _ := s.Peek()
	s.items = s.items[:len(s.items)-1]
	return top, nil
}
