package note

import (
	"regexp"
	"strings"
	"time"

	"github.com/andreasphil/one/util"
)

var normalizeExp = regexp.MustCompile(`[^\wäöüß]+`)

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
	Title      string
	Kind       Kind
	Date       time.Time
	HasQuote   bool
	HasSnippet bool
	Icon       string
	Tags       util.Set[Tag]
	Children   []Note
	Raw        string
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
	_, content, _ := strings.Cut(n.Raw, "\n")
	return strings.TrimSpace(content)
}

// IsEmpty reports whether the note has no content. Notes with children are
// never considered empty, even if their own content is empty.
func (n Note) IsEmpty() bool {
	return len(n.Children) == 0 && len(n.Content()) == 0
}

func (n Note) IsTagged() bool {
	return n.Tags.Len() > 0
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
