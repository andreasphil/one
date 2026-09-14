package note

import "strings"

func FilterNot(filter Filter) Filter {
	return func(n Note) (bool, []Match) {
		ok, _ := filter(n)
		return !ok, nil
	}
}

func FilterHasTag(tag Tag) Filter {
	return func(n Note) (bool, []Match) {
		return n.Tags.Has(tag), nil
	}
}

func FilterIsTagged() Filter {
	return func(n Note) (bool, []Match) {
		return n.Tags.Len() > 0, nil
	}
}

func FilterExactPhrase(phrase string, caseSensitive bool) Filter {
	if !caseSensitive {
		phrase = strings.ToLower(phrase)
	}

	return func(n Note) (bool, []Match) {
		noteRaw := n.Raw
		if !caseSensitive {
			noteRaw = strings.ToLower(noteRaw)
		}

		if strings.Contains(noteRaw, phrase) {
			return true, nil
		}

		return false, nil
	}
}
