package web

import (
	"net/http"
	"net/url"

	"github.com/andreasphil/one/lib/note"
)

func getTags(provider NotesProvider) handler {
	render := newRenderFunc[struct{}](provider, "get_tags.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		return render(w, r, data[struct{}]{Title: "Tags"})
	}
}

func getTag() handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		tag := r.PathValue("tag")

		target := url.URL{
			Path:     "/search/",
			RawQuery: url.Values{"query": {note.NewTag(tag).String()}}.Encode(),
		}

		http.Redirect(w, r, target.String(), http.StatusTemporaryRedirect)

		return nil
	}
}
