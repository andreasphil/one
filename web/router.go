// Package web serves the notes over HTTP.
package web

import (
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/web/service"
)

//go:embed static
var staticFS embed.FS

// NotesProvider supplies the notes the web interface renders.
type NotesProvider interface {
	Notes() []note.Note
}

// MarkdownRenderer converts the markdown source of a note into HTML.
type MarkdownRenderer interface {
	Render(input string) (template.HTML, error)
}

// RouterArgs configures a handler.
type RouterArgs struct {
	// Notes supplies the notes the web interface renders.
	Notes NotesProvider
	// Errors is where request errors are logged. Defaults to io.Discard.
	Errors io.Writer
}

// NewRouter creates a handler with the routes, templates and static files of
// the notes UI.
func NewRouter(args RouterArgs) http.Handler {
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

	return router
}
