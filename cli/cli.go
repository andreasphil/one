// Package cli implements the commands of the one executable.
package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/andreasphil/one/util"
)

func usage(w io.Writer) {
	fmt.Fprint(w, ` Usage: one <command> [flags]

 Commands:
   list, ls      List notes and their structure
   sort          Sort notes
   lint          Check notes for issues
   format, fmt   Format notes
   web           Serve notes over HTTP

 Run 'one <command> --help' for the flags of a specific command.
`)
}

// Run executes the command named by the first of args, with the rest passed to
// that command as flags. Data is written to stdout, logs and diagnostics to
// stderr. It returns an error if the command is unknown or if it failed.
func Run(args []string, stdout io.Writer, stderr io.Writer) error {
	util.Banner(stderr, "one")

	if len(args) == 0 {
		return fmt.Errorf("no command specified (run 'one help' for usage)")
	}

	params := args[1:]

	switch args[0] {
	case "help", "--help", "-help":
		usage(stdout)
		return nil

	case "list", "ls":
		listFlags := flag.NewFlagSet("list", flag.ExitOnError)
		listInput := listFlags.String("input", "one.md", "file to read")
		listFlags.Parse(params)

		return list(listArgs{input: *listInput}, stdout, stderr)

	case "sort":
		sortFlags := flag.NewFlagSet("sort", flag.ExitOnError)
		sortInput := sortFlags.String("input", "one.md", "file to read")
		sortOutput := sortFlags.String("output", "", "file to write to. writes to input if not specified")
		sortCheck := sortFlags.Bool("check", false, "if set, only reports if the file needs sorting without writing any changes")
		sortFlags.Parse(params)

		return sort(sortArgs{input: *sortInput, output: *sortOutput, check: *sortCheck}, stdout, stderr)

	case "lint":
		lintFlags := flag.NewFlagSet("lint", flag.ExitOnError)
		lintInput := lintFlags.String("input", "one.md", "file to read")
		lintFlags.Parse(params)

		return lint(lintArgs{input: *lintInput}, stdout, stderr)

	case "format", "fmt":
		formatFlags := flag.NewFlagSet("format", flag.ExitOnError)
		formatInput := formatFlags.String("input", "one.md", "file to read")
		formatOutput := formatFlags.String("output", "", "file to write to. writes to input if not specified")
		formatCheck := formatFlags.Bool("check", false, "if set, only reports if the file needs formatting without writing any changes")
		formatFlags.Parse(params)

		return format(formatArgs{input: *formatInput, output: *formatOutput, check: *formatCheck}, stdout, stderr)

	case "web":
		webFlags := flag.NewFlagSet("web", flag.ExitOnError)
		webInput := webFlags.String("input", "one.md", "file to read")
		webPort := webFlags.String("port", "8080", "port to serve on")
		webFlags.Parse(params)

		return serve(webArgs{input: *webInput, port: *webPort}, stdout, stderr)

	default:
		return fmt.Errorf("unknown command: %v (run 'one help' for usage)", args[0])
	}
}
