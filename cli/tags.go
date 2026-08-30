package cli

import (
	"fmt"
	"io"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

type tagsArgs struct {
	input string
}

func tags(args tagsArgs, stdout io.Writer, stderr io.Writer) error {
	notes, err := note.ParseFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	util.Infof(stderr, "parsed %v notes\n", len(notes))

	tags := note.Tags(notes)

	fmt.Fprintf(stdout, "%v tags:\n", len(tags))
	for _, tag := range tags {
		fmt.Fprintf(stdout, "%v\n", tag)
	}

	return nil
}
