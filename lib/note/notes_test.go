package note_test

import (
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestFindBySlug(t *testing.T) {
	type testcase struct {
		name     string
		notes    []note.Note
		slug     string
		expectOk bool
	}

	notes := []note.Note{
		{Title: "Root 1"},
		{
			Title: "Root 2",
			Children: []note.Note{
				{Title: "Child 1"},
				{
					Title: "Child 2",
					Children: []note.Note{
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
			notes:    []note.Note{},
			slug:     "any",
			expectOk: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, ok := note.FindBySlug(tc.notes, tc.slug)

			if ok != tc.expectOk {
				t.Errorf("expected ok=%v, got ok=%v", tc.expectOk, ok)
			}

			if tc.expectOk {
				if result.Slug() != tc.slug {
					t.Errorf("expected slug %q, got %q", tc.slug, result.Slug())
				}
			} else {
				exportSetInternals := cmp.AllowUnexported(util.NewSet[note.Tag]())
				if !cmp.Equal(result, note.Note{}, exportSetInternals) {
					t.Errorf("expected empty Note{}, got %+v", result)
				}
			}
		})
	}
}

func TestWalk(t *testing.T) {
	notes := []note.Note{
		{Title: "Root 1"},
		{
			Title: "Root 2",
			Children: []note.Note{
				{Title: "Child 1"},
				{
					Title: "Child 2",
					Children: []note.Note{
						{Title: "Grandchild 1"},
					},
				},
			},
		},
		{Title: "Root 3"},
	}

	t.Run("visits every note in depth-first, pre-order", func(t *testing.T) {
		var visited []string

		result := note.Walk(notes, func(n note.Note) bool {
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

		result := note.Walk(notes, func(n note.Note) bool {
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

		result := note.Walk(notes, func(n note.Note) bool {
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

		result := note.Walk([]note.Note{}, func(n note.Note) bool {
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

func TestResolveSlug(t *testing.T) {
	type testcase struct {
		name        string
		notes       []note.Note
		target      string
		expected    string
		expectFound bool
	}

	notes := []note.Note{
		{
			Title: "01.02.2026",
			Date:  time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
			Children: []note.Note{
				{
					Title: "Rehearsal",
					Date:  time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Title: "31.01.2026",
			Date:  time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC),
			Children: []note.Note{
				{
					Title: "Rehearsal",
					Date:  time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Title:    "Root 1",
			Children: []note.Note{{Title: "Child 1"}},
		},
	}

	testcases := []testcase{
		{
			name:        "note at root level",
			notes:       notes,
			target:      "Root 1",
			expected:    "root-1",
			expectFound: true,
		},
		{
			name:        "note in children",
			notes:       notes,
			target:      "Child 1",
			expected:    "child-1",
			expectFound: true,
		},
		{
			name:        "ignores case",
			notes:       notes,
			target:      "root 1",
			expected:    "root-1",
			expectFound: true,
		},
		{
			name:        "ignores punctuation",
			notes:       notes,
			target:      "Root 1!",
			expected:    "root-1",
			expectFound: true,
		},
		{
			name:        "daily note resolves to its date",
			notes:       notes,
			target:      "01.02.2026",
			expected:    "2026-02-01",
			expectFound: true,
		},
		{
			name:        "child note of the first matching day wins",
			notes:       notes,
			target:      "Rehearsal",
			expected:    "2026-02-01-rehearsal",
			expectFound: true,
		},
		{
			name:        "no match falls back to the slugified target",
			notes:       notes,
			target:      "Some Other Note",
			expected:    "some-other-note",
			expectFound: false,
		},
		{
			name:        "empty slice",
			notes:       []note.Note{},
			target:      "Root 1",
			expected:    "root-1",
			expectFound: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, found := note.ResolveSlug(tc.notes, tc.target)

			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}

			if found != tc.expectFound {
				t.Errorf("expected found=%v, got found=%v", tc.expectFound, found)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type testcase struct {
		name          string
		notes         []note.Note
		expected      []string
		expectDidSort bool
	}

	testcases := []testcase{
		{
			name: "sorts daily notes by date (descending)",
			notes: []note.Note{
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
			name: "sorts undated notes alphabetically (ascending)",
			notes: []note.Note{
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
			name: "groups all daily notes before undated notes",
			notes: []note.Note{
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
			notes: []note.Note{
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
			notes: []note.Note{
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
			result, didSort := note.Sort(tc.notes)
			if didSort != tc.expectDidSort {
				t.Errorf("expected sorted to be %v, got %v", tc.expectDidSort, didSort)
			}

			for i, n := range result {
				if n.Title != tc.expected[i] {
					t.Errorf("expected note at %d to be %v, got %v", i, tc.expected[i], n.Title)
					t.FailNow()
				}
			}
		})
	}
}

func TestSortNormalizesNewline(t *testing.T) {
	input := "# B\n\nLine 1\n\n# A\n\nLine 2\n"
	expected := "# A\n\nLine 2\n\n# B\n\nLine 1\n"

	notes, err := note.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	notes, _ = note.Sort(notes)

	result := note.String(notes)
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestString(t *testing.T) {
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
			notes, err := note.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			result := note.String(notes)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestTags(t *testing.T) {
	type testcase struct {
		name     string
		notes    []note.Note
		expected []note.Tag
	}

	testcases := []testcase{
		{
			name:     "returns nothing for no notes",
			notes:    []note.Note{},
			expected: []note.Tag{},
		},
		{
			name: "returns nothing for notes without tags",
			notes: []note.Note{
				{Title: "A"},
				{Title: "B", Tags: util.NewSet[note.Tag]()},
			},
			expected: []note.Tag{},
		},
		{
			name: "returns tags sorted alphabetically",
			notes: []note.Note{
				{Title: "A", Tags: util.NewSetFrom([]note.Tag{"#foo", "#baz"})},
				{Title: "B", Tags: util.NewSetFrom([]note.Tag{"#bar"})},
			},
			expected: []note.Tag{"#bar", "#baz", "#foo"},
		},
		{
			name: "sorts without regard to case",
			notes: []note.Note{
				{Title: "A", Tags: util.NewSetFrom([]note.Tag{"#Beta", "#alpha", "#Gamma"})},
			},
			expected: []note.Tag{"#alpha", "#Beta", "#Gamma"},
		},
		{
			name: "returns each tag only once",
			notes: []note.Note{
				{Title: "A", Tags: util.NewSetFrom([]note.Tag{"#foo", "#bar"})},
				{Title: "B", Tags: util.NewSetFrom([]note.Tag{"#foo"})},
			},
			expected: []note.Tag{"#bar", "#foo"},
		},
		{
			name: "keeps tags that differ only in case apart",
			notes: []note.Note{
				{Title: "A", Tags: util.NewSetFrom([]note.Tag{"#Foo", "#foo"})},
			},
			expected: []note.Tag{"#Foo", "#foo"},
		},
		{
			name: "includes tags of children",
			notes: []note.Note{
				{
					Title: "A",
					Tags:  util.NewSetFrom([]note.Tag{"#foo"}),
					Children: []note.Note{
						{Title: "B", Tags: util.NewSetFrom([]note.Tag{"#bar"})},
						{
							Title: "C",
							Children: []note.Note{
								{Title: "D", Tags: util.NewSetFrom([]note.Tag{"#baz"})},
							},
						},
					},
				},
			},
			expected: []note.Tag{"#bar", "#baz", "#foo"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := note.Tags(tc.notes)

			if diff := cmp.Diff(tc.expected, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("unexpected tags (-want +got):\n%v", diff)
			}
		})
	}
}
