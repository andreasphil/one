package web_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
	"github.com/andreasphil/one/web"
	"github.com/google/go-cmp/cmp"
)

func newTestMarkdownRenderer(notes []note.Note) web.MarkdownRenderer {
	return web.NewMarkdownRenderer(func(target string) (string, bool) {
		return note.ResolveSlug(notes, target)
	})
}

func TestNewNoteMeta(t *testing.T) {
	type testcase struct {
		name     string
		note     note.Note
		expected web.NoteMeta
	}

	testcases := []testcase{
		{
			name:     "maps title and slug",
			note:     note.Note{Title: "Hello World"},
			expected: web.NoteMeta{Title: "Hello World", Slug: "hello-world"},
		},
		{
			name:     "maps daily note",
			note:     note.Note{Title: "01.02.2026", Kind: note.KindDaily, Date: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
			expected: web.NoteMeta{Title: "01.02.2026", Slug: "2026-02-01", Date: "2026-02-01"},
		},
		{
			name:     "maps child note",
			note:     note.Note{Title: "Groceries run", Kind: note.KindChild, Date: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
			expected: web.NoteMeta{Title: "Groceries run", Slug: "2026-02-01-groceries-run", Date: "2026-02-01", IsChildNote: true},
		},
		{
			name:     "maps empty note",
			note:     note.Note{},
			expected: web.NoteMeta{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := web.NewNoteMeta(tc.note)

			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("NewNoteMeta() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewNoteMetaOmitsZeroDate(t *testing.T) {
	got, err := json.Marshal(web.NewNoteMeta(note.Note{Title: "Hello World"}))
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	if strings.Contains(string(got), "Date") {
		t.Errorf("json.Marshal(NewNoteMeta()) = %s, want no date", got)
	}
}

func TestMapToNoteMeta(t *testing.T) {
	date := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	// mapToNoteMeta maps the flat list of notes each page is rendered with, see
	// note.Flat.
	notes := note.Flat([]note.Note{
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
		got := web.MapToNoteMeta(notes)

		want := []web.NoteMeta{
			{Title: "Root 1", Slug: "root-1"},
			{Title: "01.02.2026", Slug: "2026-02-01", Date: "2026-02-01"},
			{Title: "Child 1", Slug: "2026-02-01-child-1", Date: "2026-02-01", IsChildNote: true},
			{Title: "Child 2", Slug: "2026-02-01-child-2", Date: "2026-02-01", IsChildNote: true},
			{Title: "Root 3", Slug: "root-3"},
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("MapToNoteMeta() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("returns an empty slice for no notes", func(t *testing.T) {
		for name, notes := range map[string][]note.Note{"nil": nil, "empty": {}} {
			t.Run(name, func(t *testing.T) {
				got := web.MapToNoteMeta(notes)

				if got == nil {
					t.Fatalf("MapToNoteMeta(%s) = nil, want a non-nil slice", name)
				}

				if len(got) != 0 {
					t.Errorf("MapToNoteMeta(%s) = %+v, want an empty slice", name, got)
				}
			})
		}
	})

	t.Run("serializes an empty slice to an empty JSON array", func(t *testing.T) {
		got, err := json.Marshal(web.MapToNoteMeta(nil))
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		if string(got) != "[]" {
			t.Errorf("json.Marshal(MapToNoteMeta(nil)) = %q, want %q", got, "[]")
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
		got := web.MapToTags(notes)
		want := []string{"Idea", "recipe", "work"}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("MapToTags() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("returns an empty slice for no notes", func(t *testing.T) {
		for name, notes := range map[string][]note.Note{"nil": nil, "empty": {}} {
			t.Run(name, func(t *testing.T) {
				got := web.MapToTags(notes)

				if got == nil {
					t.Fatalf("MapToTags(%s) = nil, want a non-nil slice", name)
				}

				if len(got) != 0 {
					t.Errorf("MapToTags(%s) = %+v, want an empty slice", name, got)
				}
			})
		}
	})

	t.Run("serializes an empty slice to an empty JSON array", func(t *testing.T) {
		got, err := json.Marshal(web.MapToTags(nil))
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}

		if string(got) != "[]" {
			t.Errorf("json.Marshal(MapToTags(nil)) = %q, want %q", got, "[]")
		}
	})
}

func toSearchResults(notes []note.Note) []note.Result {
	results := make([]note.Result, 0, len(notes))
	for _, n := range notes {
		results = append(results, note.Result{Note: n})
	}

	return results
}

func TestNewSearchResult(t *testing.T) {
	renderer := newTestMarkdownRenderer(nil)
	date := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	type testcase struct {
		name     string
		note     note.Note
		expected web.SearchResult
	}

	testcases := []testcase{
		{
			name: "maps title, slug and content",
			note: note.Note{Title: "Groceries", Raw: "# Groceries\n\nBuy **milk**.\n"},
			expected: web.SearchResult{
				Title: "Groceries",
				Slug:  "groceries",
				HTML:  "<p>Buy <strong>milk</strong>.</p>\n",
			},
		},
		{
			name: "maps the date of a child note",
			note: note.Note{Title: "Groceries run", Kind: note.KindChild, Date: date, Raw: "## Groceries run\n"},
			expected: web.SearchResult{
				Title: "Groceries run",
				Slug:  "2026-02-01-groceries-run",
				Date:  date,
			},
		},
		{
			name: "omits the date of a daily note",
			note: note.Note{Title: "01.02.2026", Kind: note.KindDaily, Date: date, Raw: "# 01.02.2026\n"},
			expected: web.SearchResult{
				Title: "01.02.2026",
				Slug:  "2026-02-01",
			},
		},
		{
			name:     "maps empty note",
			note:     note.Note{},
			expected: web.SearchResult{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := web.NewSearchResult(note.Result{Note: tc.note}, renderer)
			if err != nil {
				t.Fatalf("NewSearchResult() error = %v", err)
			}

			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("NewSearchResult() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestNewSearchResultResolvesWikiLinks(t *testing.T) {
	notes := parseNotes(t, "# Groceries\n\nSee [[Reading list]] and [[Nonexistent]].\n\n# Reading list\n\nBooks.\n")
	renderer := newTestMarkdownRenderer(notes)

	result, err := web.NewSearchResult(note.Result{Note: notes[0]}, renderer)
	if err != nil {
		t.Fatalf("NewSearchResult() error = %v", err)
	}

	html := string(result.HTML)

	if !strings.Contains(html, `<a class="wikilink" href="/notes/reading-list/">`) {
		t.Errorf("HTML does not link the existing note, got:\n%s", html)
	}

	if !strings.Contains(html, `<a class="wikilink unresolved" href="/notes/nonexistent/">`) {
		t.Errorf("HTML does not mark the link to the missing note unresolved, got:\n%s", html)
	}
}

func TestMapToSearchResults(t *testing.T) {
	renderer := newTestMarkdownRenderer(nil)

	t.Run("maps each note to a result, without flattening children", func(t *testing.T) {
		notes := parseNotes(t, "# Groceries\n\nBuy milk.\n\n## Groceries run\n\nWent to the store.\n\n# Reading list\n\nBooks.\n")

		results, err := web.MapToSearchResults(toSearchResults(notes), renderer)
		if err != nil {
			t.Fatalf("MapToSearchResults() error = %v", err)
		}

		got := make([]string, 0, len(results))
		for _, result := range results {
			got = append(got, result.Title)
		}

		want := []string{"Groceries", "Reading list"}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("MapToSearchResults() titles mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("returns an empty slice for no notes", func(t *testing.T) {
		for name, results := range map[string][]note.Result{"nil": nil, "empty": {}} {
			t.Run(name, func(t *testing.T) {
				got, err := web.MapToSearchResults(results, renderer)
				if err != nil {
					t.Fatalf("MapToSearchResults() error = %v", err)
				}

				if got == nil {
					t.Fatalf("MapToSearchResults(%s) = nil, want a non-nil slice", name)
				}

				if len(got) != 0 {
					t.Errorf("MapToSearchResults(%s) = %+v, want an empty slice", name, got)
				}
			})
		}
	})
}
