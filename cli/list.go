package cli

import (
	"fmt"
	"io"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

type listCmdInit struct {
	input string
}

func list(args listCmdInit, stdout io.Writer, _ io.Writer) error {
	notes, err := note.ParseFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	util.Infof("parsed %v notes", len(notes))

	fmt.Fprintf(stdout, "%v notes\n", len(notes))
	printTree(stdout, notes, "")

	return nil
}

func printTree(w io.Writer, notes []note.Note, prefix string) {
	for i, n := range notes {
		last := i == len(notes)-1

		connector := "├─ "
		childPrefix := prefix + "│  "
		if last {
			connector = "└─ "
			childPrefix = prefix + "   "
		}

		fmt.Fprintf(w, "%v%v%v\n", prefix, connector, n.Title)

		if len(n.Children) > 0 {
			printTree(w, n.Children, childPrefix)
		}
	}
}
