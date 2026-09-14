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
				{Title: "Child 1", Kind: note.KindChild},
				{Title: "Child 2", Kind: note.KindChild},
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
			got, ok := note.FindBySlug(tc.notes, tc.slug)

			if ok != tc.expectOk {
				t.Errorf("FindBySlug(notes, %q) ok = %v, want %v", tc.slug, ok, tc.expectOk)
			}

			if tc.expectOk {
				if got.Slug() != tc.slug {
					t.Errorf("FindBySlug(notes, %q) slug = %q, want %q", tc.slug, got.Slug(), tc.slug)
				}
			} else if diff := cmp.Diff(note.Note{}, got); diff != "" {
				t.Errorf("FindBySlug(notes, %q) mismatch (-want +got):\n%s", tc.slug, diff)
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
				{Title: "Child 1", Kind: note.KindChild},
				{Title: "Child 2", Kind: note.KindChild},
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

		want := []string{
			"Root 1", "Root 2", "Child 1", "Child 2", "Root 3",
		}

		if diff := cmp.Diff(want, visited); diff != "" {
			t.Errorf("Walk() visited mismatch (-want +got):\n%s", diff)
		}

		if !result {
			t.Errorf("Walk() = false, want true")
		}
	})

	t.Run("stops early when fn returns false", func(t *testing.T) {
		var visited []string

		result := note.Walk(notes, func(n note.Note) bool {
			visited = append(visited, n.Title)
			return n.Title != "Child 1"
		})

		want := []string{"Root 1", "Root 2", "Child 1"}

		if diff := cmp.Diff(want, visited); diff != "" {
			t.Errorf("Walk() visited mismatch (-want +got):\n%s", diff)
		}

		if result {
			t.Errorf("Walk() = true, want false")
		}
	})

	t.Run("stopping in children also stops parent traversal", func(t *testing.T) {
		var visited []string

		result := note.Walk(notes, func(n note.Note) bool {
			visited = append(visited, n.Title)
			return n.Title != "Child 2"
		})

		want := []string{"Root 1", "Root 2", "Child 1", "Child 2"}

		if diff := cmp.Diff(want, visited); diff != "" {
			t.Errorf("Walk() visited mismatch (-want +got):\n%s", diff)
		}

		if result {
			t.Errorf("Walk() = true, want false")
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		var visited []string

		result := note.Walk([]note.Note{}, func(n note.Note) bool {
			visited = append(visited, n.Title)
			return true
		})

		if len(visited) != 0 {
			t.Errorf("Walk() visited = %v, want none", visited)
		}

		if !result {
			t.Errorf("Walk() = false, want true")
		}
	})
}

func TestFlat(t *testing.T) {
	notes := []note.Note{
		{Title: "Root 1"},
		{
			Title: "Root 2",
			Children: []note.Note{
				{Title: "Child 1", Kind: note.KindChild},
				{Title: "Child 2", Kind: note.KindChild},
			},
		},
		{Title: "Root 3"},
	}

	t.Run("returns children directly after the note they belong to", func(t *testing.T) {
		var got []string
		for _, n := range note.Flat(notes) {
			got = append(got, n.Title)
		}

		want := []string{"Root 1", "Root 2", "Child 1", "Child 2", "Root 3"}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Flat() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("returns a non-nil empty slice for no notes", func(t *testing.T) {
		for name, notes := range map[string][]note.Note{"nil": nil, "empty": {}} {
			t.Run(name, func(t *testing.T) {
				got := note.Flat(notes)

				if got == nil {
					t.Fatalf("Flat(%s) = nil, want a non-nil slice", name)
				}

				if len(got) != 0 {
					t.Errorf("Flat(%s) = %v, want an empty slice", name, got)
				}
			})
		}
	})
}

func TestCount(t *testing.T) {
	type testcase struct {
		name     string
		notes    []note.Note
		expected int
	}

	testcases := []testcase{
		{
			name:     "no notes",
			notes:    []note.Note{},
			expected: 0,
		},
		{
			name:     "notes without children",
			notes:    []note.Note{{Title: "Root 1"}, {Title: "Root 2"}},
			expected: 2,
		},
		{
			name: "counts children as notes of their own",
			notes: []note.Note{
				{
					Title: "Root 1",
					Children: []note.Note{
						{Title: "Child 1", Kind: note.KindChild},
						{Title: "Child 2", Kind: note.KindChild},
					},
				},
				{Title: "Root 2"},
			},
			expected: 4,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := note.Count(tc.notes); got != tc.expected {
				t.Errorf("Count() = %d, want %d", got, tc.expected)
			}

			if got := len(note.Flat(tc.notes)); got != tc.expected {
				t.Errorf("len(Flat()) = %d, want %d", got, tc.expected)
			}
		})
	}
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
			Kind:  note.KindDaily,
			Date:  time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
			Children: []note.Note{
				{
					Title: "Rehearsal",
					Kind:  note.KindChild,
					Date:  time.Date(2026, time.February, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Title: "31.01.2026",
			Kind:  note.KindDaily,
			Date:  time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC),
			Children: []note.Note{
				{
					Title: "Rehearsal",
					Kind:  note.KindChild,
					Date:  time.Date(2026, time.January, 31, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			Title:    "Root 1",
			Children: []note.Note{{Title: "Child 1", Kind: note.KindChild}},
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
			name:        "full slug picks a specific child note",
			notes:       notes,
			target:      "2026-01-31-rehearsal",
			expected:    "2026-01-31-rehearsal",
			expectFound: true,
		},
		{
			name:        "full slug resolves a daily note",
			notes:       notes,
			target:      "2026-02-01",
			expected:    "2026-02-01",
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
			got, found := note.ResolveSlug(tc.notes, tc.target)

			if got != tc.expected {
				t.Errorf("ResolveSlug(notes, %q) = %q, want %q", tc.target, got, tc.expected)
			}

			if found != tc.expectFound {
				t.Errorf("ResolveSlug(notes, %q) found = %v, want %v", tc.target, found, tc.expectFound)
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
				t.Errorf("Sort() didSort = %v, want %v", didSort, tc.expectDidSort)
			}

			got := make([]string, 0, len(result))
			for _, n := range result {
				got = append(got, n.Title)
			}

			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Fatalf("Sort() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSortNormalizesNewline(t *testing.T) {
	input := "# B\n\nLine 1\n\n# A\n\nLine 2\n"
	expected := "# A\n\nLine 2\n\n# B\n\nLine 1\n"

	notes, err := note.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	notes, _ = note.Sort(notes)

	if got := note.String(notes); got != expected {
		t.Errorf("String(Sort(notes)) = %q, want %q", got, expected)
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
				t.Fatalf("Parse() error = %v", err)
			}

			if got := note.String(notes); got != tc.expected {
				t.Errorf("String(%q) = %q, want %q", tc.input, got, tc.expected)
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
				{Title: "A", Tags: util.NewSetFrom([]note.Tag{note.NewTag("foo"), note.NewTag("baz")})},
				{Title: "B", Tags: util.NewSetFrom([]note.Tag{note.NewTag("bar")})},
			},
			expected: []note.Tag{note.NewTag("bar"), note.NewTag("baz"), note.NewTag("foo")},
		},
		{
			name: "sorts without regard to case",
			notes: []note.Note{
				{Title: "A", Tags: util.NewSetFrom([]note.Tag{note.NewTag("Beta"), note.NewTag("alpha"), note.NewTag("Gamma")})},
			},
			expected: []note.Tag{note.NewTag("alpha"), note.NewTag("Beta"), note.NewTag("Gamma")},
		},
		{
			name: "returns each tag only once",
			notes: []note.Note{
				{Title: "A", Tags: util.NewSetFrom([]note.Tag{note.NewTag("foo"), note.NewTag("bar")})},
				{Title: "B", Tags: util.NewSetFrom([]note.Tag{note.NewTag("foo")})},
			},
			expected: []note.Tag{note.NewTag("bar"), note.NewTag("foo")},
		},
		{
			name: "keeps tags that differ only in case apart",
			notes: []note.Note{
				{Title: "A", Tags: util.NewSetFrom([]note.Tag{note.NewTag("Foo"), note.NewTag("foo")})},
			},
			expected: []note.Tag{note.NewTag("Foo"), note.NewTag("foo")},
		},
		{
			name: "includes tags of children",
			notes: []note.Note{
				{
					Title: "A",
					Tags:  util.NewSetFrom([]note.Tag{note.NewTag("foo")}),
					Children: []note.Note{
						{Title: "B", Kind: note.KindChild, Tags: util.NewSetFrom([]note.Tag{note.NewTag("bar")})},
						{Title: "C", Kind: note.KindChild, Tags: util.NewSetFrom([]note.Tag{note.NewTag("baz")})},
					},
				},
			},
			expected: []note.Tag{note.NewTag("bar"), note.NewTag("baz"), note.NewTag("foo")},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := note.Tags(tc.notes)

			if diff := cmp.Diff(tc.expected, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Tags() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
