package lib_test

import (
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
	"github.com/google/go-cmp/cmp"
)

func TestNewNote(t *testing.T) {
	note := lib.NewNote("Test Title")

	if note.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %q", note.Title)
	}

	if note.Tags.Len() != 0 {
		t.Errorf("expected empty tags set, got %d tags", note.Tags.Len())
	}

	if !note.Date.IsZero() {
		t.Errorf("expected zero date, got %v", note.Date)
	}

	if len(note.Children) != 0 {
		t.Errorf("expected nil or empty children, got %d children", len(note.Children))
	}
}

func TestSlug(t *testing.T) {
	type testcase struct {
		name     string
		note     lib.Note
		expected string
	}

	testcases := []testcase{
		{
			name:     "regular note without date",
			note:     lib.Note{Title: "My Note"},
			expected: "my-note",
		},
		{
			name: "daily note",
			note: lib.Note{
				Title: "01.01.2026",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: "2026-01-01",
		},
		{
			name: "note with date and non-date title",
			note: lib.Note{
				Title: "Meeting Notes",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: "2026-01-01-meeting-notes",
		},
		{
			name:     "special characters normalization",
			note:     lib.Note{Title: "Hello! World? ? Test!"},
			expected: "hello-world-test",
		},
		{
			name:     "multiple consecutive special characters",
			note:     lib.Note{Title: "test---note"},
			expected: "test-note",
		},
		{
			name:     "german umlauts",
			note:     lib.Note{Title: "Äpfel Öl Über"},
			expected: "äpfel-öl-über",
		},
		{
			name:     "leading and trailing special chars",
			note:     lib.Note{Title: "---test---"},
			expected: "test",
		},
		{
			name:     "empty title",
			note:     lib.Note{Title: ""},
			expected: "",
		},
		{
			name:     "only special characters",
			note:     lib.Note{Title: "!!!"},
			expected: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.note.Slug()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestContent(t *testing.T) {
	type testcase struct {
		name     string
		raw      string
		expected string
	}

	testcases := []testcase{
		{
			name:     "basic note with h1 title and content",
			raw:      "# My Title\n\nThis is the content",
			expected: "This is the content",
		},
		{
			name:     "child note with h2 title and content",
			raw:      "## Child Note\n\nThis is child content",
			expected: "This is child content",
		},
		{
			name:     "only title, no content",
			raw:      "# Just a title\n",
			expected: "",
		},
		{
			name:     "multiline content",
			raw:      "# Title\n\nLine 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "content with level 2 heading preserved",
			raw:      "# Title\n\n## Subtitle\n\nContent",
			expected: "## Subtitle\n\nContent",
		},
		{
			name:     "empty raw content",
			raw:      "",
			expected: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			note := lib.Note{Raw: tc.raw}
			result := note.Content()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestExcerpt(t *testing.T) {
	type testcase struct {
		name     string
		raw      string
		expected string
	}

	testcases := []testcase{
		{
			name:     "strips level 1 heading",
			raw:      "# Heading\n\nThis is content",
			expected: "This is content",
		},
		{
			name:     "strips level 2 heading",
			raw:      "## Heading\n\nThis is content",
			expected: "This is content",
		},
		{
			name:     "less than 40 words",
			raw:      "# Title\n\nShort content here",
			expected: "Short content here",
		},
		{
			name:     "more than 40 words",
			raw:      "# Title\n\n" + strings.Repeat("word ", 50),
			expected: strings.TrimSpace(strings.Repeat("word ", 40)),
		},
		{
			name:     "only heading",
			raw:      "# Only Title\n",
			expected: "",
		},
		{
			name:     "trims whitespace",
			raw:      "# Title\n\n\nContent after newlines",
			expected: "Content after newlines",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			note := lib.Note{Raw: tc.raw}
			result := note.Excerpt()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestIsDailyNote(t *testing.T) {
	type testcase struct {
		name     string
		note     lib.Note
		expected bool
	}

	testcases := []testcase{
		{
			name: "valid daily note",
			note: lib.Note{
				Title: "01.01.2026",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: true,
		},
		{
			name: "date set but title doesn't match",
			note: lib.Note{
				Title: "Meeting Notes",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
		{
			name: "title matches date format but no date set",
			note: lib.Note{
				Title: "01.01.2026",
				Date:  time.Time{},
			},
			expected: false,
		},
		{
			name: "neither date nor matching title",
			note: lib.Note{
				Title: "Regular Note",
				Date:  time.Time{},
			},
			expected: false,
		},
		{
			name: "child note with date from parent",
			note: lib.Note{
				Title: "Child Note",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.note.IsDailyNote()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestFindRecursive(t *testing.T) {
	type testcase struct {
		name      string
		notes     []lib.Note
		slug      string
		expectOk  bool
		expectIdx int
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
			result, ok := lib.FindRecursive(tc.notes, tc.slug)

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

func TestIsEmpty(t *testing.T) {
	type testcase struct {
		name     string
		note     lib.Note
		expected bool
	}

	testcases := []testcase{
		{
			name:     "returns that note is empty",
			note:     lib.Note{Raw: "# Title\n\n"},
			expected: true,
		},
		{
			name:     "returns that note with title is empty",
			note:     lib.Note{Raw: ""},
			expected: true,
		},
		{
			name:     "returns that note has content",
			note:     lib.Note{Raw: "# Title\n\nContent"},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.note.IsEmpty()
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
