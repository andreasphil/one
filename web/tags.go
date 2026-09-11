package web

import (
	"net/http"
	"net/url"

	"github.com/andreasphil/one/lib/note"
)

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
