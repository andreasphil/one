package cli

import (
	"fmt"
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
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	util.Infof("parsed %v notes", len(notes))

	notes, didSort := lib.Sort(notes)
	if !didSort {
		util.Infof("notes already sorted")
	} else if !args.check {
		util.Infof("notes need sorting")
	}

	if args.check {
		util.Warnf("check only. no changes have been made")
		if didSort {
			return fmt.Errorf("notes need sorting")
		}
		return nil
	}

	output := args.output
	if len(output) == 0 {
		output = args.input
	}

	err = util.WriteTextFile(lib.ToString(notes), output, 0644)
	if err != nil {
		return fmt.Errorf("could not write to %v, %v", output, err)
	}

	util.Okf("sorted")
	return nil
}
