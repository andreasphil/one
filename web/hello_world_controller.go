package web

import (
	"net/http"

	"github.com/andreasphil/one/util"
)

func getHelloWorld() http.HandlerFunc {
	render := newRenderFunc[nilPage]("get_hello_world.html")

	return func(w http.ResponseWriter, r *http.Request) {
		err := render(w, data[nilPage]{
			Title:      "Under construction",
			CurrentUrl: r.URL.Path,
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}
