package lib

// DuplicateSlugs returns the slugs of any notes (including children) that share
// their slug with another note. Slugs are returned in the order they first
// occur, and each duplicate slug is only returned once, regardless of how many
// notes share it.
func DuplicateSlugs(notes []Note) []string {
	var order []string
	counts := make(map[string]int)

	Walk(notes, func(note Note) bool {
		slug := note.Slug()
		if _, seen := counts[slug]; !seen {
			order = append(order, slug)
		}
		counts[slug]++
		return true
	})

	var duplicates []string
	for _, slug := range order {
		if counts[slug] > 1 {
			duplicates = append(duplicates, slug)
		}
	}

	return duplicates
}

// EmptyTitles returns the number of notes (including children) that have an
// empty title.
func EmptyTitles(notes []Note) int {
	count := 0

	Walk(notes, func(note Note) bool {
		if note.Title == "" {
			count++
		}
		return true
	})

	return count
}

// EmptyNotes returns the titles of any notes (including children) that have no
// content. Notes with children are never considered empty, even if their own
// content is empty, since they still hold information through their children.
func EmptyNotes(notes []Note) []string {
	var empty []string

	Walk(notes, func(note Note) bool {
		if len(note.Children) == 0 && note.IsEmpty() {
			empty = append(empty, note.Title)
		}
		return true
	})

	return empty
}
