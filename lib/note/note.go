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

var codeFenceExp = regexp.MustCompile("(?m)^[ \t]*(?:```|~~~).*$")
var ruleExp = regexp.MustCompile(`(?m)^[ \t]*(?:[-*_][ \t]*){3,}$`)
var tableDividerExp = regexp.MustCompile(`(?m)^[ \t]*[|:-][ \t:|-]*$`)
var tablePipeExp = regexp.MustCompile(`\|`)
var wikilinkExp = regexp.MustCompile(`\[\[([^\]\n]+)\]\]`)
var linkExp = regexp.MustCompile(`!?\[([^\]]*)\]\([^)]*\)`)
var autolinkExp = regexp.MustCompile(`<(https?://[^>]+)>`)
var blockMarkerExp = regexp.MustCompile(`(?m)^[ \t]*(?:(?:#{1,6}|>|[-*+]|\d+\.)(?:[ \t]+|$))+(?:\[[ xX]\][ \t]+)?`)
var underscoreEmphasisExp = regexp.MustCompile(`\b_{1,2}([^_\n]+)_{1,2}\b`)
var inlineMarkerExp = regexp.MustCompile("[*`]+|~~")

const excerptWords = 30

// Tag is a label attached to a note. Its value includes the leading "#".
type Tag string

// NewTag returns the Tag with the given name. The name may be given with or
// without a leading "#".
func NewTag(name string) Tag {
	return Tag("#" + strings.TrimPrefix(name, "#"))
}

// Name returns the tag without its leading "#".
func (t Tag) Name() string {
	return strings.TrimPrefix(string(t), "#")
}

// String returns the tag including its leading "#".
func (t Tag) String() string {
	return string(t)
}

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
	// Tags are the tags occurring anywhere in the note.
	Tags util.Set[Tag]
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
		Tags:  util.NewSet[Tag](),
	}
}

// Slug converts input into a lowercase, URL-friendly identifier.
func Slug(input string) string {
	slug := strings.ToLower(input)
	slug = normalizeExp.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

// Slug returns a unique, URL-friendly identifier for the note, derived from
// its date (if any) and title.
func (n Note) Slug() string {
	slug := strings.Builder{}

	if !n.Date.IsZero() {
		slug.WriteString(n.Date.Format("2006-01-02"))
	}

	if !n.IsDailyNote() {
		if slug.Len() > 0 {
			slug.WriteString("-")
		}

		slug.WriteString(Slug(n.Title))
	}

	return slug.String()
}

// Content returns the note's raw content with the title heading removed and
// surrounding whitespace trimmed.
func (n Note) Content() string {
	content := titleHeadingExp.ReplaceAllString(n.Raw, "")
	return strings.TrimSpace(content)
}

// Excerpt returns the first 30 words of the note's content, with markdown
// formatting removed and whitespace normalized. Longer content is truncated
// and ends with "...". The cleanup is best effort and may leave markers of
// less common syntax intact.
func (n Note) Excerpt() string {
	text := n.Content()
	text = codeFenceExp.ReplaceAllString(text, "")
	text = ruleExp.ReplaceAllString(text, "")
	text = tableDividerExp.ReplaceAllString(text, "")
	text = tablePipeExp.ReplaceAllString(text, " ")
	text = wikilinkExp.ReplaceAllString(text, "$1")
	text = linkExp.ReplaceAllString(text, "$1")
	text = autolinkExp.ReplaceAllString(text, "$1")
	text = blockMarkerExp.ReplaceAllString(text, "")
	text = underscoreEmphasisExp.ReplaceAllString(text, "$1")
	text = inlineMarkerExp.ReplaceAllString(text, "")

	words := strings.Fields(text)
	if len(words) > excerptWords {
		return strings.Join(words[:excerptWords], " ") + "..."
	}

	return strings.Join(words, " ")
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
