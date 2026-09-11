package web_test

import (
	"net/http"
	"strings"
	"testing"
)

const searchNotes = `# Groceries

Buy **milk** and eggs.

# Reading list

Books to read.

# 01.02.2026

Standup at 9.

## Groceries run

Went to the store for milk.
`

func TestGetSearchWithoutQueryShowsOnlySearchBox(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	assertContainsAll(t, body,
		"<title>Search | One</title>",
		`id="searchbox"`,
		`value=""`,
		"3 Notes",
	)

	for _, unwanted := range []string{"result for", "results for", "No search results."} {
		if strings.Contains(body, unwanted) {
			t.Errorf("did not expect body to contain %q, got:\n%s", unwanted, body)
		}
	}
}

func TestGetSearchMarksSearchActiveInNavigation(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/")

	if n := strings.Count(rec.Body.String(), `aria-current="page"`); n != 1 {
		t.Errorf("expected exactly 1 active nav entry, got %d, body:\n%s", n, rec.Body.String())
	}
}

func TestGetSearchListsMatchingNotes(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=milk")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	assertContainsAll(t, body,
		"2 results for",
		`<a href="/notes/groceries/">Groceries</a>`,
		`<a href="/notes/2026-02-01-groceries-run/">Groceries run</a>`,
	)

	if strings.Contains(body, "Reading list</a>") {
		t.Errorf("did not expect a result for the non-matching note, got:\n%s", body)
	}
}

func TestGetSearchRendersResultContentAsHTML(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=eggs")

	assertContainsAll(t, rec.Body.String(),
		"<p>Buy <strong>milk</strong> and eggs.</p>",
	)
}

func TestGetSearchCountIsSingularForOneResult(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=Books")

	body := rec.Body.String()

	assertContainsAll(t, body, "1 result for")

	if strings.Contains(body, "1 results for") {
		t.Errorf("expected a singular result count, got:\n%s", body)
	}
}

func TestGetSearchKeepsQueryInSearchBox(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=milk")

	assertContainsAll(t, rec.Body.String(), `value="milk"`)
}

func TestGetSearchShowsQueryInTitle(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=milk")

	assertContainsAll(t, rec.Body.String(), `<title>Search for &#34;milk&#34; | One</title>`)
}

func TestGetSearchEscapesQuery(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=%3Cscript%3Ealert(1)%3C%2Fscript%3E")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()

	assertContainsAll(t, body, "&lt;script&gt;alert(1)&lt;/script&gt;")

	if strings.Contains(body, "<script>alert(1)</script>") {
		t.Errorf("expected the query to be escaped, got:\n%s", body)
	}
}

func TestGetSearchShowsDateOfChildNotes(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=store")

	assertContainsAll(t, rec.Body.String(),
		`<a href="/notes/2026-02-01-groceries-run/">Groceries run</a>`,
		"(01.02.2026)",
	)
}

func TestGetSearchOmitsDateOfDailyNote(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=Standup")

	body := rec.Body.String()

	assertContainsAll(t, body, `<a href="/notes/2026-02-01/">01.02.2026</a>`)

	if strings.Contains(body, "(on 01.02.2026)") {
		t.Errorf("did not expect the date of the daily note to be repeated, got:\n%s", body)
	}
}

func TestGetSearchOmitsDateOfUndatedNote(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=Books")

	if body := rec.Body.String(); strings.Contains(body, "(on ") {
		t.Errorf("did not expect a date for a note without one, got:\n%s", body)
	}
}

func TestGetSearchWithoutResultsShowsFallback(t *testing.T) {
	router, _ := newTestRouter(t, searchNotes)

	rec := get(t, router, "/search/?query=nonexistent")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	assertContainsAll(t, rec.Body.String(),
		"0 results for",
		`has-fallback="empty"`,
		"No search results.",
	)
}

func TestGetSearchWithNoNotesShowsFallback(t *testing.T) {
	router, _ := newTestRouter(t, "")

	rec := get(t, router, "/search/?query=anything")

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	assertContainsAll(t, rec.Body.String(), "0 results for", "No search results.")
}
