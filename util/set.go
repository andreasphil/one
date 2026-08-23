package util

import (
	"maps"
	"slices"
)

// Set is an unordered collection of unique values.
type Set[T comparable] struct {
	values map[T]bool
}

// NewSet creates a new, empty Set.
func NewSet[T comparable]() Set[T] {
	set := Set[T]{values: make(map[T]bool)}
	return set
}

// NewSetFrom creates a new Set containing the unique values from initial.
func NewSetFrom[T comparable](initial []T) Set[T] {
	set := NewSet[T]()
	set.Add(initial...)
	return set
}

// Add inserts values into the set, ignoring any that are already present. It
// returns the number of values that were actually added.
func (s Set[T]) Add(values ...T) int {
	added := 0

	for _, value := range values {
		if !s.Has(value) {
			s.values[value] = true
			added += 1
		}
	}

	return added
}

// Delete removes values from the set, ignoring any that are not present. It
// returns the number of values that were actually deleted.
func (s Set[T]) Delete(values ...T) int {
	deleted := 0

	for _, value := range values {
		if s.Has(value) {
			delete(s.values, value)
			deleted += 1
		}
	}

	return deleted
}

// Has reports whether value is present in the set.
func (s Set[T]) Has(value T) bool {
	_, ok := s.values[value]
	return ok
}

// Values returns the set's values as a slice, in no particular order.
func (s Set[T]) Values() []T {
	values := maps.Keys(s.values)
	return slices.Collect(values)
}

// Len returns the number of values in the set.
func (s Set[T]) Len() int {
	return len(s.values)
}
