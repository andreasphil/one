package cli

import (
	"io"

	"github.com/andreasphil/one/lib"
	"github.com/andreasphil/one/util"
)

type sortCmdInit struct {
	input  string
	output string
	check  bool
}

func sort(args sortCmdInit, _ io.Writer, _ io.Writer) error {
	notes, err := lib.ParseFile(args.input)
	if err != nil {
		util.Errorf("failed to read notes from %v, %v", args.input, err)
		return err
	}

	util.Infof("parsed %v notes", len(notes))

	notes, didSort := lib.Sort(notes)
	if !didSort {
		util.Infof("notes already sorted")
	} else {
		util.Infof("notes need sorting")
	}

	if args.check {
		util.Warnf("check only. no changes have been made")
		return nil
	}

	output := args.output
	if len(output) == 0 {
		output = args.input
	}

	err = util.WriteTextFile(lib.Stringify(notes), output)
	if err != nil {
		util.Errorf("could not write to %v, %v", output, err)
		return err
	}

	util.Okf("sorted")
	return nil
}
