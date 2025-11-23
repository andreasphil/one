package lib

import (
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/andreasphil/one/util"
)

var dateExp = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
var normalizeExp = regexp.MustCompile(`[^\wäöüß]+`)

type Note struct {
	Title    string
	Icon     string
	Date     time.Time
	Tags     util.Set[string]
	Children []Note
	Raw      string
}

func NewNote(title string) Note {
	return Note{
		Title: title,
		Tags:  util.NewSet[string](),
	}
}

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

func (n Note) Content() string {
	titleExp := regexp.MustCompile(`^#{1,2}\s+.+\n`)
	content := titleExp.ReplaceAllString(n.Raw, "")
	return strings.TrimSpace(content)
}

func (n Note) IsEmpty() bool {
	return len(n.Content()) == 0
}

func (n Note) IsDailyNote() bool {
	return !n.Date.IsZero() && dateExp.MatchString(n.Title)
}

func (n Note) Excerpt(maxWords int) (string, bool) {
	words := strings.Split(n.Content(), " ")
	partial := strings.Join(words[0:int(math.Min(float64(len(words)), float64(maxWords)))], " ")

	whitespaceExp := regexp.MustCompile(`\s+`)
	partial = whitespaceExp.ReplaceAllString(partial, " ")

	return partial, len(words) >= maxWords
}

func (n Note) String() string {
	var raw strings.Builder
	raw.WriteString(n.Raw)

	for _, note := range n.Children {
		raw.WriteString(note.String())
	}

	return raw.String()
}
