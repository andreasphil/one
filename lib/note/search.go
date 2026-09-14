package note

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
