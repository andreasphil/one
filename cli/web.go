package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/andreasphil/one/util"
	"github.com/andreasphil/one/web"
)

type webArgs struct {
	input string
	port  string
}

func serve(args webArgs, _ io.Writer, stderr io.Writer) error {
	provider, err := newFileNotesProvider(args.input, stderr)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	router := web.NewRouter(web.RouterArgs{
		Notes:  provider,
		Errors: stderr,
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%v", args.port),
		Handler: router,
	}

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
