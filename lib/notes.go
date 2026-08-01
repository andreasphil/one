package lib

import (
	"cmp"
	"slices"
	"strings"
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

func GetRecursive(notes []Note, slug string) (Note, bool) {
	var found Note
	ok := false

	Walk(notes, func(note Note) bool {
		if note.Slug() != slug {
			return true
		}

		found = note
		ok = true
		return false
	})

	return found, ok
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

func Stringify(notes []Note) string {
	if len(notes) == 0 {
		return ""
	}

	rawNotes := make([]string, len(notes))
	for i, note := range notes {
		rawNotes[i] = strings.TrimRight(note.String(), "\n")
	}

	return strings.Join(rawNotes, "\n\n") + "\n"
}
