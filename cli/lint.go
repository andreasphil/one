package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
)

type lintCmdInit struct {
	input string
}

func lint(args lintCmdInit, _ io.Writer, _ io.Writer) error {
	hadWarnings := false

	// Parses without errors
	notes, err := note.ParseFile(args.input)
	if err != nil {
		util.Warnf("failed to parse %v: %v", args.input, err)
		return fmt.Errorf("linting failed")
	}

	util.Infof("parses without errors (%v notes)", len(notes))

	// Is sorted
	if _, didSort := note.Sort(notes); didSort {
		util.Warnf("notes are not sorted")
		hadWarnings = true
	} else {
		util.Infof("notes are sorted")
	}

	// Is formatted, would format without errors
	if content, err := os.ReadFile(args.input); err != nil {
		util.Warnf("failed to read %v: %v", args.input, err)
		hadWarnings = true
	} else if formatted, err := execFormatter(content, args.input); err != nil {
		util.Warnf("%v", err)
		hadWarnings = true
	} else if formatted != string(content) {
		util.Warnf("notes are not formatted")
		hadWarnings = true
	} else {
		util.Infof("notes are formatted")
	}

	// Has duplicate slugs
	if duplicates := note.DuplicateSlugs(notes); len(duplicates) > 0 {
		for _, slug := range duplicates {
			util.Warnf("duplicate slug: %v", slug)
		}
		hadWarnings = true
	} else {
		util.Infof("no duplicate slugs")
	}

	// Has empty titles
	if count := note.CountEmptyTitles(notes); count > 0 {
		util.Warnf("%v notes have empty titles (after cleanup)", count)
		hadWarnings = true
	} else {
		util.Infof("no empty titles")
	}

	// Has empty notes
	if empty := note.EmptyNotes(notes); len(empty) > 0 {
		for _, n := range empty {
			util.Warnf("empty note: %v", n.Title)
		}
		hadWarnings = true
	} else {
		util.Infof("no empty notes")
	}

	if hadWarnings {
		return fmt.Errorf("lint found issues")
	}

	util.Okf("no issues found")
	return nil
}
