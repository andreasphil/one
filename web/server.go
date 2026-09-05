// Package web serves the notes over HTTP.
package web

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/web/service"
)

//go:embed static
var staticFS embed.FS

// NotesProvider supplies the notes the server renders.
type NotesProvider interface {
	Notes() []note.Note
}

// MarkdownRenderer converts the markdown source of a note into HTML.
type MarkdownRenderer interface {
	Render(input string) (template.HTML, error)
}

// ServerArgs configures a server.
type ServerArgs struct {
	// Port is the TCP port to listen on.
	Port string
	// Notes supplies the notes the server renders.
	Notes NotesProvider
	// Errors is where request errors are logged. Defaults to io.Discard.
	Errors io.Writer
}

// NewServer creates a server with the routes, templates and static files of
// the notes UI. The returned server is not started.
func NewServer(args ServerArgs) http.Server {
	var markdownRenderer MarkdownRenderer = service.NewMarkdown(func(target string) (string, bool) {
		return note.ResolveSlug(args.Notes.Notes(), target)
	})

	errw := args.Errors
	if errw == nil {
		errw = io.Discard
	}

	router := http.NewServeMux()

	router.Handle("/{$}", http.RedirectHandler("/notes/", http.StatusTemporaryRedirect))
	router.HandleFunc("GET /notes/{$}", handle(errw, getNotes(args.Notes)))
	router.HandleFunc("GET /notes/{slug}/{$}", handle(errw, getNote(args.Notes, markdownRenderer)))

	router.HandleFunc("GET /search/{$}", handle(errw, getSearch(args.Notes, markdownRenderer)))

	router.HandleFunc("GET /tags/{tag}/{$}", handle(errw, getTag()))

	router.Handle("/static/", http.FileServerFS(staticFS))

	return http.Server{
		Addr:    fmt.Sprintf(":%v", args.Port),
		Handler: router,
	}
}
