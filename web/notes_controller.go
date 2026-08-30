package web

import (
	"html/template"
	"io"
	"net/http"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

func getNotes(errw io.Writer, provider NotesProvider) http.HandlerFunc {
	type getNotesData struct {
		Notes []note.Note
	}

	render := newRenderFunc[getNotesData]("get_notes.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := provider.Notes()

		err := render(w, data[getNotesData]{
			Title:      "Notes",
			CurrentURL: r.URL.Path,
			Data:       getNotesData{Notes: notes},
		})

		if err != nil {
			util.HTTPErrorf(errw, w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}

func getNote(errw io.Writer, provider NotesProvider, renderer MarkdownRenderer) http.HandlerFunc {
	type getNoteData struct {
		Notes []note.Note
		Note  note.Note
		HTML  template.HTML
	}

	render := newRenderFunc[getNoteData]("get_note.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := provider.Notes()
		slug := r.PathValue("slug")

		n, found := note.FindBySlug(notes, slug)
		if !found {
			util.HTTPErrorf(errw, w, http.StatusNotFound, "note %v not found", slug)
			return
		}

		html, err := renderer.Render(n.Content())
		if err != nil {
			util.HTTPErrorf(errw, w, http.StatusUnprocessableEntity, "failed to render note to html: %v", err)
			return
		}

		err = render(w, data[getNoteData]{
			Title:      n.Title,
			CurrentURL: r.URL.Path,
			Data:       getNoteData{Notes: notes, Note: n, HTML: html},
		})

		if err != nil {
			util.HTTPErrorf(errw, w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}
