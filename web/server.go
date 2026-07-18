package web

import (
	"embed"
	"fmt"
	"net/http"
)

//go:embed static
var static embed.FS

type ServerInit struct {
	Port string
}

func NewServer(init ServerInit) http.Server {
	router := http.NewServeMux()

	router.Handle("/{$}", http.RedirectHandler("/notes/", http.StatusTemporaryRedirect))
	router.HandleFunc("GET /notes/{$}", getNotes())

	router.Handle("/static/", http.FileServerFS(static))

	return http.Server{
		Addr:    fmt.Sprintf(":%v", init.Port),
		Handler: router,
	}
}
