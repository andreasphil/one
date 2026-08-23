package note

import "strings"

// SearchContainingString returns the notes that contain query case-insensitive,
// but otherwise exactly. Notes are matched on their raw markdown source. Child
// notes are searched and returned as results of their own. Results are not
// ranked. The order of the returned slice corresponds to the order of the
// input, depth-first for child notes.
func SearchContainingString(notes []Note, query string) []Note {
	results := []Note{}
	normalizedQuery := strings.ToLower(query)

	Walk(notes, func(note Note) bool {
		normalizedRaw := strings.ToLower(note.Raw)
		if strings.Contains(normalizedRaw, normalizedQuery) {
			results = append(results, note)
		}

		return true
	})

	return results
}
