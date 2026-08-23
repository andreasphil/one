package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
	"github.com/andreasphil/one/web"
)

type staticNotesProvider []note.Note

func (s staticNotesProvider) Notes() []note.Note {
	return s
}

type webArgs struct {
	input string
	port  string
}

func serve(args webArgs, _ io.Writer, stderr io.Writer) error {
	notes, err := note.ParseFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	util.Infof(stderr, "parsed %v notes", len(notes))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server := web.NewServer(web.ServerArgs{
		Port:   args.port,
		Notes:  staticNotesProvider(notes),
		Errors: stderr,
	})

	errChan := make(chan error, 1)

	go func() {
		util.Infof(stderr, "serving at http://localhost:%v", args.port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("server exited with an error: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		server.Shutdown(ctx)
	case err := <-errChan:
		return err
	}

	fmt.Fprintln(stderr)
	util.Infof(stderr, "bye :)")

	return nil
}
