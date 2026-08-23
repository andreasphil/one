// Package note parses, inspects and serializes the notes in a notes file.
package note

import (
	"regexp"
	"strings"
	"time"

	"github.com/andreasphil/one/util"
)

var dateTitleExp = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
var normalizeExp = regexp.MustCompile(`[^\wäöüß]+`)
var titleHeadingExp = regexp.MustCompile(`^#{1,2}\s+.+\n`)

// Note represents a single note parsed from a notes file. Daily notes (notes
// whose title is a date in the format of DD.MM.YYYY) may have children,
// which represent the level 2 headings within that daily note.
type Note struct {
	// Title is the note's heading, with tags and emoji removed.
	Title string
	// Icon is the first emoji occurring in the note, if any.
	Icon string
	// Date is set for daily notes, and inherited by their children.
	Date time.Time
	// Tags are the tags occurring anywhere in the note, each including its
	// leading "#".
	Tags util.Set[string]
	// Children are the notes formed by the level 2 headings of a daily note.
	Children []Note
	// Raw is the note's own markdown source, including its heading but
	// excluding the source of any children.
	Raw string
}

// New creates a new, empty Note with the given title.
func New(title string) Note {
	return Note{
		Title: title,
		Tags:  util.NewSet[string](),
	}
}

// Slug returns a unique, URL-friendly identifier for the note, derived from
// its date (if any) and title.
func (n Note) Slug() string {
	slug := strings.Builder{}

	if !n.Date.IsZero() {
		slug.WriteString(n.Date.Format("2006-01-02"))
	}

	if !n.IsDailyNote() {
		normalized := strings.ToLower(n.Title)
		normalized = normalizeExp.ReplaceAllString(normalized, "-")
		normalized = strings.Trim(normalized, "-")

		if slug.Len() > 0 {
			slug.WriteString("-")
		}

		slug.WriteString(normalized)
	}

	return slug.String()
}

// Content returns the note's raw content with the title heading removed and
// surrounding whitespace trimmed.
func (n Note) Content() string {
	content := titleHeadingExp.ReplaceAllString(n.Raw, "")
	return strings.TrimSpace(content)
}

// IsEmpty reports whether the note has no content besides its title heading.
func (n Note) IsEmpty() bool {
	return len(n.Content()) == 0
}

// IsDailyNote reports whether the note has a date and its title matches the
// daily note format of DD.MM.YYYY.
func (n Note) IsDailyNote() bool {
	return !n.Date.IsZero() && dateTitleExp.MatchString(n.Title)
}

// String returns the note's raw markdown source, including that of any
// children.
func (n Note) String() string {
	var raw strings.Builder
	raw.WriteString(n.Raw)

	for _, note := range n.Children {
		raw.WriteString(note.String())
	}

	return raw.String()
}
