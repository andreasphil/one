package note_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/andreasphil/one/lib/note"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func searchResultTitles(notes []note.Note) []string {
	titles := make([]string, len(notes))
	for i, n := range notes {
		titles[i] = n.Title
	}

	return titles
}

func parseNotes(t *testing.T, input string) []note.Note {
	t.Helper()

	notes, err := note.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	return notes
}

const searchFixture = `# Groceries

Buy milk and eggs.

# Reading list

Books to read: SICP, TAOCP.

# 01.01.2026

Standup at 9.

## Groceries run

Went to the store for milk.
`

func TestContaining(t *testing.T) {
	type testcase struct {
		name     string
		query    string
		expected []string
	}

	testcases := []testcase{
		{
			name:     "matches in the title",
			query:    "Reading",
			expected: []string{"Reading list"},
		},
		{
			name:     "matches in the content",
			query:    "Standup",
			expected: []string{"01.01.2026"},
		},
		{
			name:     "matches a substring of a word",
			query:    "roceri",
			expected: []string{"Groceries", "Groceries run"},
		},
		{
			name:     "ignores case of the query",
			query:    "READING",
			expected: []string{"Reading list"},
		},
		{
			name:     "ignores case of the note",
			query:    "sicp",
			expected: []string{"Reading list"},
		},
		{
			name:     "returns matches in input order",
			query:    "milk",
			expected: []string{"Groceries", "Groceries run"},
		},
		{
			name:     "returns child notes as results of their own",
			query:    "store",
			expected: []string{"Groceries run"},
		},
		{
			name:     "does not match a parent for content of its children",
			query:    "Went to the store",
			expected: []string{"Groceries run"},
		},
		{
			name:     "matches across words",
			query:    "milk and eggs",
			expected: []string{"Groceries"},
		},
		{
			name:     "does not match when the query is not contained exactly",
			query:    "milkeggs",
			expected: []string{},
		},
		{
			name:     "does not match approximately",
			query:    "grocerys",
			expected: []string{},
		},
		{
			name:     "returns no matches for a query nothing contains",
			query:    "nonexistent",
			expected: []string{},
		},
		{
			name:  "returns all notes for an empty query",
			query: "",
			expected: []string{
				"Groceries", "Reading list", "01.01.2026", "Groceries run",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			notes := parseNotes(t, searchFixture)

			got := searchResultTitles(note.Containing(notes, tc.query))

			if diff := cmp.Diff(tc.expected, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Containing(notes, %q) mismatch (-want +got):\n%s", tc.query, diff)
			}
		})
	}
}

func TestContainingVisitsChildrenDepthFirst(t *testing.T) {
	input := `# 01.01.2026

Note about x.

## Child 1

About x.

## Child 2

About x.

# 02.01.2026

Also about x.
`

	notes := parseNotes(t, input)

	got := searchResultTitles(note.Containing(notes, "x"))
	want := []string{
		"01.01.2026", "Child 1", "Child 2", "02.01.2026",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Containing(notes, \"x\") mismatch (-want +got):\n%s", diff)
	}
}

func TestContainingTreatsQueryLiterally(t *testing.T) {
	type testcase struct {
		name     string
		query    string
		expected []string
	}

	input := `# Regex

Matching .* is fun (sometimes).

# Plain

Matching everything is fun.
`

	testcases := []testcase{
		{
			name:     "does not treat the query as a pattern",
			query:    ".*",
			expected: []string{"Regex"},
		},
		{
			name:     "matches parentheses literally",
			query:    "(sometimes)",
			expected: []string{"Regex"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			notes := parseNotes(t, input)

			got := searchResultTitles(note.Containing(notes, tc.query))

			if diff := cmp.Diff(tc.expected, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Containing(notes, %q) mismatch (-want +got):\n%s", tc.query, diff)
			}
		})
	}
}

func TestContainingHandlesEmptyInput(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		if got := note.Containing([]note.Note{}, "any"); len(got) != 0 {
			t.Errorf("Containing([], \"any\") = %v, want no matches", searchResultTitles(got))
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		if got := note.Containing(nil, "any"); len(got) != 0 {
			t.Errorf("Containing(nil, \"any\") = %v, want no matches", searchResultTitles(got))
		}
	})
}

func TestContainingDoesNotModifyInput(t *testing.T) {
	notes := parseNotes(t, searchFixture)
	before := parseNotes(t, searchFixture)

	note.Containing(notes, "milk")

	if diff := cmp.Diff(searchResultTitles(before), searchResultTitles(notes)); diff != "" {
		t.Errorf("Containing() modified its input (-before +after):\n%s", diff)
	}

	if got, want := len(notes[2].Children), len(before[2].Children); got != want {
		t.Errorf("Containing() changed children = %d, want %d", got, want)
	}
}

// Filtering ----------------------------------------------

func resultTitles(results []note.Result) []string {
	titles := make([]string, len(results))
	for i, r := range results {
		titles[i] = r.Note.Title
	}

	return titles
}

func matching(matches int) note.Filter {
	return func(note.Note) (bool, []note.Match) {
		if matches == 0 {
			return true, nil
		}

		return true, make([]note.Match, matches)
	}
}

func rejecting() note.Filter {
	return func(note.Note) (bool, []note.Match) {
		return false, nil
	}
}

func withTitle(titles ...string) note.Filter {
	return func(n note.Note) (bool, []note.Match) {
		return slices.Contains(titles, n.Title), nil
	}
}

func TestFilterChainApply(t *testing.T) {
	type testcase struct {
		name        string
		chain       note.FilterChain
		expected    bool
		expectedLen int
	}

	testcases := []testcase{
		{
			name:     "does not match for an empty chain",
			chain:    note.FilterChain{},
			expected: false,
		},
		{
			name:     "does not match for a nil chain",
			chain:    nil,
			expected: false,
		},
		{
			name:     "matches if the only filter matches",
			chain:    note.FilterChain{matching(0)},
			expected: true,
		},
		{
			name:     "does not match if the only filter rejects",
			chain:    note.FilterChain{rejecting()},
			expected: false,
		},
		{
			name:     "matches if every filter matches",
			chain:    note.FilterChain{matching(0), matching(0)},
			expected: true,
		},
		{
			name:     "does not match if a later filter rejects",
			chain:    note.FilterChain{matching(0), rejecting()},
			expected: false,
		},
		{
			name:     "does not match if an earlier filter rejects",
			chain:    note.FilterChain{rejecting(), matching(0)},
			expected: false,
		},
		{
			name:        "collects the matches of a single filter",
			chain:       note.FilterChain{matching(2)},
			expected:    true,
			expectedLen: 2,
		},
		{
			name:        "collects the matches of every filter",
			chain:       note.FilterChain{matching(1), matching(2)},
			expected:    true,
			expectedLen: 3,
		},
		{
			name:        "skips filters that match without reporting matches",
			chain:       note.FilterChain{matching(1), matching(0), matching(2)},
			expected:    true,
			expectedLen: 3,
		},
		{
			name:     "matches without matches if no filter reports any",
			chain:    note.FilterChain{matching(0), matching(0)},
			expected: true,
		},
		{
			name:     "discards the matches collected before a rejection",
			chain:    note.FilterChain{matching(2), rejecting()},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, gotMatches := tc.chain.Apply(note.New("Any"))

			if got != tc.expected {
				t.Errorf("Apply(note) = %v, want %v", got, tc.expected)
			}

			if len(gotMatches) != tc.expectedLen {
				t.Errorf("Apply(note) matches = %d, want %d", len(gotMatches), tc.expectedLen)
			}
		})
	}
}

func TestFilterChainApplyStopsAtTheFirstRejection(t *testing.T) {
	calls := 0
	counting := func(note.Note) (bool, []note.Match) {
		calls++
		return true, nil
	}

	chain := note.FilterChain{rejecting(), counting}

	if ok, _ := chain.Apply(note.New("Any")); ok {
		t.Error("Apply(note) = true, want false")
	}

	if calls != 0 {
		t.Errorf("Apply(note) called %d filters after a rejection, want 0", calls)
	}
}

func TestSearch(t *testing.T) {
	type testcase struct {
		name     string
		chain    note.FilterChain
		expected []string
	}

	testcases := []testcase{
		{
			name:  "returns every note if the chain matches everything",
			chain: note.FilterChain{matching(0)},
			expected: []string{
				"Groceries", "Reading list", "01.01.2026", "Groceries run",
			},
		},
		{
			name:     "returns no results if the chain rejects everything",
			chain:    note.FilterChain{rejecting()},
			expected: []string{},
		},
		{
			name:     "returns no results for an empty chain",
			chain:    note.FilterChain{},
			expected: []string{},
		},
		{
			name:     "returns only the notes the chain matches",
			chain:    note.FilterChain{withTitle("Reading list")},
			expected: []string{"Reading list"},
		},
		{
			name:     "returns child notes of non-matching parents as results of their own",
			chain:    note.FilterChain{withTitle("Groceries run")},
			expected: []string{"Groceries run"},
		},
		{
			name:     "applies every filter in the chain",
			chain:    note.FilterChain{withTitle("Groceries", "Reading list"), withTitle("Groceries")},
			expected: []string{"Groceries"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			notes := parseNotes(t, searchFixture)

			got := resultTitles(note.Search(notes, tc.chain))

			if diff := cmp.Diff(tc.expected, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Search(notes, chain) mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSearchVisitsChildrenDepthFirst(t *testing.T) {
	input := `# 01.01.2026

Note about x.

## Child 1

About x.

## Child 2

About x.

# 02.01.2026

Also about x.
`

	notes := parseNotes(t, input)

	got := resultTitles(note.Search(notes, note.FilterChain{matching(0)}))
	want := []string{
		"01.01.2026", "Child 1", "Child 2", "02.01.2026",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Search(notes, chain) mismatch (-want +got):\n%s", diff)
	}
}

func TestSearchReturnsMatchesWithEachResult(t *testing.T) {
	notes := parseNotes(t, searchFixture)
	chain := note.FilterChain{matching(1), matching(2)}

	got := note.Search(notes, chain)
	if len(got) == 0 {
		t.Fatal("Search(notes, chain) returned no results, want all notes")
	}

	for _, r := range got {
		if len(r.Matches) != 3 {
			t.Errorf("Search(notes, chain) matches for %q = %d, want 3", r.Note.Title, len(r.Matches))
		}
	}
}

func TestSearchHandlesEmptyInput(t *testing.T) {
	chain := note.FilterChain{matching(0)}

	t.Run("empty slice", func(t *testing.T) {
		if got := note.Search([]note.Note{}, chain); len(got) != 0 {
			t.Errorf("Search([], chain) = %v, want no results", resultTitles(got))
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		if got := note.Search(nil, chain); len(got) != 0 {
			t.Errorf("Search(nil, chain) = %v, want no results", resultTitles(got))
		}
	})
}

func TestSearchDoesNotModifyInput(t *testing.T) {
	notes := parseNotes(t, searchFixture)
	before := parseNotes(t, searchFixture)

	note.Search(notes, note.FilterChain{matching(0)})

	if diff := cmp.Diff(searchResultTitles(before), searchResultTitles(notes)); diff != "" {
		t.Errorf("Search() modified its input (-before +after):\n%s", diff)
	}

	if got, want := len(notes[2].Children), len(before[2].Children); got != want {
		t.Errorf("Search() changed children = %d, want %d", got, want)
	}
}
