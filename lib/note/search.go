package note

import (
	"strings"
)

// Filtering ----------------------------------------------

type Match struct{}

type Filter func(Note) (bool, []Match)

type FilterChain []Filter

func (c FilterChain) Apply(n Note) (bool, []Match) {
	if len(c) == 0 {
		return false, nil
	}

	var result []Match

	for _, f := range c {
		ok, m := f(n)
		if !ok {
			return false, nil
		}

		result = append(result, m...)
	}

	return true, result
}

type Result struct {
	Note    Note
	Matches []Match
}

func Search(notes []Note, chain FilterChain) []Result {
	found := []Result{}

	Walk(notes, func(n Note) bool {
		if ok, m := chain.Apply(n); ok {
			found = append(found, Result{Note: n, Matches: m})
		}

		return true
	})

	return found
}

// Simple string search -----------------------------------

func containsExactPhrase(query string) Filter {
	normalizedQuery := strings.ToLower(query)

	return func(n Note) (bool, []Match) {
		normalizedRaw := strings.ToLower(n.Raw)
		if strings.Contains(normalizedRaw, normalizedQuery) {
			return true, nil
		}

		return false, nil
	}
}

func Containing(notes []Note, query string) []Note {
	results := []Note{}
	chain := FilterChain{containsExactPhrase(query)}

	for _, i := range Search(notes, chain) {
		results = append(results, i.Note)
	}

	return results
}
