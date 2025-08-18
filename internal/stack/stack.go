package stack

import "fmt"

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Peek() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, fmt.Errorf("Peek() failed: stack is empty")
	}

	return s.items[len(s.items)-1], nil
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, fmt.Errorf("Pop() failed: stack is empty")
	}

	top, _ := s.Peek()
	s.items = s.items[:len(s.items)-1]
	return top, nil
}
