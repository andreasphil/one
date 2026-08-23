package web

import (
	"html/template"
	"net/http"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

func getNotes(provider NotesProvider) http.HandlerFunc {
	type getNotesData struct {
		Notes []note.Note
	}

	render := newRenderFunc[getNotesData]("get_notes.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := provider.Notes()

		err := render(w, data[getNotesData]{
			Title:      "Notes",
			CurrentUrl: r.URL.Path,
			Data:       getNotesData{Notes: notes},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}

func getNote(provider NotesProvider, renderer MarkdownRenderer) http.HandlerFunc {
	type getNoteData struct {
		Notes []note.Note
		Note  note.Note
		Html  template.HTML
	}

	render := newRenderFunc[getNoteData]("get_note.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := provider.Notes()
		slug := r.PathValue("slug")

		n, found := note.FindBySlug(notes, slug)
		if !found {
			util.HttpErrorf(w, http.StatusNotFound, "note %v not found", slug)
			return
		}

		html, err := renderer.Render(n.Content())
		if err != nil {
			util.HttpErrorf(w, http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
			return
		}

		err = render(w, data[getNoteData]{
			Title:      n.Title,
			CurrentUrl: r.URL.Path,
			Data:       getNoteData{Notes: notes, Note: n, Html: html},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}
