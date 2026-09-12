package note

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/forPelevin/gomoji"
)

func last[T ~[]I, I any](slice T) *I {
	if len(slice) == 0 {
		return nil
	}

	return &slice[len(slice)-1]
}

var tagsExp = regexp.MustCompile(`(?:^|\s)#(\w+)`)
var dateTitleExp = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)

func cleanupTitle(title string) string {
	title = gomoji.RemoveEmojis(title)
	title = tagsExp.ReplaceAllString(title, "")
	return strings.TrimSpace(title)
}

func newRootNote(title string) Note {
	n := New(cleanupTitle(title))

	if !dateTitleExp.MatchString(n.Title) {
		return n
	}

	date, err := time.Parse("02.01.2006", n.Title)
	if err != nil {
		return n
	}

	n.Kind = KindDaily
	n.Date = date

	return n
}

func isFence(line string) bool {
	rest, found := strings.CutPrefix(line, "```")
	return found && !strings.HasPrefix(rest, "`")
}

// Parse parses an input into a structured list of notes. The input value
// follows an opinionated subset of markdown, with the following conventions:
//
//   - the file must start with a level one heading (or be empty)
//   - a level 1 heading indicates the start of a new note, with all content
//     until the next level 1 heading considered part of that note, and the
//     content of the heading being the title of the note.
//   - if a level 1 heading is exactly a date in the format of DD.MM.YYYY, the
//     note is considered a "daily note", and will have the Note.Date set to
//     that date. A date anywhere else in the heading has no special meaning.
//   - level 2 headings in daily notes will be added to the children of that
//     note, inheriting its date. Children never have children of their own. In
//     notes without a date, the level 2 heading has no special significance and
//     no child notes will be created.
//   - notes can be tagged. A tag starts with a "#", followed by letters,
//     numbers, and underscores (word characters)
//   - for code blocks, only fenced code blocks are supported. A fence is
//     exactly 3 backticks at the beginning of a line, optionally followed by an
//     info string, which is ignored. Lines starting with 4 or more backticks
//     are not fences, so they can be used for nesting inside a block. Tilde
//     fences and code blocks by indentation are not supported.
func Parse(input io.Reader) ([]Note, error) {
	scanner := bufio.NewScanner(input)
	var notes []Note
	var root *Note
	var current *Note

	var inFencedBlock bool = false

	for scanner.Scan() {
		line := scanner.Text()

		shouldParse := !inFencedBlock

		if shouldParse {
			// Level 1 heading = new note
			if title, found := strings.CutPrefix(line, "# "); found {
				notes = append(notes, newRootNote(title))
				root = last(notes)
				current = root
			} else if current == nil {
				return nil, fmt.Errorf("invalid format: file must start with a heading")
			}

			// Level 2 heading =
			// - if the root note is a daily note, create it as a child note with
			// 	 that date
			// - otherwise ignore
			if childTitle, found := strings.CutPrefix(line, "## "); found && root != nil && root.IsDailyNote() {
				childNote := New(cleanupTitle(childTitle))
				childNote.Kind = KindChild
				childNote.Date = root.Date

				root.Children = append(root.Children, childNote)
				current = last(root.Children)
			}

			// Parse tags
			tags := tagsExp.FindAllStringSubmatch(line, -1)
			for _, tag := range tags {
				current.Tags.Add(NewTag(tag[1]))
			}
		}

		if isFence(line) {
			inFencedBlock = !inFencedBlock
		}

		// Extract first emoji for icon
		if current.Icon == "" {
			if emojis := gomoji.FindAll(line); len(emojis) > 0 {
				current.Icon = emojis[0].Character
			}
		}

		// Plain text note value
		current.Raw += line + "\n"
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if inFencedBlock {
		return nil, errors.New("invalid notes file content, fenced code block was not closed")
	}

	return notes, nil
}

// ParseFile reads the file at path and parses it into a structured list of
// notes. See Parse for details on the expected file format.
func ParseFile(path string) ([]Note, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return Parse(file)
}
