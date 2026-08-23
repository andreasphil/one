package note

import (
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/andreasphil/one/util"
)

var dateTitleExp = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
var normalizeExp = regexp.MustCompile(`[^\wäöüß]+`)

// Note represents a single note parsed from a notes file. Daily notes (notes
// whose title is a date in the format of DD.MM.YYYY) may have children,
// which represent the level 2 headings within that daily note.
type Note struct {
	Title    string
	Icon     string
	Date     time.Time
	Tags     util.Set[string]
	Children []Note
	Raw      string
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
	titleExp := regexp.MustCompile(`^#{1,2}\s+.+\n`)
	content := titleExp.ReplaceAllString(n.Raw, "")
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

// Excerpt returns up to maxWords words from the note's content, with whitespace
// normalized. The second return value reports whether the content was actually
// longer than maxWords, i.e. whether the excerpt is a partial view of the
// content.
func (n Note) Excerpt(maxWords int) (string, bool) {
	words := strings.Split(n.Content(), " ")
	partial := strings.Join(words[0:int(math.Min(float64(len(words)), float64(maxWords)))], " ")

	whitespaceExp := regexp.MustCompile(`\s+`)
	partial = whitespaceExp.ReplaceAllString(partial, " ")

	return partial, len(words) >= maxWords
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
