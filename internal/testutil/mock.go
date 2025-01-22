package testutil

import (
	"context"
	"sync"
)

// MockSource implements a mock configuration source for testing
type MockSource[T any] struct {
	mu      sync.RWMutex
	current T
	ch      chan T
}

func NewMockSource[T any]() *MockSource[T] {
	return &MockSource[T]{
		ch: make(chan T, 1),
	}
}

func (s *MockSource[T]) Load(ctx context.Context) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current, nil
}

func (s *MockSource[T]) Watch(ctx context.Context) (<-chan T, error) {
	return s.ch, nil
}

func (s *MockSource[T]) Update(value T) {
	s.mu.Lock()
	s.current = value
	s.mu.Unlock()

	select {
	case s.ch <- value:
	default:
		// Channel is full or closed
	}
}

func (s *MockSource[T]) Close() {
	close(s.ch)
}

// MockValidator implements a mock validator for testing
type MockValidator[T any] struct {
	ValidateFunc func(old, new T) error
}

func (v *MockValidator[T]) Validate(old, new T) error {
	if v.ValidateFunc != nil {
		return v.ValidateFunc(old, new)
	}
	return nil
}
