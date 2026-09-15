package web

import (
	"net/http"

	"github.com/andreasphil/one/lib/note"
)

type getTagData struct {
	Results []searchResult
	Tag     string
}

func getTags(provider NotesProvider) handler {
	render := newRenderFunc[getTagData](provider, "get_tags.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		return render(w, r, data[getTagData]{Title: "Tags"})
	}
}

func getTag(provider NotesProvider, renderer markdownRenderer) handler {
	render := newRenderFunc[getTagData](provider, "get_tags.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		tag := note.NewTag(r.PathValue("tag"))

		notes := note.Search(provider.Notes(), note.FilterChain{note.FilterHasTag(tag)})

		results, err := mapToSearchResults(notes, renderer)
		if err != nil {
			return httpStatusErrorf(http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
		}

		if len(results) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}

		title := tag.String()

		return render(w, r, data[getTagData]{
			Title: title,
			Data:  getTagData{Results: results, Tag: tag.Name()},
		})
	}
}
