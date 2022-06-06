package quang

import "sync"

type Stack[T any] struct {
	items []T
	lock  sync.RWMutex
}

func NewStack[T any]() *Stack[T] {
	s := []T{}

	return &Stack[T]{
		items: s,
	}
}

func (s *Stack[T]) Count() int {
	return len(s.items)
}

func (s *Stack[T]) Push(item T) {
	s.lock.Lock()
	s.items = append(s.items, item)
	s.lock.Unlock()
}

func (s *Stack[T]) Pop() (item *T) {
	if len(s.items) == 0 {
		return
	}

	s.lock.Lock()
	item = &s.items[len(s.items)-1]
	s.items = s.items[0 : len(s.items)-1]
	s.lock.Unlock()

	return item
}

func (s *Stack[T]) Peek() (item *T) {
	if len(s.items) == 0 {
		return
	}

	s.lock.Lock()
	item = &s.items[len(s.items)-1]
	s.lock.Unlock()

	return item
}
