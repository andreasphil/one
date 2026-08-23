package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/andreasphil/one/util"
)

type formatArgs struct {
	input  string
	output string
	check  bool
}

func execFormatter(content []byte, filename string) (string, error) {
	execPath, err := exec.LookPath("oxfmt")
	if err != nil {
		return "", fmt.Errorf("oxfmt is not installed or not in PATH: %v", err)
	}

	cmd := exec.Command(execPath, "--stdin-filepath", filename)
	cmd.Stdin = bytes.NewReader(content)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("oxfmt failed to format %v: %v: %v", filename, err, stderr.String())
	}

	return stdout.String(), nil
}

func format(args formatArgs, _ io.Writer, _ io.Writer) error {
	content, err := os.ReadFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	formatted, err := execFormatter(content, args.input)
	if err != nil {
		return err
	}

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
