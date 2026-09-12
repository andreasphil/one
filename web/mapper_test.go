package web

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
	"github.com/google/go-cmp/cmp"
)

func TestNewNoteMeta(t *testing.T) {
	type testcase struct {
		name     string
		note     note.Note
		expected noteMeta
	}

	testcases := []testcase{
		{
			name:     "maps title and slug",
			note:     note.Note{Title: "Hello World"},
			expected: noteMeta{Title: "Hello World", Slug: "hello-world"},
		},
		{
			name:     "maps daily note",
			note:     note.Note{Title: "01.02.2026", Kind: note.KindDaily, Date: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
			expected: noteMeta{Title: "01.02.2026", Slug: "2026-02-01"},
		},
		{
			name:     "maps empty note",
			note:     note.Note{},
			expected: noteMeta{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := newNoteMeta(tc.note)

			if !cmp.Equal(result, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestMapToNoteMeta(t *testing.T) {
	date := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	// mapToNoteMeta maps the flat list of notes each page is rendered with, see
	// note.Flatten.
	notes := note.Flatten([]note.Note{
		{Title: "Root 1"},
		{
			Title: "01.02.2026",
			Kind:  note.KindDaily,
			Date:  date,
			Children: []note.Note{
				{Title: "Child 1", Kind: note.KindChild, Date: date},
				{Title: "Child 2", Kind: note.KindChild, Date: date},
			},
		},
		{Title: "Root 3"},
	})

	t.Run("maps every note, keeping the order of the input", func(t *testing.T) {
		result := mapToNoteMeta(notes)

		expected := []noteMeta{
			{Title: "Root 1", Slug: "root-1"},
			{Title: "01.02.2026", Slug: "2026-02-01"},
			{Title: "Child 1", Slug: "2026-02-01-child-1"},
			{Title: "Child 2", Slug: "2026-02-01-child-2"},
			{Title: "Root 3", Slug: "root-3"},
		}

		if !cmp.Equal(result, expected) {
			t.Errorf("expected %+v, got %+v", expected, result)
		}
	})

	t.Run("returns an empty slice for no notes", func(t *testing.T) {
		for _, notes := range [][]note.Note{nil, {}} {
			result := mapToNoteMeta(notes)

			if result == nil {
				t.Fatalf("expected a non-nil slice, got nil")
			}

			if len(result) != 0 {
				t.Errorf("expected an empty slice, got %+v", result)
			}
		}
	})

	t.Run("serializes an empty slice to an empty JSON array", func(t *testing.T) {
		result, err := json.Marshal(mapToNoteMeta(nil))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if string(result) != "[]" {
			t.Errorf("expected %q, got %q", "[]", string(result))
		}
	})
}

func TestMapToTags(t *testing.T) {
	notes := []note.Note{
		{Title: "Root 1", Tags: util.NewSetFrom([]note.Tag{"#work", "#Idea"})},
		{
			Title: "Root 2",
			Tags:  util.NewSetFrom([]note.Tag{"#work"}),
			Children: []note.Note{
				{Title: "Child 1", Kind: note.KindChild, Tags: util.NewSetFrom([]note.Tag{"#recipe"})},
			},
		},
	}

	t.Run("collects unique tags of notes and children without the leading #", func(t *testing.T) {
		result := mapToTags(notes)

		expected := []string{"Idea", "recipe", "work"}

		if !cmp.Equal(result, expected) {
			t.Errorf("expected %+v, got %+v", expected, result)
		}
	})

	t.Run("returns an empty slice for no notes", func(t *testing.T) {
		for _, notes := range [][]note.Note{nil, {}} {
			result := mapToTags(notes)

			if result == nil {
				t.Fatalf("expected a non-nil slice, got nil")
			}

			if len(result) != 0 {
				t.Errorf("expected an empty slice, got %+v", result)
			}
		}
	})

	t.Run("serializes an empty slice to an empty JSON array", func(t *testing.T) {
		result, err := json.Marshal(mapToTags(nil))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if string(result) != "[]" {
			t.Errorf("expected %q, got %q", "[]", string(result))
		}
	})
}

func newTestMarkdownRenderer(notes []note.Note) markdownRenderer {
	return newMarkdownRenderer(func(target string) (string, bool) {
		return note.ResolveSlug(notes, target)
	})
}

func parseTestNotes(t *testing.T, markdown string) []note.Note {
	t.Helper()

	notes, err := note.Parse(strings.NewReader(markdown))
	if err != nil {
		t.Fatalf("failed to parse test notes: %v", err)
	}

	return notes
}

func TestNewSearchResult(t *testing.T) {
	renderer := newTestMarkdownRenderer(nil)
	date := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	type testcase struct {
		name     string
		note     note.Note
		expected searchResult
	}

	testcases := []testcase{
		{
			name: "maps title, slug and content",
			note: note.Note{Title: "Groceries", Raw: "# Groceries\n\nBuy **milk**.\n"},
			expected: searchResult{
				Title: "Groceries",
				Slug:  "groceries",
				HTML:  "<p>Buy <strong>milk</strong>.</p>\n",
			},
		},
		{
			name: "maps the date of a child note",
			note: note.Note{Title: "Groceries run", Kind: note.KindChild, Date: date, Raw: "## Groceries run\n"},
			expected: searchResult{
				Title: "Groceries run",
				Slug:  "2026-02-01-groceries-run",
				Date:  date,
			},
		},
		{
			name: "omits the date of a daily note",
			note: note.Note{Title: "01.02.2026", Kind: note.KindDaily, Date: date, Raw: "# 01.02.2026\n"},
			expected: searchResult{
				Title: "01.02.2026",
				Slug:  "2026-02-01",
			},
		},
		{
			name:     "maps empty note",
			note:     note.Note{},
			expected: searchResult{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := newSearchResult(tc.note, renderer)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !cmp.Equal(result, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, result)
			}
		})
	}
}

func TestNewSearchResultResolvesWikiLinks(t *testing.T) {
	notes := parseTestNotes(t, "# Groceries\n\nSee [[Reading list]] and [[Nonexistent]].\n\n# Reading list\n\nBooks.\n")
	renderer := newTestMarkdownRenderer(notes)

	result, err := newSearchResult(notes[0], renderer)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	html := string(result.HTML)

	if !strings.Contains(html, `<a class="wikilink" href="/notes/reading-list/">`) {
		t.Errorf("expected a link to the existing note, got:\n%s", html)
	}

	if !strings.Contains(html, `<a class="wikilink unresolved" href="/notes/nonexistent/">`) {
		t.Errorf("expected the link to the missing note to be marked unresolved, got:\n%s", html)
	}
}

func TestMapToSearchResults(t *testing.T) {
	renderer := newTestMarkdownRenderer(nil)

	t.Run("maps each note to a result, without flattening children", func(t *testing.T) {
		notes := parseTestNotes(t, "# Groceries\n\nBuy milk.\n\n## Groceries run\n\nWent to the store.\n\n# Reading list\n\nBooks.\n")

		results, err := mapToSearchResults(notes, renderer)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		expected := []string{"Groceries", "Reading list"}

		titles := make([]string, 0, len(results))
		for _, result := range results {
			titles = append(titles, result.Title)
		}

		if !cmp.Equal(titles, expected) {
			t.Errorf("expected %+v, got %+v", expected, titles)
		}
	})

	t.Run("returns an empty slice for no notes", func(t *testing.T) {
		for _, notes := range [][]note.Note{nil, {}} {
			result, err := mapToSearchResults(notes, renderer)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if result == nil {
				t.Fatalf("expected a non-nil slice, got nil")
			}

			if len(result) != 0 {
				t.Errorf("expected an empty slice, got %+v", result)
			}
		}
	})
}
