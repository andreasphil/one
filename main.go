package main

import (
	"os"

	"github.com/andreasphil/one/cli"
)

func main() {
	code := 0
	if err := cli.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		code = 1
	}

	os.Exit(code)
}
