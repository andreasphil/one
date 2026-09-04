package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/andreasphil/one/lib/note"
)

type searchResult struct {
	Title string
	Slug  string
	Date  time.Time
	HTML  template.HTML
}

func getSearch(provider NotesProvider, renderer MarkdownRenderer) handler {
	type getSearchData struct {
		Results []searchResult
		Query   string
	}

	render := newRenderFunc[getSearchData](provider, "get_search.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		query := r.URL.Query().Get("query")

		results := []searchResult{}
		for _, match := range note.Containing(provider.Notes(), query) {
			html, err := renderer.Render(match.Content())
			if err != nil {
				return httpStatusErrorf(http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
			}

			var date time.Time
			if !match.IsDailyNote() {
				date = match.Date
			}

			results = append(results, searchResult{
				Title: match.Title,
				Slug:  match.Slug(),
				Date:  date,
				HTML:  html,
			})
		}

		title := "Search"
		if query != "" {
			title += fmt.Sprintf(` for "%v"`, query)
		}

		return render(w, r, data[getSearchData]{
			Title: title,
			Data:  getSearchData{Results: results, Query: query},
		})
	}
}
