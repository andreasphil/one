package lib_test

import (
	"testing"
	"time"

	"github.com/andreasphil/one/lib"
	"github.com/google/go-cmp/cmp"
)

func TestDuplicateSlugs(t *testing.T) {
	type testcase struct {
		name     string
		notes    []lib.Note
		expected []string
	}

	testcases := []testcase{
		{
			name: "returns nil if there are no duplicates",
			notes: []lib.Note{
				{Title: "Root 1"},
				{Title: "Root 2"},
			},
			expected: nil,
		},
		{
			name: "returns duplicate slugs at the root level",
			notes: []lib.Note{
				{Title: "Root 1"},
				{Title: "Root 1"},
				{Title: "Root 2"},
			},
			expected: []string{"root-1"},
		},
		{
			name: "returns duplicate slugs across nested children",
			notes: []lib.Note{
				{
					Title: "Root 1",
					Children: []lib.Note{
						{Title: "Child 1"},
					},
				},
				{Title: "Child 1"},
			},
			expected: []string{"child-1"},
		},
		{
			name: "returns each duplicate slug only once, in first-seen order",
			notes: []lib.Note{
				{Title: "B"},
				{Title: "A"},
				{Title: "B"},
				{Title: "A"},
				{Title: "B"},
			},
			expected: []string{"b", "a"},
		},
		{
			name:     "returns nil for empty input",
			notes:    []lib.Note{},
			expected: nil,
		},
		{
			name: "treats daily notes with the same date as duplicates",
			notes: []lib.Note{
				{
					Title: "01.01.2025",
					Date:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Title: "01.01.2025",
					Date:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: []string{"2025-01-01"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := lib.DuplicateSlugs(tc.notes)
			if !cmp.Equal(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestEmptyTitles(t *testing.T) {
	type testcase struct {
		name     string
		notes    []lib.Note
		expected int
	}

	testcases := []testcase{
		{
			name: "returns 0 if there are no empty titles",
			notes: []lib.Note{
				{Title: "Root 1"},
				{Title: "Root 2"},
			},
			expected: 0,
		},
		{
			name: "counts notes with empty titles at the root level",
			notes: []lib.Note{
				{Title: ""},
				{Title: "Root 2"},
			},
			expected: 1,
		},
		{
			name: "counts notes with empty titles in nested children",
			notes: []lib.Note{
				{
					Title: "Root 1",
					Children: []lib.Note{
						{Title: ""},
						{
							Title: "Child 2",
							Children: []lib.Note{
								{Title: ""},
							},
						},
					},
				},
			},
			expected: 2,
		},
		{
			name:     "returns 0 for empty input",
			notes:    []lib.Note{},
			expected: 0,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := lib.EmptyTitles(tc.notes)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestEmptyNotes(t *testing.T) {
	type testcase struct {
		name     string
		notes    []lib.Note
		expected []string
	}

	testcases := []testcase{
		{
			name: "returns nil if there are no empty notes",
			notes: []lib.Note{
				{Title: "Root 1", Raw: "# Root 1\n\nContent"},
			},
			expected: nil,
		},
		{
			name: "returns titles of notes without content at the root level",
			notes: []lib.Note{
				{Title: "Root 1", Raw: "# Root 1\n\n"},
				{Title: "Root 2", Raw: "# Root 2\n\nContent"},
			},
			expected: []string{"Root 1"},
		},
		{
			name: "returns titles of notes without content in nested children",
			notes: []lib.Note{
				{
					Title: "Root 1",
					Raw:   "# Root 1\n\n",
					Children: []lib.Note{
						{Title: "Child 1", Raw: "## Child 1\n\n"},
						{Title: "Child 2", Raw: "## Child 2\n\nContent"},
					},
				},
			},
			expected: []string{"Child 1"},
		},
		{
			name: "does not count notes with children as empty, even if their own content is empty",
			notes: []lib.Note{
				{
					Title: "Root 1",
					Raw:   "# Root 1\n\n",
					Children: []lib.Note{
						{Title: "Child 1", Raw: "## Child 1\n\nContent"},
					},
				},
			},
			expected: nil,
		},
		{
			name:     "returns nil for empty input",
			notes:    []lib.Note{},
			expected: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := lib.EmptyNotes(tc.notes)
			if !cmp.Equal(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
