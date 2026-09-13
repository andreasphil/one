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
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTemporaryRedirect)
	}

	if loc := rec.Header().Get("Location"); loc != "/search/?query=%23groceries" {
		t.Errorf("Location = %q, want %q", loc, "/search/?query=%23groceries")
	}
}

func TestGetTagRedirectsToTrailingSlash(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/groceries")

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTemporaryRedirect)
	}

	if loc := rec.Header().Get("Location"); loc != "/tags/groceries/" {
		t.Errorf("Location = %q, want %q", loc, "/tags/groceries/")
	}
}

func TestGetNoteLinksTagsToTagRoute(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/notes/groceries/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertContainsAll(t, rec.Body.String(), `<a class="tag" href="/tags/groceries/">`)
}
