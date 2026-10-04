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
   format, fmt   Format notes
   lint          Check notes for issues
   list, ls      List notes and their structure
   merge         Merge notes from another file
   snapshot      Commit notes to git
   sort          Sort notes
   tags          List tags
   web           Serve notes over HTTP

 Run 'one <command> --help' for the flags of a specific command.
`)
}

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

	case "tags":
		tagsFlags := flag.NewFlagSet("tag", flag.ExitOnError)
		tagsInput := tagsFlags.String("input", "one.md", "file to read")
		tagsFlags.Parse(params)

		return tags(tagsArgs{input: *tagsInput}, stdout, stderr)

	case "sort":
		sortFlags := flag.NewFlagSet("sort", flag.ExitOnError)
		sortInput := sortFlags.String("input", "one.md", "file to read")
		sortOutput := sortFlags.String("output", "", "file to write to. writes to input if not specified")
		sortCheck := sortFlags.Bool("check", false, "if set, only reports if the file needs sorting without making any changes")
		sortFlags.Parse(params)

		return sort(sortArgs{input: *sortInput, output: *sortOutput, check: *sortCheck}, stdout, stderr)

	case "merge":
		mergeFlags := flag.NewFlagSet("merge", flag.ExitOnError)
		mergeInput := mergeFlags.String("input", "one.md", "file to merge into")
		mergeOutput := mergeFlags.String("output", "", "file to write to. writes to input if not specified")

		mergeFlags.Usage = func() {
			fmt.Fprintf(mergeFlags.Output(), "Usage of merge: one merge [flags] <file>\n")
			mergeFlags.PrintDefaults()
		}

		mergeFlags.Parse(params)

		if mergeFlags.NArg() == 0 {
			return fmt.Errorf("no file to merge specified")
		} else if mergeFlags.NArg() > 1 {
			return fmt.Errorf("expected exactly one file to merge")
		}

		return merge(mergeArgs{input: *mergeInput, output: *mergeOutput, file: mergeFlags.Arg(0)}, stdout, stderr)

	case "lint":
		lintFlags := flag.NewFlagSet("lint", flag.ExitOnError)
		lintInput := lintFlags.String("input", "one.md", "file to read")
		lintFlags.Parse(params)

		return lint(lintArgs{input: *lintInput}, stdout, stderr)

	case "format", "fmt":
		formatFlags := flag.NewFlagSet("format", flag.ExitOnError)
		formatInput := formatFlags.String("input", "one.md", "file to read")
		formatOutput := formatFlags.String("output", "", "file to write to. writes to input if not specified")
		formatCheck := formatFlags.Bool("check", false, "if set, only reports if the file needs formatting without making any changes")
		formatFlags.Parse(params)

		return format(formatArgs{input: *formatInput, output: *formatOutput, check: *formatCheck}, stdout, stderr)

	case "snapshot":
		snapshotFlags := flag.NewFlagSet("snapshot", flag.ExitOnError)
		snapshotInput := snapshotFlags.String("input", "one.md", "file to read")
		snapshotCheck := snapshotFlags.Bool("check", false, "if set, only reports if a snapshot would be created without making any changes")
		snapshotTidy := snapshotFlags.Bool("tidy", false, "if set, sorts and formats notes before creating the snapshot")
		snapshotFlags.Parse(params)

		return snapshot(snapshotArgs{input: *snapshotInput, check: *snapshotCheck, tidy: *snapshotTidy}, stdout, stderr)

	case "web":
		webFlags := flag.NewFlagSet("web", flag.ExitOnError)
		webInput := webFlags.String("input", "one.md", "file to read")
		webPort := webFlags.String("port", "1111", "port to serve on")
		webFlags.Parse(params)

		return serve(webArgs{input: *webInput, port: *webPort}, stdout, stderr)

	default:
		return fmt.Errorf("unknown command: %v (run 'one help' for usage)", args[0])
	}
}
