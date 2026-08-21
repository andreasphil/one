package web

import (
	"net/http"
	"net/url"
)

func getTag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tag := r.PathValue("tag")

		target := url.URL{
			Path:     "/search/",
			RawQuery: url.Values{"query": {"#" + tag}}.Encode(),
		}

		http.Redirect(w, r, target.String(), http.StatusTemporaryRedirect)
	}
}
