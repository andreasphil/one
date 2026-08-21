package web_test

import (
	"net/http"
	"testing"
)

// tagNotes contains a note tagged with #groceries.
const tagNotes = `# Groceries #groceries

Buy milk and eggs.
`

func TestGetTagRedirectsToSearch(t *testing.T) {
	handler, _ := newTestServer(t, tagNotes)

	rec := get(t, handler, "/tags/groceries/")

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "/search/?query=%23groceries" {
		t.Errorf("expected redirect to the search for the tag, got %q", loc)
	}
}

func TestGetTagRedirectsToTrailingSlash(t *testing.T) {
	handler, _ := newTestServer(t, tagNotes)

	rec := get(t, handler, "/tags/groceries")

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "/tags/groceries/" {
		t.Errorf("expected redirect to /tags/groceries/, got %q", loc)
	}
}

func TestGetNoteLinksTagsToTagRoute(t *testing.T) {
	handler, _ := newTestServer(t, tagNotes)

	rec := get(t, handler, "/notes/groceries/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	assertContainsAll(t, rec.Body.String(), `<a class="tag" href="/tags/groceries/">`)
}
