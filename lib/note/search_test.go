package note_test

import (
	"strings"
	"testing"

	"github.com/andreasphil/one/lib/note"
	"github.com/google/go-cmp/cmp"
)

func searchResultTitles(notes []note.Note) []string {
	titles := make([]string, len(notes))
	for i, n := range notes {
		titles[i] = n.Title
	}

	return titles
}

func parseForSearch(t *testing.T, input string) []note.Note {
	t.Helper()

	notes, err := note.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return notes
}

// Contains an undated note, a second undated note, and a daily
// note with one child note. "milk" and "groceries" each occur in more than one
// note, "Standup" only in the daily note itself, and "store" only in its child.
const searchFixture = `# Groceries

Buy milk and eggs.

# Reading list

Books to read: SICP, TAOCP.

# 01.01.2026

Standup at 9.

## Groceries run

Went to the store for milk.
`

func TestSearchContainingString(t *testing.T) {
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
			notes := parseForSearch(t, searchFixture)

			result := searchResultTitles(note.SearchContainingString(notes, tc.query))

			if len(result) != len(tc.expected) {
				t.Fatalf("expected %v matches %v, got %v matches %v",
					len(tc.expected), tc.expected, len(result), result)
			}

			if !cmp.Equal(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestSearchContainingStringVisitsChildrenDepthFirst(t *testing.T) {
	input := `# 01.01.2026

Note about x.

## Child 1

About x.

## Child 2

About x.

# 02.01.2026

Also about x.
`

	notes := parseForSearch(t, input)

	result := searchResultTitles(note.SearchContainingString(notes, "x"))
	expected := []string{
		"01.01.2026", "Child 1", "Child 2", "02.01.2026",
	}

	if !cmp.Equal(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestSearchContainingStringTreatsQueryLiterally(t *testing.T) {
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
			notes := parseForSearch(t, input)

			result := searchResultTitles(note.SearchContainingString(notes, tc.query))

			if !cmp.Equal(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestSearchContainingStringHandlesEmptyInput(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		result := note.SearchContainingString([]note.Note{}, "any")
		if len(result) != 0 {
			t.Errorf("expected no matches, got %v", searchResultTitles(result))
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		result := note.SearchContainingString(nil, "any")
		if len(result) != 0 {
			t.Errorf("expected no matches, got %v", searchResultTitles(result))
		}
	})
}

func TestSearchContainingStringDoesNotModifyInput(t *testing.T) {
	notes := parseForSearch(t, searchFixture)
	before := parseForSearch(t, searchFixture)

	note.SearchContainingString(notes, "milk")

	if !cmp.Equal(searchResultTitles(notes), searchResultTitles(before)) {
		t.Errorf("expected input to be unchanged, got %v", searchResultTitles(notes))
	}

	if len(notes[2].Children) != len(before[2].Children) {
		t.Errorf("expected children to be unchanged, got %v", notes[2].Children)
	}
}
