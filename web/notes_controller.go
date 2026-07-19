package web

import (
	"html/template"
	"net/http"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
	"github.com/andreasphil/one/web/adapter"
)

func getNotes(notes adapter.NotesProvider) http.HandlerFunc {
	type getNotesData struct {
		Notes []lib.Note
	}

	render := newRenderFunc[getNotesData]("get_notes.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := notes.Notes()

		err := render(w, data[getNotesData]{
			Title:      "Notes",
			CurrentUrl: r.URL.Path,
			Data:       getNotesData{Notes: notes},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
		}
	}
}

func getNote(notes adapter.NotesProvider, renderer adapter.MarkdownRenderer) http.HandlerFunc {
	type getNoteData struct {
		Note lib.Note
		Html template.HTML
	}

	render := newRenderFunc[getNoteData]("get_note.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := notes.Notes()
		slug := r.PathValue("slug")

		note, found := lib.GetRecursive(notes, slug)
		if !found {
			util.HttpErrorf(w, http.StatusNotFound, "note %v not found", slug)
			return
		}

		html, err := renderer.Render(note.Content())
		if err != nil {
			util.HttpErrorf(w, http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
		}

		err = render(w, data[getNoteData]{
			Title:      note.Title,
			CurrentUrl: r.URL.Path,
			Data:       getNoteData{Note: note, Html: html},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}
