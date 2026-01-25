package cli

import (
	"flag"

	"github.com/andreasphil/one/util"
)

type Args struct {
	Filename string
	Port     string
}

func ReadArgs() *Args {
	port := flag.String("port", "8080", "the port for serving the application")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		util.Fatalf("filename is required")
	} else if len(args) > 1 {
		util.Warnf("ignored additional arguments: %v", args[1:])
	}

	return &Args{
		Filename: flag.Arg(0),
		Port:     *port,
	}
}
