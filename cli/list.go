package cli

import (
	"fmt"
	"io"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
)

type listCmdInit struct {
	input string
}

func list(args listCmdInit, stdout io.Writer, _ io.Writer) error {
	notes, err := lib.ParseFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	util.Infof("parsed %v notes", len(notes))

	fmt.Fprintf(stdout, "%v notes\n", len(notes))
	printTree(stdout, notes, "")

	return nil
}

func printTree(w io.Writer, notes []lib.Note, prefix string) {
	for i, note := range notes {
		last := i == len(notes)-1

		connector := "├─ "
		childPrefix := prefix + "│  "
		if last {
			connector = "└─ "
			childPrefix = prefix + "   "
		}

		fmt.Fprintf(w, "%v%v%v\n", prefix, connector, note.Title)

		if len(note.Children) > 0 {
			printTree(w, note.Children, childPrefix)
		}
	}
}
