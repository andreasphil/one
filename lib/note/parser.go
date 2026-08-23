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

var tagsExp = regexp.MustCompile(`(^|\s)#\w+`)
var dateExp = regexp.MustCompile(`^# (\d{2}\.\d{2}\.\d{4})`)

func cleanupTitle(title string) string {
	title = gomoji.RemoveEmojis(title)
	title = tagsExp.ReplaceAllString(title, "")
	return strings.TrimSpace(title)
}

// isFence reports whether line opens or closes a fenced code block, which is
// the case for exactly 3 backticks at the beginning of the line. Anything
// following the backticks (an info string such as "```js") is ignored, except
// for further backticks, which make the line a non-fence.
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
//   - if a level 1 heading matches a date in the format of DD.MM.YYYY, the note
//     is considered a "daily note", and will have the Note.Date set to that date
//   - level 2 headings in daily notes will be added to the children of that
//     note. In notes without a date, the level 2 heading has no special
//     significance and no child notes will be created.
//   - notes can be tagged. A tag starts with a "#", by letters, numbers, and
//     underscores (word characters)
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

	var isFencedBlock bool = false

	for scanner.Scan() {
		line := scanner.Text()

		shouldParseDetail := !isFencedBlock

		if shouldParseDetail {
			// Level 1 heading = new note
			if title, found := strings.CutPrefix(line, "# "); found {
				notes = append(notes, New(cleanupTitle(title)))
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
				current = root

				childNote := New(cleanupTitle(childTitle))
				childNote.Date = current.Date

				current.Children = append(current.Children, childNote)
				current = last(current.Children)
			}

			// Parse tags
			tags := tagsExp.FindAllString(line, -1)
			for _, tag := range tags {
				current.Tags.Add(strings.TrimSpace(tag))
			}

			// Parse date, only in note name for now
			if match := dateExp.FindStringSubmatch(line); len(match) == 2 {
				if parsedTime, err := time.Parse("02.01.2006", match[1]); err == nil {
					current.Date = parsedTime
				}
			}
		}

		if isFence(line) {
			isFencedBlock = !isFencedBlock
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

	if isFencedBlock {
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
