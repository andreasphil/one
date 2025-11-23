package web

import (
	"net/http"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
)

func getNotes() http.HandlerFunc {
	type getNotesData struct {
		Notes []lib.Note
	}

	render := newRenderFunc[getNotesData]("get_notes.html")

	return func(w http.ResponseWriter, r *http.Request) {
		err := render(w, data[getNotesData]{
			Title:      "Notes",
			CurrentUrl: r.URL.Path,
			Data:       getNotesData{},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
		}
	}
}
