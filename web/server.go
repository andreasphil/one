package web

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/andreasphil/one/adapter"
)

type ServerInit struct {
	Static            fs.FS
	Port              string
	NotesLoader       adapter.NotesLoader
	NoteLoader        adapter.NoteLoader
	NotesByDateFinder adapter.NotesByDateFinder
	MarkdownRenderer  adapter.MarkdownRenderer
}

func NewServer(init ServerInit) http.Server {
	router := http.NewServeMux()

	router.Handle("/{$}", http.RedirectHandler("/notes/", http.StatusTemporaryRedirect))
	router.HandleFunc("GET /notes/{$}", getNotes(init.NotesLoader))
	router.HandleFunc("GET /notes/{slug}", getNote(init.NoteLoader, init.NotesLoader, init.MarkdownRenderer))

	router.Handle("/static/", http.FileServerFS(init.Static))

	return http.Server{
		Addr:    fmt.Sprintf(":%v", init.Port),
		Handler: router,
	}
}
