package note

import "strings"

func Containing(notes []Note, query string) []Note {
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
