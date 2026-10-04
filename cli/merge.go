package cli

import (
	"fmt"
	"io"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

type mergeArgs struct {
	input  string
	output string
	file   string
}

func merge(args mergeArgs, _ io.Writer, stderr io.Writer) error {
	incoming, err := note.ParseFile(args.file)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.file, err)
	}

	if len(incoming) == 0 {
		util.Warnf(stderr, "%v contains no notes. no changes have been made", args.file)
		return nil
	}

	util.Infof(stderr, "parsed %v notes to merge", note.Count(incoming))

	notes, err := note.ParseFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	util.Infof(stderr, "parsed %v notes", note.Count(notes))

	existingDuplicates := util.NewSetFrom(note.DuplicateSlugs(notes))
	merged, sorted := note.Merge(notes, incoming)

	for _, slug := range note.DuplicateSlugs(merged) {
		if !existingDuplicates.Has(slug) {
			util.Warnf(stderr, "duplicate slug: %v", slug)
		}
	}

	if !sorted {
		util.Warnf(stderr, "notes were not sorted, appended merged notes to the end. run 'one sort' to sort")
	}

	output := args.output
	if len(output) == 0 {
		output = args.input
	}

	err = util.WriteTextFile(note.String(merged), output, 0644)
	if err != nil {
		return fmt.Errorf("could not write to %v, %v", output, err)
	}

	util.Successf(stderr, "merged")
	return nil
}
