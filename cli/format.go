package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/andreasphil/one/util"
)

type formatCmdInit struct {
	input  string
	output string
	check  bool
}

func format(args formatCmdInit, _ io.Writer, _ io.Writer) error {
	execPath, err := exec.LookPath("oxfmt")
	if err != nil {
		return fmt.Errorf("oxfmt is not installed or not in PATH: %v", err)
	}

	content, err := os.ReadFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	cmd := exec.Command(execPath, "--stdin-filepath", args.input)
	cmd.Stdin = bytes.NewReader(content)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("oxfmt failed to format %v: %v: %v", args.input, err, stderr.String())
	}

	formatted := stdout.String()
	didFormat := formatted != string(content)

	if !didFormat {
		util.Infof("notes already formatted")
	} else if !args.check {
		util.Infof("notes need formatting")
	}

	if args.check {
		util.Warnf("check only. no changes have been made")
		if didFormat {
			return fmt.Errorf("notes need formatting")
		}
		return nil
	}

	output := args.output
	if len(output) == 0 {
		output = args.input
	}

	err = util.WriteTextFile(formatted, output, 0644)
	if err != nil {
		return fmt.Errorf("could not write to %v, %v", output, err)
	}

	util.Okf("formatted")
	return nil
}
