package cli

import (
	"fmt"
	"io"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
)

type listCmdInit struct {
	input string
}

func list(args listCmdInit, stdout io.Writer, _ io.Writer) error {
	notes, err := lib.ParseFile(args.input)
	if err != nil {
		util.Errorf("failed to read notes from %v, %v", args.input, err)
		return err
	}

	util.Infof("parsed %v notes", len(notes))

	t := tree.New()
	t.Root(fmt.Sprintf("%v notes", len(notes)))

	for _, note := range notes {
		node := tree.Root(note.Title)
		t.Child(node)

		for _, childNote := range note.Children {
			node.Child(childNote.Title)
		}
	}

	lipgloss.Fprintln(stdout, t)

	return nil
}
