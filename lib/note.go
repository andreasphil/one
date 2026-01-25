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

	if slug.Len() == 0 {
		util.Warnf("encountered empty slug, note will not be usable: %v", n.Raw)
	}

	return slug.String()
}

func (n Note) Content() string {
	titleExp := regexp.MustCompile(`^#{1,2}\s+.+\n`)
	content := titleExp.ReplaceAllString(n.Raw, "")
	return strings.TrimSpace(content)
}

func (n Note) Excerpt() string {
	words := strings.Split(n.Content(), " ")
	return strings.Join(words[0:int(math.Min(float64(len(words)), 40))], " ")
}

func (n Note) IsDailyNote() bool {
	return !n.Date.IsZero() && dateExp.MatchString(n.Title)
}

func FindRecursive(notes []Note, slug string) (Note, bool) {
	for _, v := range notes {
		if v.Slug() == slug {
			return v, true
		} else if len(v.Children) > 0 {
			if note, found := FindRecursive(v.Children, slug); found {
				return note, true
			}
		}
	}

	return Note{}, false
}
