package lib_test

import (
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
	"github.com/google/go-cmp/cmp"
)

func TestGetRecursive(t *testing.T) {
	type testcase struct {
		name     string
		notes    []lib.Note
		slug     string
		expectOk bool
	}

	notes := []lib.Note{
		{Title: "Root 1"},
		{
			Title: "Root 2",
			Children: []lib.Note{
				{Title: "Child 1"},
				{
					Title: "Child 2",
					Children: []lib.Note{
						{Title: "Grandchild 1"},
					},
				},
			},
		},
		{Title: "Root 3"},
	}

	testcases := []testcase{
		{
			name:     "find note at root level",
			notes:    notes,
			slug:     "root-1",
			expectOk: true,
		},
		{
			name:     "find note in children",
			notes:    notes,
			slug:     "child-1",
			expectOk: true,
		},
		{
			name:     "find note in nested children",
			notes:    notes,
			slug:     "grandchild-1",
			expectOk: true,
		},
		{
			name:     "note not found",
			notes:    notes,
			slug:     "nonexistent",
			expectOk: false,
		},
		{
			name:     "search in empty slice",
			notes:    []lib.Note{},
			slug:     "any",
			expectOk: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, ok := lib.GetRecursive(tc.notes, tc.slug)

			if ok != tc.expectOk {
				t.Errorf("expected ok=%v, got ok=%v", tc.expectOk, ok)
			}

			if tc.expectOk {
				if result.Slug() != tc.slug {
					t.Errorf("expected slug %q, got %q", tc.slug, result.Slug())
				}
			} else {
				exportSetInternals := cmp.AllowUnexported(util.NewSet[string]())
				if !cmp.Equal(result, lib.Note{}, exportSetInternals) {
					t.Errorf("expected empty Note{}, got %+v", result)
				}
			}
		})
	}
}

func TestWalk(t *testing.T) {
	notes := []lib.Note{
		{Title: "Root 1"},
		{
			Title: "Root 2",
			Children: []lib.Note{
				{Title: "Child 1"},
				{
					Title: "Child 2",
					Children: []lib.Note{
						{Title: "Grandchild 1"},
					},
				},
			},
		},
		{Title: "Root 3"},
	}

	t.Run("visits every note in depth-first, pre-order", func(t *testing.T) {
		var visited []string

		result := lib.Walk(notes, func(n lib.Note) bool {
			visited = append(visited, n.Title)
			return true
		})

		expected := []string{
			"Root 1", "Root 2", "Child 1", "Child 2", "Grandchild 1", "Root 3",
		}

		if !cmp.Equal(visited, expected) {
			t.Errorf("expected visited %v, got %v", expected, visited)
		}

		if !result {
			t.Errorf("expected result to be true, got false")
		}
	})

	t.Run("stops early when fn returns false", func(t *testing.T) {
		var visited []string

		result := lib.Walk(notes, func(n lib.Note) bool {
			visited = append(visited, n.Title)
			return n.Title != "Child 1"
		})

		expected := []string{"Root 1", "Root 2", "Child 1"}

		if !cmp.Equal(visited, expected) {
			t.Errorf("expected visited %v, got %v", expected, visited)
		}

		if result {
			t.Errorf("expected result to be false, got true")
		}
	})

	t.Run("stopping in nested children also stops parent traversal", func(t *testing.T) {
		var visited []string

		result := lib.Walk(notes, func(n lib.Note) bool {
			visited = append(visited, n.Title)
			return n.Title != "Grandchild 1"
		})

		expected := []string{
			"Root 1", "Root 2", "Child 1", "Child 2", "Grandchild 1",
		}

		if !cmp.Equal(visited, expected) {
			t.Errorf("expected visited %v, got %v", expected, visited)
		}

		if result {
			t.Errorf("expected result to be false, got true")
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		var visited []string

		result := lib.Walk([]lib.Note{}, func(n lib.Note) bool {
			visited = append(visited, n.Title)
			return true
		})

		if len(visited) != 0 {
			t.Errorf("expected no notes to be visited, got %v", visited)
		}

		if !result {
			t.Errorf("expected result to be true, got false")
		}
	})
}

func TestSort(t *testing.T) {
	type testcase struct {
		name          string
		notes         []lib.Note
		expected      []string
		expectDidSort bool
	}

	testcases := []testcase{
		{
			name: "sorts daily notes by date (descending)",
			notes: []lib.Note{
				{
					Title: "01.01.2025",
					Date:  time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
				},
				{
					Title: "01.03.2025",
					Date:  time.Date(2025, 03, 01, 0, 0, 0, 0, time.UTC),
				},
				{
					Title: "01.02.2025",
					Date:  time.Date(2025, 02, 01, 0, 0, 0, 0, time.UTC),
				},
			},
			expected:      []string{"01.03.2025", "01.02.2025", "01.01.2025"},
			expectDidSort: true,
		},
		{
			name: "sorts knowledge base notes alphabetically (ascending)",
			notes: []lib.Note{
				{
					Title: "B",
				},
				{
					Title: "C",
				},
				{
					Title: "A",
				},
			},
			expected:      []string{"A", "B", "C"},
			expectDidSort: true,
		},
		{
			name: "groups all daily notes before knowledge base notes",
			notes: []lib.Note{
				{
					Title: "01.01.2025",
					Date:  time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
				},
				{
					Title: "A",
				},
				{
					Title: "01.02.2025",
					Date:  time.Date(2025, 02, 01, 0, 0, 0, 0, time.UTC),
				},
			},
			expected:      []string{"01.02.2025", "01.01.2025", "A"},
			expectDidSort: true,
		},
		{
			name: "sort is not case-sensitive",
			notes: []lib.Note{
				{
					Title: "a",
				},
				{
					Title: "b",
				},
				{
					Title: "A",
				},
			},
			expected:      []string{"a", "A", "b"},
			expectDidSort: true,
		},
		{
			name: "does not change already sorted notes",
			notes: []lib.Note{
				{
					Title: "01.02.2025",
					Date:  time.Date(2025, 02, 01, 0, 0, 0, 0, time.UTC),
				},
				{
					Title: "01.01.2025",
					Date:  time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
				},
				{
					Title: "A",
				},
			},
			expected:      []string{"01.02.2025", "01.01.2025", "A"},
			expectDidSort: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, didSort := lib.Sort(tc.notes)
			if didSort != tc.expectDidSort {
				t.Errorf("expected sorted to be %v, got %v", tc.expectDidSort, didSort)
			}

			for i, note := range result {
				if note.Title != tc.expected[i] {
					t.Errorf("expected note at %d to be %v, got %v", i, tc.expected[i], note.Title)
					t.FailNow()
				}
			}
		})
	}
}

func TestSortNormalizesNewline(t *testing.T) {
	input := "# B\n\nLine 1\n\n# A\n\nLine 2\n"
	expected := "# A\n\nLine 2\n\n# B\n\nLine 1\n"

	notes, err := lib.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	notes, _ = lib.Sort(notes)

	result := lib.ToString(notes)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestToString(t *testing.T) {
	type testcase struct {
		name     string
		input    string
		expected string
	}

	testcases := []testcase{
		{
			name:     "returns empty string for empty notes",
			input:    "",
			expected: "",
		},
		{
			name:     "preserves spacing between top level notes",
			input:    "# Note 1\n\nLine 1\n\n# Note 2\n\nLine 2\n",
			expected: "# Note 1\n\nLine 1\n\n# Note 2\n\nLine 2\n",
		},
		{
			name:     "preserves child notes in daily note output",
			input:    "# 01.01.2026\n\nLine 1\n\n## Child Note 1\n\nLine 2\n\n# 02.01.2026\n\nLine 3\n",
			expected: "# 01.01.2026\n\nLine 1\n\n## Child Note 1\n\nLine 2\n\n# 02.01.2026\n\nLine 3\n",
		},
		{
			name:     "normalizes multiple trailing newlines to one",
			input:    "# Note 1\n\nLine 1\n\n",
			expected: "# Note 1\n\nLine 1\n",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			notes, err := lib.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			result := lib.ToString(notes)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
