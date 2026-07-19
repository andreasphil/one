package web

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/andreasphil/one/web/adapter"
)

//go:embed static
var static embed.FS

type ServerInit struct {
	Port  string
	Notes adapter.NotesProvider
}

func NewServer(init ServerInit) http.Server {
	router := http.NewServeMux()

	router.Handle("/{$}", http.RedirectHandler("/notes/", http.StatusTemporaryRedirect))
	router.HandleFunc("GET /notes/{$}", getNotes(init.Notes))

	router.Handle("/static/", http.FileServerFS(static))

	return http.Server{
		Addr:    fmt.Sprintf(":%v", init.Port),
		Handler: router,
	}
}
