package lib

import (
	"cmp"
	"slices"
	"strings"
)

func GetRecursive(notes []Note, slug string) (Note, bool) {
	for _, v := range notes {
		if v.Slug() == slug {
			return v, true
		} else if len(v.Children) > 0 {
			if note, found := GetRecursive(v.Children, slug); found {
				return note, true
			}
		}
	}

	return Note{}, false
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
