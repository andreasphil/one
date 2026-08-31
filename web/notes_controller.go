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

	render := newRenderFunc[getNotesData]("get_notes.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		notes := provider.Notes()

		return render(w, data[getNotesData]{
			Title:      "Notes",
			CurrentURL: r.URL.Path,
			Data:       getNotesData{Notes: notes},
		})
	}
}

func getNote(provider NotesProvider, renderer MarkdownRenderer) handler {
	type getNoteData struct {
		Notes []note.Note
		Note  note.Note
		HTML  template.HTML
	}

	render := newRenderFunc[getNoteData]("get_note.html")

	return func(w http.ResponseWriter, r *http.Request) error {
		notes := provider.Notes()
		slug := r.PathValue("slug")

		n, found := note.FindBySlug(notes, slug)
		if !found {
			return httpStatusErrorf(http.StatusNotFound, "note %v not found", slug)
		}

		html, err := renderer.Render(n.Content())
		if err != nil {
			return httpStatusErrorf(http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
		}

		return render(w, data[getNoteData]{
			Title:      n.Title,
			CurrentURL: r.URL.Path,
			Data:       getNoteData{Notes: notes, Note: n, HTML: html},
		})
	}
}
