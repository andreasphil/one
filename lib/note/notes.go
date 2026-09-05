package note

import (
	"cmp"
	"slices"
	"strings"

	"github.com/andreasphil/one/util"
)

// Walk recursively visits every note and its children, in depth-first,
// pre-order (i.e. a note is visited before its children). fn is called once per
// note. Traversal stops as soon as fn returns false, in which case Walk also
// returns false. If every note was visited, Walk returns true.
func Walk(notes []Note, fn func(Note) bool) bool {
	for _, note := range notes {
		if !fn(note) {
			return false
		}

		if len(note.Children) > 0 && !Walk(note.Children, fn) {
			return false
		}
	}

	return true
}

func find(notes []Note, match func(Note) bool) (Note, bool) {
	var found Note
	ok := false

	Walk(notes, func(note Note) bool {
		if !match(note) {
			return true
		}

		found = note
		ok = true
		return false
	})

	return found, ok
}

// FindBySlug searches notes and their children for a note matching slug,
// returning it along with true if found. If no note matches, it returns a
// zero-value Note and false.
func FindBySlug(notes []Note, slug string) (Note, bool) {
	return find(notes, func(note Note) bool {
		return note.Slug() == slug
	})
}

// ResolveSlug returns the slug of the note whose title matches target along
// with true, ignoring case and any characters that slugs do not preserve. If
// several notes match, the first one in depth-first, pre-order wins. If no
// note matches, it returns the slugified target and false, so that links to
// notes that do not exist point at where the note would live.
func ResolveSlug(notes []Note, target string) (string, bool) {
	slug := Slug(target)

	n, found := find(notes, func(note Note) bool {
		return Slug(note.Title) == slug
	})

	if !found {
		return slug, false
	}

	return n.Slug(), true
}

// Sort sorts notes in place, with daily notes ordered by date (descending)
// before all other notes ordered alphabetically by title (ascending, case
// insensitive). It returns the sorted notes along with whether sorting actually
// changed the order (i.e. whether notes were not already sorted).
func Sort(notes []Note) ([]Note, bool) {
	compare := func(a Note, b Note) int {
		return cmp.Or(
			b.Date.Compare(a.Date),
			cmp.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title)),
		)
	}

	isSorted := slices.IsSortedFunc(notes, compare)
	if !isSorted {
		slices.SortStableFunc(notes, compare)
	}

	return notes, !isSorted
}

// String serializes notes back into their notes file markdown representation,
// with a blank line between top-level notes and a single trailing newline.
func String(notes []Note) string {
	if len(notes) == 0 {
		return ""
	}

	rawNotes := make([]string, len(notes))
	for i, note := range notes {
		rawNotes[i] = strings.TrimRight(note.String(), "\n")
	}

	return strings.Join(rawNotes, "\n\n") + "\n"
}

// Tags returns the unique tags occurring in notes and their children, sorted
// alphabetically (ascending, case insensitive).
func Tags(notes []Note) []Tag {
	tags := util.NewSet[Tag]()

	Walk(notes, func(note Note) bool {
		tags.Add(note.Tags.Values()...)
		return true
	})

	values := tags.Values()
	slices.SortFunc(values, func(a Tag, b Tag) int {
		return cmp.Or(
			cmp.Compare(strings.ToLower(string(a)), strings.ToLower(string(b))),
			cmp.Compare(a, b),
		)
	})

	return values
}
