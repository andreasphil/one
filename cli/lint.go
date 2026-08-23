package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

type lintArgs struct {
	input string
}

func lint(args lintArgs, _ io.Writer, stderr io.Writer) error {
	hadWarnings := false

	// Parses without errors
	notes, err := note.ParseFile(args.input)
	if err != nil {
		util.Warnf(stderr, "failed to parse %v: %v", args.input, err)
		return fmt.Errorf("linting failed")
	}

	util.Infof(stderr, "parses without errors (%v notes)", len(notes))

	// Is sorted
	if _, didSort := note.Sort(notes); didSort {
		util.Warnf(stderr, "notes are not sorted")
		hadWarnings = true
	} else {
		util.Infof(stderr, "notes are sorted")
	}

	// Is formatted, would format without errors
	if content, err := os.ReadFile(args.input); err != nil {
		util.Warnf(stderr, "failed to read %v: %v", args.input, err)
		hadWarnings = true
	} else if formatted, err := execFormatter(content, args.input); err != nil {
		util.Warnf(stderr, "%v", err)
		hadWarnings = true
	} else if formatted != string(content) {
		util.Warnf(stderr, "notes are not formatted")
		hadWarnings = true
	} else {
		util.Infof(stderr, "notes are formatted")
	}

	// Has duplicate slugs
	if duplicates := note.DuplicateSlugs(notes); len(duplicates) > 0 {
		for _, slug := range duplicates {
			util.Warnf(stderr, "duplicate slug: %v", slug)
		}
		hadWarnings = true
	} else {
		util.Infof(stderr, "no duplicate slugs")
	}

	// Has empty titles
	if count := note.CountEmptyTitles(notes); count > 0 {
		util.Warnf(stderr, "%v notes have empty titles (after cleanup)", count)
		hadWarnings = true
	} else {
		util.Infof(stderr, "no empty titles")
	}

	// Has empty notes
	if empty := note.EmptyNotes(notes); len(empty) > 0 {
		for _, n := range empty {
			util.Warnf(stderr, "empty note: %v", n.Title)
		}
		hadWarnings = true
	} else {
		util.Infof(stderr, "no empty notes")
	}

	if hadWarnings {
		return fmt.Errorf("lint found issues")
	}

	util.Successf(stderr, "no issues found")
	return nil
}
