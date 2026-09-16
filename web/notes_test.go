package web_test

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/web"
	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update the golden files in testdata/")

type fakeNotesProvider []note.Note

func (f fakeNotesProvider) Notes() []note.Note { return f }

func parseNotes(t *testing.T, markdown string) []note.Note {
	t.Helper()

	notes, err := note.Parse(strings.NewReader(markdown))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	return notes
}

func newTestRouter(t *testing.T, markdown string) (http.Handler, []note.Note) {
	t.Helper()

	notes := parseNotes(t, markdown)
	router := web.NewRouter(web.RouterArgs{Notes: fakeNotesProvider(notes)})

	return router, notes
}

func get(t *testing.T, handler http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func assertContainsAll(t *testing.T, body string, want ...string) {
	t.Helper()

	for _, w := range want {
		if !strings.Contains(body, w) {
			t.Errorf("body does not contain %q, got:\n%s", w, body)
		}
	}
}

func assertGolden(t *testing.T, name string, got string) {
	t.Helper()

	path := filepath.Join("testdata", name)

	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
		}

		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", path, err)
		}

		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v (run `go test ./web -update` to create it)", path, err)
	}

	if diff := cmp.Diff(string(want), got); diff != "" {
		t.Errorf("%s mismatch (-want +got), run `go test ./web -update` to accept:\n%s", path, diff)
	}
}

func mainOf(t *testing.T, body string) string {
	t.Helper()

	start := strings.Index(body, "<main>")
	end := strings.Index(body, "</main>")
	if start < 0 || end < start {
		t.Fatalf("body has no <main> element, got:\n%s", body)
	}

	return body[start:end]
}

func TestRootRedirectsToNotesList(t *testing.T) {
	router, _ := newTestRouter(t, "")

	rec := get(t, router, "/")

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusTemporaryRedirect)
	}

	if loc := rec.Header().Get("Location"); loc != "/notes/" {
		t.Errorf("Location = %q, want %q", loc, "/notes/")
	}
}

func TestGetNotesListsAllNotes(t *testing.T) {
	router, _ := newTestRouter(t, "# First note\n\nHello.\n\n# Second note #tag\n\nWorld.\n")

	rec := get(t, router, "/notes/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertGolden(t, "notes_list.html", rec.Body.String())
}

func TestGetNotesListsChildNotes(t *testing.T) {
	router, _ := newTestRouter(t,
		"# 01.02.2026\n\nStandup at 9.\n\n## Groceries run\n\nWent to the store.\n\n# Reading list\n\nBooks.\n")

	rec := get(t, router, "/notes/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertGolden(t, "notes_list_with_children.html", rec.Body.String())
}

func TestNavigationListsChildNotesFlat(t *testing.T) {
	router, _ := newTestRouter(t,
		"# 01.02.2026\n\nStandup at 9.\n\n## Groceries run\n\nWent to the store.\n\n# Reading list\n\nBooks.\n")

	rec := get(t, router, "/notes/")

	body := rec.Body.String()
	nav := body[strings.Index(body, `<nav class="navigation">`):strings.Index(body, "</aside>")]

	assertContainsAll(t, nav,
		"3 Notes",
		`href="/notes/2026-02-01-groceries-run/"`,
	)

	// The sidebar is flat, so it only has the list of links at the top and the
	// list of notes below it, with nothing nested inside either.
	if n := strings.Count(nav, "<ul>"); n != 2 {
		t.Errorf("lists in the navigation = %d, want 2 (it should be flat):\n%s", n, nav)
	}
}

func TestGetNotesShowsExcerpts(t *testing.T) {
	router, _ := newTestRouter(t, "# First note\n\nHello **world**, this is\nthe excerpt.\n")

	rec := get(t, router, "/notes/")

	assertContainsAll(t, mainOf(t, rec.Body.String()), "Hello world, this is the excerpt.")
}

func TestGetNotesOmitsExcerptForEmptyNote(t *testing.T) {
	router, _ := newTestRouter(t, "# Empty note\n")

	rec := get(t, router, "/notes/")

	content := mainOf(t, rec.Body.String())

	assertContainsAll(t, content, "Empty note")

	if strings.Contains(content, "clamp") {
		t.Errorf("content has an excerpt, want none for a note without content:\n%s", content)
	}
}

func TestGetNotesShowsIcons(t *testing.T) {
	router, _ := newTestRouter(t,
		"# \U0001F389 Party Planning\n\nLet us celebrate.\n\n# Plain note\n\nNo icon here.\n")

	rec := get(t, router, "/notes/")

	content := mainOf(t, rec.Body.String())

	assertContainsAll(t, content, `data-content="`+"\U0001F389"+`"`)

	if n := strings.Count(content, `class="glow"`); n != 1 {
		t.Errorf("icon = %d, want 1, content:\n%s", n, content)
	}
}

func TestGetNotesWithNoNotesShowsEmptyState(t *testing.T) {
	router, _ := newTestRouter(t, "")

	rec := get(t, router, "/notes/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertContainsAll(t, mainOf(t, rec.Body.String()),
		`has-fallback="empty"`,
		"0 notes",
		"Welcome!",
	)
}

func TestGetNoteRendersUndatedNote(t *testing.T) {
	router, notes := newTestRouter(t, "# My Guide #golang\n\nSome helpful content.\n")
	slug := notes[0].Slug()

	rec := get(t, router, "/notes/"+slug+"/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertGolden(t, "note_undated.html", rec.Body.String())
}

func TestGetNoteRendersDailyNoteWithChild(t *testing.T) {
	router, notes := newTestRouter(t, "# 01.02.2026\n\nDaily content.\n\n## Child A\n\nChild content.\n")
	root := notes[0]

	rec := get(t, router, "/notes/"+root.Slug()+"/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()

	// The formatted date is the one piece of the page that is not literally in
	// the fixture, so assert it by hand before pinning the rest.
	wantDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Format("02.01.2006 (Mon)")
	assertContainsAll(t, body, wantDate)

	assertGolden(t, "note_daily_with_child.html", body)
}

func TestGetNoteResolvesWikiLinks(t *testing.T) {
	router, _ := newTestRouter(t,
		"# 01.02.2026\n\nSee [[Child A]], [[01.02.2026]] and [[nope]].\n\n## Child A\n\nChild content.\n")

	rec := get(t, router, "/notes/2026-02-01/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertContainsAll(t, rec.Body.String(),
		`<a class="wikilink" href="/notes/2026-02-01-child-a/">Child A</a>`,
		`<a class="wikilink" href="/notes/2026-02-01/">01.02.2026</a>`,
		`<a class="wikilink unresolved" href="/notes/nope/">nope</a>`,
	)
}

func TestGetNoteChildLinksBackToParentDate(t *testing.T) {
	router, notes := newTestRouter(t, "# 01.02.2026\n\nDaily content.\n\n## Child A\n\nChild content.\n")
	child := notes[0].Children[0]

	rec := get(t, router, "/notes/"+child.Slug()+"/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()

	assertContainsAll(t, body,
		"<p>Child content.</p>",
		`href="/notes/2026-02-01/"`,
	)
}

func TestGetNoteWithoutContentShowsFallback(t *testing.T) {
	router, notes := newTestRouter(t, "# Empty Note\n")
	slug := notes[0].Slug()

	rec := get(t, router, "/notes/"+slug+"/")

	assertContainsAll(t, rec.Body.String(), "This note has no text.")
}

func TestGetNoteWithoutTagsShowsFallback(t *testing.T) {
	router, notes := newTestRouter(t, "# Untagged Note\n\nSome content.\n")
	slug := notes[0].Slug()

	rec := get(t, router, "/notes/"+slug+"/")

	assertContainsAll(t, rec.Body.String(), "This note has no tags.")
}

func TestGetNoteWithIconShowsGlow(t *testing.T) {
	router, notes := newTestRouter(t, "# \U0001F389 Party Planning\n\nLet's celebrate.\n")
	slug := notes[0].Slug()

	rec := get(t, router, "/notes/"+slug+"/")

	assertContainsAll(t, rec.Body.String(),
		"<title>Party Planning | One</title>",
		`data-content="`+"\U0001F389"+`"`,
	)
}

func TestGetNoteNotFoundReturns404(t *testing.T) {
	router, _ := newTestRouter(t, "# Only Note\n\nContent.\n")

	rec := get(t, router, "/notes/does-not-exist/")

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	assertContainsAll(t, rec.Body.String(), "does-not-exist")
}

func TestGetNoteMarksActiveNoteInNavigation(t *testing.T) {
	router, notes := newTestRouter(t, "# First note\n\nHello.\n\n# Second note\n\nWorld.\n")

	rec := get(t, router, "/notes/"+notes[0].Slug()+"/")

	if n := strings.Count(rec.Body.String(), `aria-current="page"`); n != 1 {
		t.Errorf("active nav entry = %d, want 1, body:\n%s", n, rec.Body.String())
	}
}

func TestStaticAssetsAreServed(t *testing.T) {
	router, _ := newTestRouter(t, "")

	rec := get(t, router, "/static/styles/styles.css")

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
