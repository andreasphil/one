package web_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/web"
)

type fakeNotesProvider []note.Note

func (f fakeNotesProvider) Notes() []note.Note { return f }

func newTestRouter(t *testing.T, markdown string) (http.Handler, []note.Note) {
	t.Helper()

	notes, err := note.Parse(strings.NewReader(markdown))
	if err != nil {
		t.Fatalf("failed to parse test notes: %v", err)
	}

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
			t.Errorf("expected body to contain %q, got:\n%s", w, body)
		}
	}
}

func mainOf(t *testing.T, body string) string {
	t.Helper()

	start := strings.Index(body, "<main>")
	end := strings.Index(body, "</main>")
	if start < 0 || end < start {
		t.Fatalf("expected body to contain a main element, got:\n%s", body)
	}

	return body[start:end]
}

func TestRootRedirectsToNotesList(t *testing.T) {
	router, _ := newTestRouter(t, "")

	rec := get(t, router, "/")

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "/notes/" {
		t.Errorf("expected redirect to /notes/, got %q", loc)
	}
}

func TestGetNotesListsAllNotes(t *testing.T) {
	router, _ := newTestRouter(t, "# First note\n\nHello.\n\n# Second note #tag\n\nWorld.\n")

	rec := get(t, router, "/notes/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	content := mainOf(t, body)

	assertContainsAll(t, body, "<title>Notes | One</title>")

	assertContainsAll(t, content,
		`has-fallback=""`,
		"2 notes",
		"First note",
		"Second note",
		`href="/notes/first-note/"`,
		`href="/notes/second-note/"`,
	)

	if n := strings.Count(content, `class="card"`); n != 2 {
		t.Errorf("expected exactly 2 cards, got %d, content:\n%s", n, content)
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
		t.Errorf("expected no excerpt for a note without content, got:\n%s", content)
	}
}

func TestGetNotesShowsIcons(t *testing.T) {
	router, _ := newTestRouter(t,
		"# \U0001F389 Party Planning\n\nLet us celebrate.\n\n# Plain note\n\nNo icon here.\n")

	rec := get(t, router, "/notes/")

	content := mainOf(t, rec.Body.String())

	assertContainsAll(t, content, `data-content="`+"\U0001F389"+`"`)

	if n := strings.Count(content, `class="glow"`); n != 1 {
		t.Errorf("expected exactly 1 icon, got %d, content:\n%s", n, content)
	}
}

func TestGetNotesWithNoNotesShowsEmptyState(t *testing.T) {
	router, _ := newTestRouter(t, "")

	rec := get(t, router, "/notes/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
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
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	assertContainsAll(t, body,
		"<title>My Guide | One</title>",
		"Some helpful content.",
		"Knowledge Base",
		"golang",
	)

	if n := strings.Count(body, `class="tag"`); n != 1 {
		t.Errorf("expected exactly 1 tag, got %d, body:\n%s", n, body)
	}
}

func TestGetNoteRendersDailyNoteWithChild(t *testing.T) {
	router, notes := newTestRouter(t, "# 01.02.2026\n\nDaily content.\n\n## Child A\n\nChild content.\n")
	root := notes[0]

	rec := get(t, router, "/notes/"+root.Slug()+"/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	wantDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Format("Mon, 2. Jan 2006")

	assertContainsAll(t, body,
		wantDate,
		"Daily content.",
		"Also on this day",
		"Child A",
	)

	if strings.Contains(body, `href="/notes/2026-02-01"`) {
		t.Errorf("did not expect date to be a link on the daily note's own page, got:\n%s", body)
	}
}

func TestGetNoteResolvesWikiLinks(t *testing.T) {
	router, _ := newTestRouter(t,
		"# 01.02.2026\n\nSee [[Child A]], [[01.02.2026]] and [[nope]].\n\n## Child A\n\nChild content.\n")

	rec := get(t, router, "/notes/2026-02-01/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
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
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	assertContainsAll(t, body,
		"Child content.",
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
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	assertContainsAll(t, rec.Body.String(), "does-not-exist")
}

func TestGetNoteMarksActiveNoteInNavigation(t *testing.T) {
	router, notes := newTestRouter(t, "# First note\n\nHello.\n\n# Second note\n\nWorld.\n")

	rec := get(t, router, "/notes/"+notes[0].Slug()+"/")

	if n := strings.Count(rec.Body.String(), `aria-current="page"`); n != 1 {
		t.Errorf("expected exactly 1 active nav entry, got %d, body:\n%s", n, rec.Body.String())
	}
}

func TestStaticAssetsAreServed(t *testing.T) {
	router, _ := newTestRouter(t, "")

	rec := get(t, router, "/static/styles/styles.css")

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
