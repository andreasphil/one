package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
	"github.com/andreasphil/one/web"
)

type webCmdInit struct {
	input string
	port  string
}

func serve(args webCmdInit, stdout io.Writer, _ io.Writer) error {
	notes, err := lib.ParseFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	util.Infof("parsed %v notes", len(notes))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server := web.NewServer(web.ServerInit{
		Port: args.port,
	})

	errChan := make(chan error, 1)

	go func() {
		util.Infof("serving at http://localhost:%v", args.port)

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

	fmt.Println()
	util.Infof("bye :)")

	return nil
}
