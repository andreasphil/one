package note_test

import (
	"testing"

	"github.com/andreasphil/one/lib/note"
	"github.com/andreasphil/one/util"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func applies(f note.Filter, n note.Note) bool {
	ok, _ := f(n)
	return ok
}

func tagged(title string, tags ...string) note.Note {
	values := make([]note.Tag, 0, len(tags))
	for _, tag := range tags {
		values = append(values, note.NewTag(tag))
	}

	return note.Note{Title: title, Tags: util.NewSetFrom(values)}
}

func TestFilterNot(t *testing.T) {
	type testcase struct {
		name     string
		filter   note.Filter
		expected bool
	}

	testcases := []testcase{
		{
			name:     "rejects what the filter matches",
			filter:   matching(0),
			expected: false,
		},
		{
			name:     "matches what the filter rejects",
			filter:   rejecting(),
			expected: true,
		},
		{
			name:     "cancels out when applied twice",
			filter:   note.FilterNot(rejecting()),
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := applies(note.FilterNot(tc.filter), note.New("Any"))

			if got != tc.expected {
				t.Errorf("FilterNot(filter)(note) = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestFilterNotDiscardsMatches(t *testing.T) {
	_, got := note.FilterNot(rejecting())(note.New("Any"))

	if len(got) != 0 {
		t.Errorf("FilterNot(filter)(note) matches = %d, want 0", len(got))
	}
}

func TestFilterHasTag(t *testing.T) {
	type testcase struct {
		name     string
		note     note.Note
		tag      note.Tag
		expected bool
	}

	testcases := []testcase{
		{
			name:     "matches a note carrying the tag",
			note:     tagged("A", "work"),
			tag:      note.NewTag("work"),
			expected: true,
		},
		{
			name:     "matches one of several tags",
			note:     tagged("A", "work", "idea", "recipe"),
			tag:      note.NewTag("idea"),
			expected: true,
		},
		{
			name:     "does not match a note carrying other tags",
			note:     tagged("A", "work"),
			tag:      note.NewTag("idea"),
			expected: false,
		},
		{
			name:     "does not match an untagged note",
			note:     note.New("A"),
			tag:      note.NewTag("work"),
			expected: false,
		},
		{
			name:     "does not match a tag the note's tag starts with",
			note:     tagged("A", "workshop"),
			tag:      note.NewTag("work"),
			expected: false,
		},
		{
			name:     "does not match a tag that starts with the note's tag",
			note:     tagged("A", "work"),
			tag:      note.NewTag("workshop"),
			expected: false,
		},
		{
			name:     "accepts the tag name with a leading hash",
			note:     tagged("A", "work"),
			tag:      note.NewTag("#work"),
			expected: true,
		},
		{
			name:     "distinguishes tags differing in case",
			note:     tagged("A", "Work"),
			tag:      note.NewTag("work"),
			expected: false,
		},
		{
			name:     "does not match a note with a nil tag set",
			note:     note.Note{Title: "A"},
			tag:      note.NewTag("work"),
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := applies(note.FilterHasTag(tc.tag), tc.note)

			if got != tc.expected {
				t.Errorf("FilterHasTag(%v)(note) = %v, want %v", tc.tag, got, tc.expected)
			}
		})
	}
}

func TestFilterIsTagged(t *testing.T) {
	type testcase struct {
		name     string
		note     note.Note
		expected bool
	}

	testcases := []testcase{
		{
			name:     "matches a note with one tag",
			note:     tagged("A", "work"),
			expected: true,
		},
		{
			name:     "matches a note with several tags",
			note:     tagged("A", "work", "idea"),
			expected: true,
		},
		{
			name:     "does not match a note with no tags",
			note:     note.New("A"),
			expected: false,
		},
		{
			name:     "does not match a note with a nil tag set",
			note:     note.Note{Title: "A"},
			expected: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := applies(note.FilterIsTagged(), tc.note)

			if got != tc.expected {
				t.Errorf("FilterIsTagged()(note) = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestFilterExactPhrase(t *testing.T) {
	type testcase struct {
		name     string
		phrase   string
		expected []string
	}

	testcases := []testcase{
		{
			name:     "matches in the title",
			phrase:   "Reading",
			expected: []string{"Reading list"},
		},
		{
			name:     "matches in the content",
			phrase:   "Standup",
			expected: []string{"01.01.2026"},
		},
		{
			name:     "matches a substring of a word",
			phrase:   "roceri",
			expected: []string{"Groceries", "Groceries run"},
		},
		{
			name:     "ignores case of the phrase",
			phrase:   "READING",
			expected: []string{"Reading list"},
		},
		{
			name:     "ignores case of the note",
			phrase:   "sicp",
			expected: []string{"Reading list"},
		},
		{
			name:     "matches across words",
			phrase:   "milk and eggs",
			expected: []string{"Groceries"},
		},
		{
			name:     "returns child notes as results of their own",
			phrase:   "store",
			expected: []string{"Groceries run"},
		},
		{
			name:     "does not match a parent for content of its children",
			phrase:   "Went to the store",
			expected: []string{"Groceries run"},
		},
		{
			name:     "does not match when the phrase is not contained exactly",
			phrase:   "milkeggs",
			expected: []string{},
		},
		{
			name:     "does not match approximately",
			phrase:   "grocerys",
			expected: []string{},
		},
		{
			name:     "does not match a phrase nothing contains",
			phrase:   "nonexistent",
			expected: []string{},
		},
		{
			name:   "matches every note for an empty phrase",
			phrase: "",
			expected: []string{
				"Groceries", "Reading list", "01.01.2026", "Groceries run",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			notes := parseNotes(t, searchFixture)
			chain := note.FilterChain{note.FilterExactPhrase(tc.phrase, false)}

			got := resultTitles(note.Search(notes, chain))

			if diff := cmp.Diff(tc.expected, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("FilterExactPhrase(%q, false) mismatch (-want +got):\n%s", tc.phrase, diff)
			}
		})
	}
}

func TestFilterExactPhraseCaseSensitive(t *testing.T) {
	type testcase struct {
		name     string
		phrase   string
		expected []string
	}

	testcases := []testcase{
		{
			name:     "matches a phrase in the same case",
			phrase:   "SICP",
			expected: []string{"Reading list"},
		},
		{
			name:     "does not match a phrase in a different case",
			phrase:   "sicp",
			expected: []string{},
		},
		{
			name:     "does not match a title in a different case",
			phrase:   "reading",
			expected: []string{},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			notes := parseNotes(t, searchFixture)
			chain := note.FilterChain{note.FilterExactPhrase(tc.phrase, true)}

			got := resultTitles(note.Search(notes, chain))

			if diff := cmp.Diff(tc.expected, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("FilterExactPhrase(%q, true) mismatch (-want +got):\n%s", tc.phrase, diff)
			}
		})
	}
}

func TestFilterExactPhraseTreatsPhraseLiterally(t *testing.T) {
	type testcase struct {
		name     string
		phrase   string
		expected []string
	}

	input := `# Regex

Matching .* is fun (sometimes).

# Plain

Matching everything is fun.
`

	testcases := []testcase{
		{
			name:     "does not treat the phrase as a pattern",
			phrase:   ".*",
			expected: []string{"Regex"},
		},
		{
			name:     "matches parentheses literally",
			phrase:   "(sometimes)",
			expected: []string{"Regex"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			notes := parseNotes(t, input)
			chain := note.FilterChain{note.FilterExactPhrase(tc.phrase, false)}

			got := resultTitles(note.Search(notes, chain))

			if diff := cmp.Diff(tc.expected, got, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("FilterExactPhrase(%q, false) mismatch (-want +got):\n%s", tc.phrase, diff)
			}
		})
	}
}
