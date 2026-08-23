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

// fakeNotesProvider satisfies adapter.NotesProvider without needing an
// interface import - Go interface satisfaction is structural.
type fakeNotesProvider []note.Note

func (f fakeNotesProvider) Notes() []note.Note { return f }

// newTestServer parses markdown into notes and wires up a real web.Server
// backed by those notes, exercising the actual router, handlers and
// templates - not any internal function directly.
func newTestServer(t *testing.T, markdown string) (http.Handler, []note.Note) {
	t.Helper()

	notes, err := note.Parse(strings.NewReader(markdown))
	if err != nil {
		t.Fatalf("failed to parse test notes: %v", err)
	}

	server := web.NewServer(web.ServerArgs{Notes: fakeNotesProvider(notes), Port: "0"})
	return server.Handler, notes
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

func TestRootRedirectsToNotesList(t *testing.T) {
	handler, _ := newTestServer(t, "")

	rec := get(t, handler, "/")

	if rec.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "/notes/" {
		t.Errorf("expected redirect to /notes/, got %q", loc)
	}
}

func TestGetNotesListsAllNotes(t *testing.T) {
	handler, _ := newTestServer(t, "# First note\n\nHello.\n\n# Second note #tag\n\nWorld.\n")

	rec := get(t, handler, "/notes/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	assertContainsAll(t, rec.Body.String(),
		"<title>Notes | One</title>",
		"2 Notes",
		"First note",
		"Second note",
		`href="/notes/first-note/"`,
		`href="/notes/second-note/"`,
	)
}

func TestGetNotesWithNoNotesShowsZeroCount(t *testing.T) {
	handler, _ := newTestServer(t, "")

	rec := get(t, handler, "/notes/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	assertContainsAll(t, rec.Body.String(), "0 Notes")
}

func TestGetNoteRendersUndatedNote(t *testing.T) {
	handler, notes := newTestServer(t, "# My Guide #golang\n\nSome helpful content.\n")
	slug := notes[0].Slug()

	rec := get(t, handler, "/notes/"+slug+"/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	assertContainsAll(t, body,
		"<title>My Guide | One</title>",
		"Some helpful content.",
		"Knowledge Base", // the UI label for a note without a date
		"golang",         // tag rendered without its leading "#"
	)

	if n := strings.Count(body, `class="tag"`); n != 1 {
		t.Errorf("expected exactly 1 tag, got %d, body:\n%s", n, body)
	}
}

func TestGetNoteRendersDailyNoteWithChild(t *testing.T) {
	handler, notes := newTestServer(t, "# 01.02.2026\n\nDaily content.\n\n## Child A\n\nChild content.\n")
	root := notes[0]

	rec := get(t, handler, "/notes/"+root.Slug()+"/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	wantDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC).Format("Monday, 2. January 2006")

	assertContainsAll(t, body,
		wantDate,
		"Daily content.",
		"Also on this day",
		"Child A",
	)

	// The daily note's own page shows its date as plain text, not a link.
	if strings.Contains(body, `href="/notes/2026-02-01"`) {
		t.Errorf("did not expect date to be a link on the daily note's own page, got:\n%s", body)
	}
}

func TestGetNoteChildLinksBackToParentDate(t *testing.T) {
	handler, notes := newTestServer(t, "# 01.02.2026\n\nDaily content.\n\n## Child A\n\nChild content.\n")
	child := notes[0].Children[0]

	rec := get(t, handler, "/notes/"+child.Slug()+"/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	assertContainsAll(t, body,
		"Child content.",
		`href="/notes/2026-02-01/"`, // links back to the parent day
	)
}

func TestGetNoteWithoutContentShowsFallback(t *testing.T) {
	handler, notes := newTestServer(t, "# Empty Note\n")
	slug := notes[0].Slug()

	rec := get(t, handler, "/notes/"+slug+"/")

	assertContainsAll(t, rec.Body.String(), "This note has no text.")
}

func TestGetNoteWithoutTagsShowsFallback(t *testing.T) {
	handler, notes := newTestServer(t, "# Untagged Note\n\nSome content.\n")
	slug := notes[0].Slug()

	rec := get(t, handler, "/notes/"+slug+"/")

	assertContainsAll(t, rec.Body.String(), "This note has no tags.")
}

func TestGetNoteWithIconShowsGlow(t *testing.T) {
	handler, notes := newTestServer(t, "# \U0001F389 Party Planning\n\nLet's celebrate.\n")
	slug := notes[0].Slug()

	rec := get(t, handler, "/notes/"+slug+"/")

	assertContainsAll(t, rec.Body.String(),
		"<title>Party Planning | One</title>",
		`data-content="`+"\U0001F389"+`"`,
	)
}

func TestGetNoteNotFoundReturns404(t *testing.T) {
	handler, _ := newTestServer(t, "# Only Note\n\nContent.\n")

	rec := get(t, handler, "/notes/does-not-exist/")

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}

	assertContainsAll(t, rec.Body.String(), "does-not-exist")
}

func TestGetNoteMarksActiveNoteInNavigation(t *testing.T) {
	handler, notes := newTestServer(t, "# First note\n\nHello.\n\n# Second note\n\nWorld.\n")

	rec := get(t, handler, "/notes/"+notes[0].Slug()+"/")

	if n := strings.Count(rec.Body.String(), `aria-current="page"`); n != 1 {
		t.Errorf("expected exactly 1 active nav entry, got %d, body:\n%s", n, rec.Body.String())
	}
}

func TestStaticAssetsAreServed(t *testing.T) {
	handler, _ := newTestServer(t, "")

	rec := get(t, handler, "/static/styles/styles.css")

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
