package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

type snapshotArgs struct {
	input string
	check bool
	tidy  bool
}

func execGit(args ...string) (string, error) {
	execPath, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("git is not installed or not in PATH: %v", err)
	}

	cmd := exec.Command(execPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %v failed: %v: %v", args[0], err, strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
}

func tidySnapshot(args snapshotArgs, stderr io.Writer) error {
	content, err := os.ReadFile(args.input)
	if err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	notes, err := note.Parse(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to parse notes from %v, %v", args.input, err)
	}

	tidied := string(content)

	if notes, didSort := note.Sort(notes); !didSort {
		util.Infof(stderr, "notes already sorted")
	} else {
		tidied = note.String(notes)
		if args.check {
			util.Infof(stderr, "notes would be sorted")
		} else {
			util.Successf(stderr, "sorted")
		}
	}

	formatted, err := execFormatter([]byte(tidied), args.input)
	if err != nil {
		return err
	}

	if formatted == tidied {
		util.Infof(stderr, "notes already formatted")
	} else if args.check {
		util.Infof(stderr, "notes would be formatted")
	} else {
		util.Successf(stderr, "formatted")
	}

	if args.check || formatted == string(content) {
		return nil
	}

	if err := util.WriteTextFile(formatted, args.input, 0644); err != nil {
		return fmt.Errorf("could not write to %v, %v", args.input, err)
	}

	return nil
}

func snapshot(args snapshotArgs, _ io.Writer, stderr io.Writer) error {
	if inside, err := execGit("rev-parse", "--is-inside-work-tree"); err != nil || inside != "true" {
		return fmt.Errorf("current directory is not inside a git repository")
	}

	if _, err := os.Stat(args.input); err != nil {
		return fmt.Errorf("failed to read notes from %v, %v", args.input, err)
	}

	status, err := execGit("status", "--porcelain", "--", args.input)
	if err != nil {
		return err
	}

	if len(status) == 0 {
		util.Warnf(stderr, "no changes since last snapshot")
		return fmt.Errorf("no snapshot created")
	}

	otherStaged, err := execGit("diff", "--cached", "--name-only", "--", ":/", ":(exclude,literal)"+args.input)
	if err != nil {
		return err
	}

	if len(otherStaged) > 0 {
		return fmt.Errorf("other files are currently staged for commit. Please unstage them first.")
	}

	if args.tidy {
		if err := tidySnapshot(args, stderr); err != nil {
			return err
		}
	}

	if args.check {
		util.Infof(stderr, "snapshot would be created")
		util.Warnf(stderr, "check only. no changes have been made")
		return nil
	}

	if _, err := execGit("add", "--", args.input); err != nil {
		return err
	}

	message := "snapshot: " + time.Now().Format("2006-01-02 15:04")
	if _, err := execGit("commit", "--message", message, "--", args.input); err != nil {
		return err
	}

	util.Successf(stderr, "created %q", message)
	return nil
}
