// Package note parses, inspects and serializes the notes in a notes file.
package note

import (
	"regexp"
	"strings"
	"time"

	"github.com/andreasphil/one/util"
)

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

// Kind describes how a note relates to the structure of the notes file. Every
// note has exactly one kind, assigned while parsing.
type Kind int

const (
	// KindStandalone is a note without a date, not tied to a particular day.
	KindStandalone Kind = iota
	// KindDaily is a note whose title is exactly a date in the format of
	// DD.MM.YYYY.
	KindDaily
	// KindChild is a note formed by a level 2 heading inside a daily note. It
	// inherits the date of that daily note.
	KindChild
)

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
// which represent the level 2 headings within that daily note. Children never
// have children of their own, so notes are at most two levels deep.
type Note struct {
	// Title is the note's heading, with tags and emoji removed.
	Title string
	// Kind describes how the note relates to the structure of the notes file.
	Kind Kind
	// Icon is the first emoji occurring in the note, if any.
	Icon string
	// Date is set for daily notes, and inherited by their children.
	Date time.Time
	// Tags are the tags occurring anywhere in the note.
	Tags util.Set[Tag]
	// Children are the notes formed by the level 2 headings of a daily note.
	// Only daily notes have children.
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

// Slug returns a URL-friendly identifier for the note. Daily notes are
// identified by their date alone, child notes by the date of their parent
// followed by their title, and standalone notes by their title. Slugs are not
// necessarily unique, see DuplicateSlugs.
func (n Note) Slug() string {
	date := ""
	if !n.Date.IsZero() {
		date = n.Date.Format("2006-01-02")
	}

	if n.IsDailyNote() {
		return date
	}

	title := Slug(n.Title)
	if date == "" {
		return title
	}

	return date + "-" + title
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

// IsDailyNote reports whether the note is a daily note.
func (n Note) IsDailyNote() bool {
	return n.Kind == KindDaily
}

// IsChildNote reports whether the note is a child note.
func (n Note) IsChildNote() bool {
	return n.Kind == KindChild
}

// IsStandalone reports whether the note is not tied to a particular day.
func (n Note) IsStandalone() bool {
	return n.Kind == KindStandalone
}

// String returns the note's raw markdown source, including that of any
// children.
func (n Note) String() string {
	var raw strings.Builder
	raw.WriteString(n.Raw)

	for _, child := range n.Children {
		raw.WriteString(child.Raw)
	}

	return raw.String()
}
