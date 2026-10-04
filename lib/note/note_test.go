package note_test

import (
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
)

func TestTag(t *testing.T) {
	for _, name := range []string{"kermit", "#kermit"} {
		t.Run(name, func(t *testing.T) {
			tag := note.NewTag(name)

			if tag != note.NewTag("kermit") {
				t.Errorf("NewTag(%q) = %q, want %q", name, tag, note.NewTag("kermit"))
			}

			if got := tag.Name(); got != "kermit" {
				t.Errorf("NewTag(%q).Name() = %q, want %q", name, got, "kermit")
			}

			if got := tag.String(); got != "#kermit" {
				t.Errorf("NewTag(%q).String() = %q, want %q", name, got, "#kermit")
			}
		})
	}
}

func TestNew(t *testing.T) {
	n := note.New("Test Title")

	if n.Title != "Test Title" {
		t.Errorf("New().Title = %q, want %q", n.Title, "Test Title")
	}

	if got := n.Tags.Len(); got != 0 {
		t.Errorf("New().Tags.Len() = %d, want 0", got)
	}

	if !n.Date.IsZero() {
		t.Errorf("New().Date = %v, want the zero time", n.Date)
	}

	if got := len(n.Children); got != 0 {
		t.Errorf("len(New().Children) = %d, want 0", got)
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
				Kind:  note.KindDaily,
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: "2026-01-01",
		},
		{
			name: "child note",
			note: note.Note{
				Title: "Meeting Notes",
				Kind:  note.KindChild,
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: "2026-01-01-meeting-notes",
		},
		{
			name: "child note whose title is a date",
			note: note.Note{
				Title: "31.12.2025",
				Kind:  note.KindChild,
				Date:  time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: "2026-01-01-31-12-2025",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.note.Slug()
			if got != tc.expected {
				t.Errorf("Note{Title: %q}.Slug() = %q, want %q", tc.note.Title, got, tc.expected)
			}

			if again := note.Slug(got); again != got {
				t.Errorf("Slug(%q) = %q, want it to stay the same", got, again)
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
			got := note.Slug(tc.input)
			if got != tc.expected {
				t.Errorf("Slug(%q) = %q, want %q", tc.input, got, tc.expected)
			}

			if again := note.Slug(got); again != got {
				t.Errorf("Slug(%q) = %q, want it to stay the same", got, again)
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

			if got := n.Content(); got != tc.expected {
				t.Errorf("Note{Raw: %q}.Content() = %q, want %q", tc.raw, got, tc.expected)
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
		{
			name:     "returns that note with children is not empty",
			note:     note.Note{Raw: "# Title\n\n", Children: []note.Note{{Raw: "## Child\n"}}},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.note.IsEmpty(); got != tc.expected {
				t.Errorf("Note{Raw: %q}.IsEmpty() = %v, want %v", tc.note.Raw, got, tc.expected)
			}
		})
	}
}

func TestKindPredicates(t *testing.T) {
	type testcase struct {
		name       string
		note       note.Note
		daily      bool
		child      bool
		standalone bool
	}

	date := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	testcases := []testcase{
		{
			name:  "daily note",
			note:  note.Note{Title: "01.01.2026", Kind: note.KindDaily, Date: date},
			daily: true,
		},
		{
			name:  "child note",
			note:  note.Note{Title: "Child Note", Kind: note.KindChild, Date: date},
			child: true,
		},
		{
			name:       "standalone note",
			note:       note.Note{Title: "Regular Note"},
			standalone: true,
		},
		{
			name:       "standalone note whose title looks like a date",
			note:       note.Note{Title: "01.01.2026"},
			standalone: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.note.IsDailyNote(); got != tc.daily {
				t.Errorf("IsDailyNote() = %v, want %v", got, tc.daily)
			}

			if got := tc.note.IsChildNote(); got != tc.child {
				t.Errorf("IsChildNote() = %v, want %v", got, tc.child)
			}

			if got := tc.note.IsStandalone(); got != tc.standalone {
				t.Errorf("IsStandalone() = %v, want %v", got, tc.standalone)
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
				t.Fatalf("Parse() error = %v", err)
			}

			if len(notes) == 0 {
				t.Fatalf("Parse(%q) = 0 notes, want at least 1", tc.input)
			}

			if got := notes[0].String(); got != tc.expected {
				t.Errorf("Note.String() = %q, want %q", got, tc.expected)
			}
		})
	}
}
