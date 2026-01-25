package util

import (
	"testing"
)

func TestFilter(t *testing.T) {
	t.Run("filters integers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		result := Filter(input, func(n int) bool {
			return n%2 == 0
		})

		expected := []int{2, 4, 6}
		if len(result) != len(expected) {
			t.Fatalf("expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("at index %d: expected %d, got %d", i, v, result[i])
			}
		}
	})

	t.Run("filters strings", func(t *testing.T) {
		input := []string{"apple", "banana", "apricot", "cherry"}
		result := Filter(input, func(s string) bool {
			return len(s) > 5
		})

		expected := []string{"banana", "apricot", "cherry"}
		if len(result) != len(expected) {
			t.Fatalf("expected length %d, got %d", len(expected), len(result))
		}

		for i, v := range expected {
			if result[i] != v {
				t.Errorf("at index %d: expected %s, got %s", i, v, result[i])
			}
		}
	})

	t.Run("returns empty slice when no matches", func(t *testing.T) {
		input := []int{1, 3, 5, 7}
		result := Filter(input, func(n int) bool {
			return n%2 == 0
		})

		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("returns all items when all match", func(t *testing.T) {
		input := []int{2, 4, 6, 8}
		result := Filter(input, func(n int) bool {
			return n%2 == 0
		})

		if len(result) != len(input) {
			t.Fatalf("expected length %d, got %d", len(input), len(result))
		}

		for i, v := range input {
			if result[i] != v {
				t.Errorf("at index %d: expected %d, got %d", i, v, result[i])
			}
		}
	})

	t.Run("handles empty input slice", func(t *testing.T) {
		input := []int{}
		result := Filter(input, func(n int) bool {
			return n%2 == 0
		})

		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("handles nil slice", func(t *testing.T) {
		var input []int
		result := Filter(input, func(n int) bool {
			return n%2 == 0
		})

		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})
}
