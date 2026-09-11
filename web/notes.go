package web

import (
	"html/template"
	"net/http"

	"github.com/andreasphil/one/lib/note"
)

func getNotes(provider NotesProvider) handler {
	type getNotesData struct {
		Notes []note.Note
	}

	render := newRenderFunc[getNotesData](provider, "get_notes.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		return render(w, r, data[getNotesData]{
			Title: "Notes",
			Data:  getNotesData{Notes: provider.Notes()},
		})
	}
}

func getNote(provider NotesProvider, renderer markdownRenderer) handler {
	type getNoteData struct {
		Note note.Note
		HTML template.HTML
	}

	render := newRenderFunc[getNoteData](provider, "get_note.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		slug := r.PathValue("slug")

		n, found := note.FindBySlug(provider.Notes(), slug)
		if !found {
			return httpStatusErrorf(http.StatusNotFound, "note %v not found", slug)
		}

		html, err := renderer.render(n.Content())
		if err != nil {
			return httpStatusErrorf(http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
		}

		return render(w, r, data[getNoteData]{
			Title: n.Title,
			Data:  getNoteData{Note: n, HTML: html},
		})
	}
}
