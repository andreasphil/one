// Package mapper converts notes into the view models used for rendering the
// web templates and for sending state to the client.
package mapper

import "github.com/andreasphil/one/lib/note"

// NoteMeta describes a note with the minimum of information needed for
// listing and linking to it.
type NoteMeta struct {
	Title string
	Slug  string
}

// NewNoteMeta returns the NoteMeta describing n.
func NewNoteMeta(n note.Note) NoteMeta {
	return NoteMeta{Title: n.Title, Slug: n.Slug()}
}

// ToNoteMeta flattens notes and their children into one NoteMeta each, in
// depth-first, pre-order. The result is never nil, so that it serializes to
// an empty JSON array rather than to null.
func ToNoteMeta(n []note.Note) []NoteMeta {
	m := make([]NoteMeta, 0, len(n))

	note.Walk(n, func(i note.Note) bool {
		m = append(m, NewNoteMeta(i))
		return true
	})

	return m
}
