package web

import (
	"html/template"
	"time"

	"github.com/andreasphil/one/lib/note"
)

// Note meta ----------------------------------------------

// Should mirror NoteMeta in static/scripts/lib/types.ts.
type noteMeta struct {
	Title       string
	Slug        string
	Date        string `json:",omitempty"`
	IsChildNote bool
}

func newNoteMeta(n note.Note) noteMeta {
	var date string
	if !n.Date.IsZero() {
		date = n.Date.Format("2006-01-02")
	}

	return noteMeta{
		Title:       n.Title,
		Slug:        n.Slug(),
		Date:        date,
		IsChildNote: n.IsChildNote(),
	}
}

func mapToNoteMeta(n []note.Note) []noteMeta {
	m := make([]noteMeta, 0, len(n))

	for _, i := range n {
		m = append(m, newNoteMeta(i))
	}

	return m
}

// Tags ---------------------------------------------------

func mapToTags(n []note.Note) []string {
	tags := note.Tags(n)
	m := make([]string, 0, len(tags))

	for _, tag := range tags {
		m = append(m, tag.Name())
	}

	return m
}

// Search results -----------------------------------------

type searchResult struct {
	Title string
	Slug  string
	Date  time.Time
	HTML  template.HTML
}

func newSearchResult(n note.Note, renderer markdownRenderer) (searchResult, error) {
	html, err := renderer.render(n.Content())
	if err != nil {
		return searchResult{}, err
	}

	var date time.Time
	if n.IsChildNote() {
		date = n.Date
	}

	return searchResult{Title: n.Title, Slug: n.Slug(), Date: date, HTML: html}, nil
}

func mapToSearchResults(n []note.Note, renderer markdownRenderer) ([]searchResult, error) {
	m := make([]searchResult, 0, len(n))

	for _, i := range n {
		result, err := newSearchResult(i, renderer)
		if err != nil {
			return nil, err
		}

		m = append(m, result)
	}

	return m, nil
}
