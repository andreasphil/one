package web

import (
	"html/template"
	"net/http"
	"time"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
	"github.com/andreasphil/one/web/adapter"
)

type searchResult struct {
	Title string
	Slug  string
	Date  time.Time
	Html  template.HTML
}

func getSearch(notes adapter.NotesProvider, renderer adapter.MarkdownRenderer) http.HandlerFunc {
	type getSearchData struct {
		Notes   []lib.Note
		Results []searchResult
		Query   string
	}

	render := newRenderFunc[getSearchData]("get_search.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := notes.Notes()

		query := r.URL.Query().Get("query")

		results := []searchResult{}
		for _, result := range lib.SearchContainingString(notes, query) {
			html, err := renderer.Render(result.Content())
			if err != nil {
				util.HttpErrorf(w, http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
				return
			}

			var date time.Time
			if !result.IsDailyNote() {
				date = result.Date
			}

			results = append(results, searchResult{
				Title: result.Title,
				Slug:  result.Slug(),
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
