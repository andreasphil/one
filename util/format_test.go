package util_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/andreasphil/one/util"
)

func TestFormatNestedMap(t *testing.T) {
	type testcase struct {
		name     string
		input    map[string][]string
		expected string
	}

	testcases := []testcase{
		{
			name:     "empty map",
			input:    map[string][]string{},
			expected: "",
		},
		{
			name:     "nil map",
			input:    nil,
			expected: "",
		},
		{
			name:     "single key with single value",
			input:    map[string][]string{"one": {"1"}},
			expected: "one=1",
		},
		{
			name:     "single key with multiple values",
			input:    map[string][]string{"one": {"1", "2"}},
			expected: "one=1 one=2",
		},
		{
			name:     "key without values",
			input:    map[string][]string{"one": {}},
			expected: "",
		},
		{
			name:     "empty key and value",
			input:    map[string][]string{"": {""}},
			expected: "=",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := util.FormatNestedMap(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestFormatNestedMapMultipleKeys(t *testing.T) {
	input := map[string][]string{
		"one": {"1"},
		"two": {"2", "3"},
	}

	result := util.FormatNestedMap(input)

	// Map iteration order is not stable, so compare the pairs instead of the
	// whole string.
	pairs := strings.Fields(result)
	slices.Sort(pairs)
	expected := []string{"one=1", "two=2", "two=3"}

	if !slices.Equal(pairs, expected) {
		t.Errorf("expected %v, got %v", expected, pairs)
	}
}
