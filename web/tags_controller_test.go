package web_test

import (
	"net/http"
	"testing"
)

const tagNotes = `# Groceries #groceries

Buy milk and eggs.
`

func TestGetTagRedirectsToSearch(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/groceries/")

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "/search/?query=%23groceries" {
		t.Errorf("expected redirect to the search for the tag, got %q", loc)
	}
}

func TestGetTagRedirectsToTrailingSlash(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/groceries")

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected status %d, got %d", http.StatusTemporaryRedirect, rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "/tags/groceries/" {
		t.Errorf("expected redirect to /tags/groceries/, got %q", loc)
	}
}

func TestGetNoteLinksTagsToTagRoute(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/notes/groceries/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	assertContainsAll(t, rec.Body.String(), `<a class="tag" href="/tags/groceries/">`)
}
