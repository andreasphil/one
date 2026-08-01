package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("no command specified")
	}

	params := args[1:]

	switch args[0] {
	case "help", "--help", "-help":
		return nil

	case "list", "ls":
		listFlags := flag.NewFlagSet("list", flag.ExitOnError)
		listInput := listFlags.String("input", "one.md", "file to read")
		listFlags.Parse(params)

		return list(listCmdInit{input: *listInput}, stdout, stderr)

	case "sort":
		sortFlags := flag.NewFlagSet("sort", flag.ExitOnError)
		sortInput := sortFlags.String("input", "one.md", "file to read")
		sortOutput := sortFlags.String("output", "", "file to write to. writes to input if not specified")
		sortCheck := sortFlags.Bool("check", false, "if set, only reports if the file needs sorting without writing any changes")
		sortFlags.Parse(params)

		return sort(sortCmdInit{input: *sortInput, output: *sortOutput, check: *sortCheck}, stdout, stderr)

	case "lint":
		lintFlags := flag.NewFlagSet("lint", flag.ExitOnError)
		lintInput := lintFlags.String("input", "one.md", "file to read")
		lintFlags.Parse(params)

		return lint(lintCmdInit{input: *lintInput}, stdout, stderr)

	case "format", "fmt":
		formatFlags := flag.NewFlagSet("format", flag.ExitOnError)
		formatInput := formatFlags.String("input", "one.md", "file to read")
		formatOutput := formatFlags.String("output", "", "file to write to. writes to input if not specified")
		formatCheck := formatFlags.Bool("check", false, "if set, only reports if the file needs formatting without writing any changes")
		formatFlags.Parse(params)

		return format(formatCmdInit{input: *formatInput, output: *formatOutput, check: *formatCheck}, stdout, stderr)

	case "web":
		webFlags := flag.NewFlagSet("web", flag.ExitOnError)
		webInput := webFlags.String("input", "one.md", "file to read")
		webPort := webFlags.String("port", "8080", "port to serve on")
		webFlags.Parse(params)

		return serve(webCmdInit{input: *webInput, port: *webPort}, stdout, stderr)

	default:
		return fmt.Errorf("unknown command: %v", args[0])
	}
}
