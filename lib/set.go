package lib

import (
	"maps"
	"slices"
)

type Set[T comparable] struct {
	values map[T]bool
}

func NewSet[T comparable]() Set[T] {
	set := Set[T]{values: make(map[T]bool)}
	return set
}

func NewSetFrom[T comparable](initial []T) Set[T] {
	set := NewSet[T]()
	set.Add(initial...)
	return set
}

func (s *Set[T]) Add(values ...T) int {
	added := 0

	for _, value := range values {
		if !s.Has(value) {
			s.values[value] = true
			added += 1
		}
	}

	return added
}

func (s *Set[T]) Delete(values ...T) int {
	deleted := 0

	for _, value := range values {
		if s.Has(value) {
			delete(s.values, value)
			deleted += 1
		}
	}

	return deleted
}

func (s *Set[T]) Has(value T) bool {
	_, ok := s.values[value]
	return ok
}

func (s *Set[T]) ToSlice() []T {
	values := maps.Keys(s.values)
	return slices.Collect(values)
}

func (s *Set[T]) Len() int {
	return len(s.values)
}
