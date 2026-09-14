package note

import (
	"regexp"
	"strings"
	"time"

	"github.com/andreasphil/one/util"
)

var normalizeExp = regexp.MustCompile(`[^\wäöüß]+`)
var titleHeadingExp = regexp.MustCompile(`^#{1,2}\s+.+\n`)

// Expressions for cleaning up excerpts -------------------

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

// Note kind ----------------------------------------------

type Kind int

const (
	KindStandalone Kind = iota
	KindDaily
	KindChild
)

// Tags ---------------------------------------------------

type Tag struct{ name string }

// NewTag returns the Tag with the given name. The name may be given with or
// without a leading "#".
func NewTag(name string) Tag {
	return Tag{name: strings.TrimPrefix(name, "#")}
}

// Name returns the tag without its leading "#".
func (t Tag) Name() string {
	return t.name
}

// String returns the tag including its leading "#".
func (t Tag) String() string {
	return "#" + t.name
}

func (t Tag) Equal(other Tag) bool {
	return t == other
}

// Note ---------------------------------------------------

type Note struct {
	Title    string
	Kind     Kind
	Icon     string
	Date     time.Time
	Tags     util.Set[Tag]
	Children []Note
	Raw      string
}

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

// Slug returns a URL-friendly identifier for the note. Slugs are not
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

func (n Note) IsEmpty() bool {
	return len(n.Content()) == 0
}

func (n Note) IsDailyNote() bool {
	return n.Kind == KindDaily
}

func (n Note) IsChildNote() bool {
	return n.Kind == KindChild
}

func (n Note) IsStandalone() bool {
	return n.Kind == KindStandalone
}

// String returns the note's raw markdown source (including children).
func (n Note) String() string {
	var raw strings.Builder
	raw.WriteString(n.Raw)

	for _, child := range n.Children {
		raw.WriteString(child.Raw)
	}

	return raw.String()
}
