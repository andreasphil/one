package web_test

import (
	"net/http"
	"strings"
	"testing"
)

const tagNotes = `# Groceries #groceries #errands

Buy milk and eggs.

# Reading list #books

Books to read.

# 01.02.2026

Standup at 9.

## Groceries run #groceries

Went to the store for milk.
`

func TestGetTagsListsAllTags(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertGolden(t, "tags_list.html", rec.Body.String())
}

func TestGetTagsOmitsResultsWithoutTag(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/")

	body := rec.Body.String()

	for _, unwanted := range []string{"note tagged", "notes tagged", "No search results."} {
		if strings.Contains(body, unwanted) {
			t.Errorf("body contains %q, want no result section without a tag:\n%s", unwanted, body)
		}
	}
}

func TestGetTagsWithNoTagsShowsFallback(t *testing.T) {
	router, _ := newTestRouter(t, "# Untagged note\n\nContent.\n")

	rec := get(t, router, "/tags/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertContainsAll(t, mainOf(t, rec.Body.String()),
		`has-fallback="empty"`,
		"0 tags",
		"No tags.",
	)
}

func TestGetTagsMarksTagsActiveInNavigation(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	for _, path := range []string{"/tags/", "/tags/groceries/"} {
		rec := get(t, router, path)

		if n := strings.Count(rec.Body.String(), `aria-current="page"`); n != 1 {
			t.Errorf("active nav entry for %v = %d, want 1, body:\n%s", path, n, rec.Body.String())
		}
	}
}

func TestGetTagListsTaggedNotes(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/groceries/")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertGolden(t, "tag_results.html", rec.Body.String())
}

func TestGetTagRendersNoteContentAsHTML(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/books/")

	assertContainsAll(t, rec.Body.String(), "<p>Books to read.</p>")
}

func TestGetTagCountIsSingularForOneNote(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/books/")

	body := rec.Body.String()

	assertContainsAll(t, body, "1 note tagged")

	if strings.Contains(body, "1 notes tagged") {
		t.Errorf("body contains %q, want the singular note count:\n%s", "1 notes tagged", body)
	}
}

func TestGetTagShowsTagInTitle(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/groceries/")

	assertContainsAll(t, rec.Body.String(), "<title>#groceries | One</title>")
}

func TestGetTagMarksCurrentTagInList(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/groceries/")

	body := rec.Body.String()

	if n := strings.Count(body, `aria-current="location"`); n != 1 {
		t.Errorf("current tag = %d, want 1, body:\n%s", n, body)
	}

	link := mainOf(t, body)
	start := strings.Index(link, `href="/tags/groceries/"`)
	if start < 0 {
		t.Fatalf("body has no link to the current tag, got:\n%s", body)
	}
	link = link[start : start+strings.Index(link[start:], "</a>")]

	assertContainsAll(t, link, `aria-current="location"`)
}

func TestGetTagUnknownShowsFallbackWith404(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/nonexistent/")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	assertGolden(t, "tag_not_found.html", rec.Body.String())
}

func TestGetTagIsCaseSensitive(t *testing.T) {
	router, _ := newTestRouter(t, tagNotes)

	rec := get(t, router, "/tags/Groceries/")

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
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
