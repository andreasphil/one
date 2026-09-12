package web

import (
	"embed"
	"io"
	"net/http"

	"github.com/andreasphil/one/lib/note"
)

//go:embed static
var staticFS embed.FS

// NotesProvider supplies the notes the web interface renders.
type NotesProvider interface {
	Notes() []note.Note
}

// RouterArgs configures the router.
type RouterArgs struct {
	Notes  NotesProvider
	Errors io.Writer
}

// NewRouter creates a handler with the routes, templates and static files of
// the notes UI.
func NewRouter(args RouterArgs) http.Handler {
	markdownRenderer := newMarkdownRenderer(func(target string) (string, bool) {
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
