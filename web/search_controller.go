package web

import (
	"html/template"
	"net/http"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

type searchResult struct {
	Title string
	Slug  string
	Date  time.Time
	Html  template.HTML
}

func getSearch(provider NotesProvider, renderer MarkdownRenderer) http.HandlerFunc {
	type getSearchData struct {
		Notes   []note.Note
		Results []searchResult
		Query   string
	}

	render := newRenderFunc[getSearchData]("get_search.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := provider.Notes()

		query := r.URL.Query().Get("query")

		results := []searchResult{}
		for _, match := range note.Containing(notes, query) {
			html, err := renderer.Render(match.Content())
			if err != nil {
				util.HttpErrorf(w, http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
				return
			}

			var date time.Time
			if !match.IsDailyNote() {
				date = match.Date
			}

			results = append(results, searchResult{
				Title: match.Title,
				Slug:  match.Slug(),
				Date:  date,
				Html:  html,
			})
		}

		err := render(w, data[getSearchData]{
			Title:      "Search",
			CurrentUrl: r.URL.Path,
			Data:       getSearchData{Notes: notes, Results: results, Query: query},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}
