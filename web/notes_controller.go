package web

import (
	"html/template"
	"net/http"

	"github.com/andreasphil/one/adapter"
	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
)

func getNotes(loader adapter.NotesLoader) http.HandlerFunc {
	type notesPageData struct {
		Notes []lib.Note
	}

	render := newRenderFunc[notesPageData]("get_notes.html")

	return func(w http.ResponseWriter, r *http.Request) {
		notes := loader.LoadNotes()

		err := render(w, data[notesPageData]{
			Title:      "Notes",
			CurrentUrl: r.URL.Path,
			Data:       notesPageData{Notes: notes},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}

func getNote(noteLoader adapter.NoteLoader, notesLoader adapter.NotesLoader, renderer adapter.MarkdownRenderer) http.HandlerFunc {
	type notePageData struct {
		Notes []lib.Note
		Note  lib.Note
		Html  template.HTML
	}

	render := newRenderFunc[notePageData]("get_note.html")

	return func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")

		notes := notesLoader.LoadNotes()

		note, found := noteLoader.LoadNote(slug)
		if !found {
			util.HttpErrorf(w, http.StatusNotFound, "note %v not found", slug)
			return
		}

		html, err := renderer.RenderMarkdown(note.Raw)
		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render note to html: %v", err)
		}

		err = render(w, data[notePageData]{
			Title:      note.Title,
			CurrentUrl: r.URL.Path,
			Data: notePageData{
				Notes: notes,
				Note:  note,
				Html:  template.HTML(html),
			},
		})

		if err != nil {
			util.HttpErrorf(w, http.StatusInternalServerError, "failed to render page template: %v", err)
			return
		}
	}
}
