package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/web/service"
)

//go:embed static
var static embed.FS

// NotesProvider supplies the notes the server renders.
type NotesProvider interface {
	Notes() []note.Note
}

// MarkdownRenderer converts the markdown source of a note into HTML.
type MarkdownRenderer interface {
	Render(input string) (template.HTML, error)
}

type ServerArgs struct {
	Port  string
	Notes NotesProvider
}

func NewServer(args ServerArgs) http.Server {
	var markdownRenderer MarkdownRenderer = service.NewMarkdown()

	router := http.NewServeMux()

	router.Handle("/{$}", http.RedirectHandler("/notes/", http.StatusTemporaryRedirect))
	router.HandleFunc("GET /notes/{$}", getNotes(args.Notes))
	router.HandleFunc("GET /notes/{slug}/{$}", getNote(args.Notes, markdownRenderer))

	router.HandleFunc("GET /search/{$}", getSearch(args.Notes, markdownRenderer))

	router.HandleFunc("GET /tags/{tag}/{$}", getTag())

	router.Handle("/static/", http.FileServerFS(static))

	return http.Server{
		Addr:    fmt.Sprintf(":%v", args.Port),
		Handler: router,
	}
}
