package web

import (
	"fmt"
	"net/http"

	"github.com/andreasphil/one/lib/note"
)

func getSearch(provider NotesProvider, renderer markdownRenderer) handler {
	type getSearchData struct {
		Results []searchResult
		Query   string
	}

	render := newRenderFunc[getSearchData](provider, "get_search.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		query := r.URL.Query().Get("query")

		notes := note.Search(provider.Notes(), note.FilterChain{note.FilterExactPhrase(query, false)})

		results, err := mapToSearchResults(notes, renderer)
		if err != nil {
			return httpStatusErrorf(http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
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
