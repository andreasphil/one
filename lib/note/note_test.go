package note_test

import (
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
)

func TestTag(t *testing.T) {
	for _, name := range []string{"kermit", "#kermit"} {
		tag := note.NewTag(name)

		if tag != note.Tag("#kermit") {
			t.Errorf("expected #kermit from %q, got %v", name, tag)
		}

		if tag.Name() != "kermit" {
			t.Errorf("expected kermit, got %q", tag.Name())
		}

		if tag.String() != "#kermit" {
			t.Errorf("expected #kermit, got %q", tag.String())
		}
	}
}

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
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.note.Slug()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}

			if again := note.Slug(result); again != result {
				t.Errorf("expected slug of %q to stay the same, got %q", result, again)
			}
		})
	}
}

func TestSlugFunc(t *testing.T) {
	type testcase struct {
		name     string
		input    string
		expected string
	}

	testcases := []testcase{
		{
			name:     "simple title",
			input:    "My Note",
			expected: "my-note",
		},
		{
			name:     "input that is already a slug",
			input:    "my-note",
			expected: "my-note",
		},
		{
			name:     "input that is already a note slug with date",
			input:    "2026-01-01-meeting-notes",
			expected: "2026-01-01-meeting-notes",
		},
		{
			name:     "special characters normalization",
			input:    "Hello! World? ? Test!",
			expected: "hello-world-test",
		},
		{
			name:     "multiple consecutive special characters",
			input:    "test---note",
			expected: "test-note",
		},
		{
			name:     "german umlauts",
			input:    "Äpfel Öl Über",
			expected: "äpfel-öl-über",
		},
		{
			name:     "leading and trailing special chars",
			input:    "---test---",
			expected: "test",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!!!",
			expected: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := note.Slug(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}

			if again := note.Slug(result); again != result {
				t.Errorf("expected slug of %q to stay the same, got %q", result, again)
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
		name     string
		raw      string
		expected string
	}

	testcases := []testcase{
		{
			name:     "empty note",
			raw:      "# Title\n",
			expected: "",
		},
		{
			name:     "plain content",
			raw:      "# Title\n\nJust some plain text",
			expected: "Just some plain text",
		},
		{
			name:     "emphasis and strong",
			raw:      "# Title\n\nSome *emphasized* and **strong** and _underscored_ and __double__ words",
			expected: "Some emphasized and strong and underscored and double words",
		},
		{
			name:     "inline code and strikethrough",
			raw:      "# Title\n\nRun `go test` and ~~forget~~ it",
			expected: "Run go test and forget it",
		},
		{
			name:     "links and images",
			raw:      "# Title\n\nSee [the docs](https://example.com) and ![a picture](pic.png) here",
			expected: "See the docs and a picture here",
		},
		{
			name:     "autolinks",
			raw:      "# Title\n\nFound at <https://example.com/page> today",
			expected: "Found at https://example.com/page today",
		},
		{
			name:     "headings, quotes and lists",
			raw:      "# Title\n\n## Section\n\n> A quote\n\n- First item\n* Second item\n+ Third item\n1. Fourth item",
			expected: "Section A quote First item Second item Third item Fourth item",
		},
		{
			name:     "nested list markers and checkboxes",
			raw:      "# Title\n\n- Outer\n  - Inner\n- [ ] Open\n- [x] Done",
			expected: "Outer Inner Open Done",
		},
		{
			name:     "callout with empty quote lines",
			raw:      "# Title\n\n> [!NOTE]\n>\n> Not responsible for singed eyebrows.\n>\n> Especially Beaker.",
			expected: "[!NOTE] Not responsible for singed eyebrows. Especially Beaker.",
		},
		{
			name:     "snake case preserved across lines",
			raw:      "# Title\n\nprevious_password: BorkBorkBork1\nnext_scheduled_change: soon",
			expected: "previous_password: BorkBorkBork1 next_scheduled_change: soon",
		},
		{
			name:     "underscore emphasis next to snake case",
			raw:      "# Title\n\nThe _old_ value of previous_password is gone",
			expected: "The old value of previous_password is gone",
		},
		{
			name:     "hashtag preserved at start of line",
			raw:      "# Title\n\nSome content\n#tech 💻",
			expected: "Some content #tech 💻",
		},
		{
			name:     "wiki links",
			raw:      "# Title\n\nSee [[Reading list]] and [[01.02.2026]] for more",
			expected: "See Reading list and 01.02.2026 for more",
		},
		{
			name:     "unclosed wiki link left alone",
			raw:      "# Title\n\nSee [[unclosed here",
			expected: "See [[unclosed here",
		},
		{
			name:     "table",
			raw:      "# Title\n\n| Fruit | Color |\n| --- | --- |\n| Apple | Green |\n| Plum | Purple |",
			expected: "Fruit Color Apple Green Plum Purple",
		},
		{
			name:     "table with alignment and no outer pipes",
			raw:      "# Title\n\nFruit | Color\n:--- | ---:\nApple | Green",
			expected: "Fruit Color Apple Green",
		},
		{
			name:     "code fences and horizontal rules",
			raw:      "# Title\n\nBefore\n\n```go\nfmt.Println()\n```\n\n---\n\nAfter",
			expected: "Before fmt.Println() After",
		},
		{
			name:     "line breaks and whitespace normalized",
			raw:      "# Title\n\nOne\n\n\nTwo\t\tthree   four\n",
			expected: "One Two three four",
		},
		{
			name:     "truncated to 30 words",
			raw:      "# Title\n\n" + strings.Repeat("word ", 50),
			expected: strings.TrimSpace(strings.Repeat("word ", 30)) + "...",
		},
		{
			name:     "exactly 30 words kept",
			raw:      "# Title\n\n" + strings.Repeat("word ", 30),
			expected: strings.TrimSpace(strings.Repeat("word ", 30)),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			n := note.Note{Raw: tc.raw}
			result := n.Excerpt()
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
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
