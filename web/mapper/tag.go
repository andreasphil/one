package mapper

import "github.com/andreasphil/one/lib/note"

// ToTags returns the names of the unique tags occurring in notes and their
// children, without the leading "#", sorted alphabetically (ascending, case
// insensitive). The result is never nil, so that it serializes to an empty
// JSON array rather than to null.
func ToTags(n []note.Note) []string {
	tags := note.Tags(n)
	m := make([]string, 0, len(tags))

	for _, tag := range tags {
		m = append(m, tag.Name())
	}

	return m
}
