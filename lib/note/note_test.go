package note_test

import (
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
)

func TestNew(t *testing.T) {
	n := note.New("Test Title")

	if n.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %q", n.Title)
	}

	if n.Tags.Len() != 0 {
		t.Errorf("expected empty tags set, got %d tags", n.Tags.Len())
	}

	if !n.Date.IsZero() {
		t.Errorf("expected zero date, got %v", n.Date)
	}

	if len(n.Children) != 0 {
		t.Errorf("expected nil or empty children, got %d children", len(n.Children))
	}
}

func TestSlug(t *testing.T) {
	type testcase struct {
		name     string
		note     note.Note
		expected string
	}

	testcases := []testcase{
		{
			name:     "regular note without date",
			note:     note.Note{Title: "My Note"},
			expected: "my-note",
		},
		{
			name: "daily note",
			note: note.Note{
				Title: "01.01.2026",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: "2026-01-01",
		},
		{
			name: "note with date and non-date title",
			note: note.Note{
				Title: "Meeting Notes",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: "2026-01-01-meeting-notes",
		},
		{
			name:     "special characters normalization",
			note:     note.Note{Title: "Hello! World? ? Test!"},
			expected: "hello-world-test",
		},
		{
			name:     "multiple consecutive special characters",
			note:     note.Note{Title: "test---note"},
			expected: "test-note",
		},
		{
			name:     "german umlauts",
			note:     note.Note{Title: "Äpfel Öl Über"},
			expected: "äpfel-öl-über",
		},
		{
			name:     "leading and trailing special chars",
			note:     note.Note{Title: "---test---"},
			expected: "test",
		},
		{
			name:     "empty title",
			note:     note.Note{Title: ""},
			expected: "",
		},
		{
			name:     "only special characters",
			note:     note.Note{Title: "!!!"},
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
			n := note.Note{Raw: tc.raw}
			result := n.Content()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestExcerpt(t *testing.T) {
	type testcase struct {
		name      string
		raw       string
		expected  string
		shortened bool
	}

	testcases := []testcase{
		{
			name:      "strips level 1 heading",
			raw:       "# Heading\n\nThis is content",
			expected:  "This is content",
			shortened: false,
		},
		{
			name:      "strips level 2 heading",
			raw:       "## Heading\n\nThis is content",
			expected:  "This is content",
			shortened: false,
		},
		{
			name:      "less than word limit",
			raw:       "# Title\n\nShort content here",
			expected:  "Short content here",
			shortened: false,
		},
		{
			name:      "more than word limit",
			raw:       "# Title\n\n" + strings.Repeat("word ", 50),
			expected:  strings.TrimSpace(strings.Repeat("word ", 40)),
			shortened: true,
		},
		{
			name:      "only heading",
			raw:       "# Only Title\n",
			expected:  "",
			shortened: false,
		},
		{
			name:      "trims whitespace",
			raw:       "# Title\n\n\nContent after newlines",
			expected:  "Content after newlines",
			shortened: false,
		},
		{
			name:      "removes line breaks",
			raw:       "# Title\nThis is\ncontent",
			expected:  "This is content",
			shortened: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			n := note.Note{Raw: tc.raw}
			result, shortened := n.Excerpt(40)
			if result != tc.expected {
				t.Errorf("expected excerpt %q, got %q", tc.expected, result)
			}
			if shortened != tc.shortened {
				t.Errorf("expected shortened to be %v, got %v", tc.shortened, shortened)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	type testcase struct {
		name     string
		note     note.Note
		expected bool
	}

	testcases := []testcase{
		{
			name:     "returns that note is empty",
			note:     note.Note{Raw: "# Title\n\n"},
			expected: true,
		},
		{
			name:     "returns that note with title is empty",
			note:     note.Note{Raw: ""},
			expected: true,
		},
		{
			name:     "returns that note has content",
			note:     note.Note{Raw: "# Title\n\nContent"},
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

func TestIsDailyNote(t *testing.T) {
	type testcase struct {
		name     string
		note     note.Note
		expected bool
	}

	testcases := []testcase{
		{
			name: "valid daily note",
			note: note.Note{
				Title: "01.01.2026",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: true,
		},
		{
			name: "date set but title doesn't match",
			note: note.Note{
				Title: "Meeting Notes",
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: false,
		},
		{
			name: "title matches date format but no date set",
			note: note.Note{
				Title: "01.01.2026",
				Date:  time.Time{},
			},
			expected: false,
		},
		{
			name: "neither date nor matching title",
			note: note.Note{
				Title: "Regular Note",
				Date:  time.Time{},
			},
			expected: false,
		},
		{
			name: "child note with date from parent",
			note: note.Note{
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

func TestNoteString(t *testing.T) {
	type testcase struct {
		name     string
		input    string
		expected string
	}

	testcases := []testcase{
		{
			name:     "returns raw for note without children",
			input:    "# Note 1\n\nLine 1\n",
			expected: "# Note 1\n\nLine 1\n",
		},
		{
			name:     "includes child notes in serialized daily note",
			input:    "# 01.01.2026\n\nLine 1\n\n## Child Note 1\n\nLine 2\n\n## Child Note 2\n\nLine 3\n",
			expected: "# 01.01.2026\n\nLine 1\n\n## Child Note 1\n\nLine 2\n\n## Child Note 2\n\nLine 3\n",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			notes, err := note.Parse(strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			result := notes[0].String()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
