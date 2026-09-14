package note

import (
	"cmp"
	"slices"
	"strings"

	"github.com/andreasphil/one/util"
)

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

// Walk visits every note and its children, in depth-first, pre-order (i.e. a
// note is visited before its children). fn is called once per note. Traversal
// stops as soon as fn returns false, in which case Walk also returns false. If
// every note was visited, Walk returns true.
func Walk(notes []Note, fn func(Note) bool) bool {
	for _, note := range notes {
		if !fn(note) {
			return false
		}

		for _, child := range note.Children {
			if !fn(child) {
				return false
			}
		}
	}

	return true
}

func Flat(notes []Note) []Note {
	flat := make([]Note, 0, Count(notes))

	Walk(notes, func(note Note) bool {
		flat = append(flat, note)
		return true
	})

	return flat
}

func FindBySlug(notes []Note, slug string) (Note, bool) {
	return find(notes, func(note Note) bool {
		return note.Slug() == slug
	})
}

// ResolveSlug returns the slug of the note matching target along with true,
// ignoring case and any characters that slugs do not preserve. A note matches
// if either its title or its own slug matches, so that a note can also be
// linked by the slug shown in its URL. If several notes match, the first one in
// depth-first, pre-order wins. If no note matches, it returns the slugified
// target and false, so that links to notes that do not exist point at where the
// note would live.
func ResolveSlug(notes []Note, target string) (string, bool) {
	slug := Slug(target)

	n, found := find(notes, func(note Note) bool {
		return Slug(note.Title) == slug || note.Slug() == slug
	})

	if !found {
		return slug, false
	}

	return n.Slug(), true
}

func Tags(notes []Note) []Tag {
	tags := util.NewSet[Tag]()

	Walk(notes, func(note Note) bool {
		tags.Add(note.Tags.Values()...)
		return true
	})

	values := tags.Values()
	slices.SortFunc(values, func(a Tag, b Tag) int {
		return cmp.Or(
			cmp.Compare(strings.ToLower(a.Name()), strings.ToLower(b.Name())),
			cmp.Compare(a.Name(), b.Name()),
		)
	})

	return values
}

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

func Count(notes []Note) int {
	count := 0

	for _, note := range notes {
		count += 1 + len(note.Children)
	}

	return count
}
