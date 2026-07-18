package main

import (
	"os"

	"github.com/andreasphil/one/cli"
	"github.com/andreasphil/one/util"
)

func main() {
	code := 0
	if err := cli.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		util.Errorf("%v", err)
		code = 1
	}

	os.Exit(code)
}
