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
			if got := util.FormatNestedMap(tc.input); got != tc.expected {
				t.Errorf("FormatNestedMap(%v) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestFormatNestedMapMultipleKeys(t *testing.T) {
	input := map[string][]string{
		"one": {"1"},
		"two": {"2", "3"},
	}

	// Map iteration order is not stable, so compare the pairs instead of the
	// whole string.
	got := strings.Fields(util.FormatNestedMap(input))
	slices.Sort(got)
	want := []string{"one=1", "two=2", "two=3"}

	if !slices.Equal(got, want) {
		t.Errorf("FormatNestedMap(%v) pairs = %v, want %v", input, got, want)
	}
}
