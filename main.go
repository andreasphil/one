package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/andreasphil/one/service"
	"github.com/andreasphil/one/util"
	"github.com/andreasphil/one/web"
)

//go:embed static
var static embed.FS

type Config struct {
	Port string
	Path string
}

func serve(config Config) {
	context, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	notesService, err := service.NewNotesService(config.Path)
	if err != nil {
		util.Fatalf("failed to parse notes at path %v: %v", config.Path, err)
	}

	markdownService := service.NewMarkdownService()

	server := web.NewServer(web.ServerInit{
		Static:            static,
		Port:              config.Port,
		NotesLoader:       notesService,
		NoteLoader:        notesService,
		NotesByDateFinder: notesService,
		MarkdownRenderer:  &markdownService,
	})

	go func() {
		util.Infof("serving at http://localhost:%v", config.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Fatalf("server exited with an error: %v", err)
		}
	}()

	<-context.Done()
	server.Shutdown(context)

	fmt.Println() // Clear the ^C from terminal
	util.Infof("bye :)")
}

func main() {
	serve(Config{Port: "8080", Path: "test.md"})
}
