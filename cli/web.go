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
		util.Errorf("failed to read notes from %v, %v", args.input, err)
		return err
	}

	util.Infof("parsed %v notes", len(notes))

	context, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	server := web.NewServer(web.ServerInit{
		Port: args.port,
	})

	go func() {
		util.Infof("serving at http://localhost:%v", args.port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Errorf("server exited with an error: %v", err)
		}
	}()

	<-context.Done()
	server.Shutdown(context)

	fmt.Println()
	util.Infof("bye :)")

	return nil
}
