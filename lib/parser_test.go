package lib_test

import (
	"strings"
	"testing"
	"time"

	"github.com/andreasphil/one/lib"
	"github.com/google/go-cmp/cmp"
)

func TestChunkNotes(t *testing.T) {
	type testcase struct {
		name     string
		input    string
		expected []lib.Note
	}

	testcases := []testcase{
		{
			name:  "simple note",
			input: "# Note 1\n\nLine 1\n\nLine 2\n",
			expected: []lib.Note{
				{
					Title: "Note 1",
					Raw:   "# Note 1\n\nLine 1\n\nLine 2\n",
				},
			},
		},
		{
			name:  "two notes",
			input: "# Note 1\n\nLine 1\n\nLine 2\n\n# Note 2\n\nLine 3\n\nLine 4\n",
			expected: []lib.Note{
				{
					Title: "Note 1",
					Raw:   "# Note 1\n\nLine 1\n\nLine 2\n\n",
				},
				{
					Title: "Note 2",
					Raw:   "# Note 2\n\nLine 3\n\nLine 4\n",
				},
			},
		},
		{
			name:  "simple note with sub headings",
			input: "# Note 1\n\nLine 1\n\nLine 2\n\n## Child Note 1\n\nLine 3\n\nLine 4\n\n# Note 2\n\nLine 5\n\nLine 6\n",
			expected: []lib.Note{
				{
					Title: "Note 1",
					Raw:   "# Note 1\n\nLine 1\n\nLine 2\n\n## Child Note 1\n\nLine 3\n\nLine 4\n\n",
				},
				{
					Title: "Note 2",
					Raw:   "# Note 2\n\nLine 5\n\nLine 6\n",
				},
			},
		},
		{
			name:  "daily note",
			input: "# 01.01.2026\n\nLine 1\n\nLine 2\n",
			expected: []lib.Note{
				{
					Title: "01.01.2026",
					Raw:   "# 01.01.2026\n\nLine 1\n\nLine 2\n",
					Date:  time.Date(2026, time.January, 01, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "note with one child note",
			input: "# 01.01.2026\n\nLine 1\n\nLine 2\n\n## Child Note 1\n\nLine 3\n\nLine 4\n\n# 02.01.2026\n\nLine 5\n\nLine 6\n",
			expected: []lib.Note{
				{
					Title: "01.01.2026",
					Raw:   "# 01.01.2026\n\nLine 1\n\nLine 2\n\n",
					Date:  time.Date(2026, time.January, 01, 0, 0, 0, 0, time.UTC),
					Children: []lib.Note{
						{
							Title: "Child Note 1",
							Raw:   "## Child Note 1\n\nLine 3\n\nLine 4\n\n",
							Date:  time.Date(2026, time.January, 01, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				{
					Title: "02.01.2026",
					Raw:   "# 02.01.2026\n\nLine 5\n\nLine 6\n",
					Date:  time.Date(2026, time.January, 02, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "note with multiple child notes",
			input: "# 01.01.2026\n\nLine 1\n\nLine 2\n\n## Child Note 1\n\nLine 3\n\nLine 4\n\n## Child Note 2\n\nLine 5\n\nLine 6\n\n# 02.01.2026\n\nLine 7\n\nLine 8\n",
			expected: []lib.Note{
				{
					Title: "01.01.2026",
					Raw:   "# 01.01.2026\n\nLine 1\n\nLine 2\n\n",
					Date:  time.Date(2026, time.January, 01, 0, 0, 0, 0, time.UTC),
					Children: []lib.Note{
						{
							Title: "Child Note 1",
							Raw:   "## Child Note 1\n\nLine 3\n\nLine 4\n\n",
							Date:  time.Date(2026, time.January, 01, 0, 0, 0, 0, time.UTC),
						},
						{
							Title: "Child Note 2",
							Raw:   "## Child Note 2\n\nLine 5\n\nLine 6\n\n",
							Date:  time.Date(2026, time.January, 01, 0, 0, 0, 0, time.UTC),
						},
					},
				},
				{
					Title: "02.01.2026",
					Raw:   "# 02.01.2026\n\nLine 7\n\nLine 8\n",
					Date:  time.Date(2026, time.January, 02, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "note with fenced code block",
			input: "# 01.01.2026\n\nLine 1\n\n```\n\n# Block comment\n\n```\n",
			expected: []lib.Note{
				{
					Title: "01.01.2026",
					Raw:   "# 01.01.2026\n\nLine 1\n\n```\n\n# Block comment\n\n```\n",
					Date:  time.Date(2026, time.January, 01, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "level 1 heading in fenced code block",
			input: "# Note 1\n\n```\n# Not a new note\n```\n",
			expected: []lib.Note{
				{
					Title: "Note 1",
					Raw:   "# Note 1\n\n```\n# Not a new note\n```\n",
				},
			},
		},
		{
			name:  "level 2 heading in fenced code block",
			input: "# 01.01.2026\n\n```\n## Not a child note\n```\n",
			expected: []lib.Note{
				{
					Title: "01.01.2026",
					Raw:   "# 01.01.2026\n\n```\n## Not a child note\n```\n",
					Date:  time.Date(2026, time.January, 01, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:  "no trailing newline",
			input: "# Note 1",
			expected: []lib.Note{
				{
					Title: "Note 1",
					Raw:   "# Note 1\n",
				},
			},
		},
	}

	ignoreTags := cmp.FilterPath(func(p cmp.Path) bool {
		return strings.Contains(p.String(), "Tags") || strings.Contains(p.String(), "Icon")
	}, cmp.Ignore())

	for _, i := range testcases {
		t.Run(i.name, func(t *testing.T) {
			result, err := lib.Parse(strings.NewReader(i.input))

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(i.expected) {
				t.Errorf("got %d notes, expected %v", len(result), len(i.expected))
			}

			if diff := cmp.Diff(i.expected, result, ignoreTags); diff != "" {
				t.Errorf("note mismatch:\n%s", diff)
			}
		})
	}
}

func TestRequireHeading(t *testing.T) {
	input := "test\n\n# Note 1"
	_, err := lib.Parse(strings.NewReader(input))

	if err == nil {
		t.Errorf("expected to return error when heading is missing")
	}
}

func TestRequireContent(t *testing.T) {
	input := ""
	result, err := lib.Parse(strings.NewReader(input))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected to return empty list when input is empty")
	}
}

func TestExtractTags(t *testing.T) {
	type testcase struct {
		name     string
		input    string
		expected lib.Set[string]
	}

	testcases := []testcase{
		{
			name:     "no tags",
			input:    "# Note 1\n\nLine 1\n",
			expected: lib.NewSet[string](),
		},
		{
			name:     "tag in title",
			input:    "# Note 1 #tag\n",
			expected: lib.NewSetFrom([]string{"#tag"}),
		},
		{
			name:     "tag in body",
			input:    "# Note 1\n\nLine 1 #tag\n",
			expected: lib.NewSetFrom([]string{"#tag"}),
		},
		{
			name:     "multiple tags",
			input:    "# Note 1 #tag_1\n\nLine 1 #tag_2\n",
			expected: lib.NewSetFrom([]string{"#tag_1", "#tag_2"}),
		},
		{
			name:     "multiple tags in one line",
			input:    "# Note 1\n\nLine #tag_1 1 #tag_2\n",
			expected: lib.NewSetFrom([]string{"#tag_1", "#tag_2"}),
		},
		{
			name:     "duplicate tags",
			input:    "# Note 1\n\nLine 1 #tag_1\n\nLine 2 #tag_1",
			expected: lib.NewSetFrom([]string{"#tag_1"}),
		},
		{
			name:     "tags in fenced code block",
			input:    "# Note 1\n\n```\n#tag_in_code\n```\n",
			expected: lib.NewSet[string](),
		},
	}

	exportSetInternals := cmp.AllowUnexported(lib.NewSet[string]())

	for _, i := range testcases {
		t.Run(i.name, func(t *testing.T) {
			result, err := lib.Parse(strings.NewReader(i.input))

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(i.expected, result[0].Tags, exportSetInternals); diff != "" {
				t.Errorf("tag list mismatch:\n%s", diff)
			}
		})
	}
}

func TestExtractDate(t *testing.T) {
	type testcase struct {
		name     string
		input    string
		expected time.Time
	}

	testcases := []testcase{
		{
			name:     "with date",
			input:    "# 23.11.2025\n",
			expected: time.Date(2025, time.November, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "with invalid date",
			input:    "# 33.11.2025\n",
			expected: time.Time{},
		},
		{
			name:     "no date",
			input:    "# Note 1\n\nLine 1\n",
			expected: time.Time{},
		},
		{
			name:     "date mixed with text",
			input:    "# Test 23.11.2025 - test\n",
			expected: time.Time{},
		},
		{
			name:     "date in note body",
			input:    "# Note 1\n\nLine 23.11.2025\n",
			expected: time.Time{},
		},
		{
			name:     "child note dates",
			input:    "# Note 1\n\n ## 23.11.2025\n",
			expected: time.Time{},
		},
	}

	for _, i := range testcases {
		t.Run(i.name, func(t *testing.T) {
			result, err := lib.Parse(strings.NewReader(i.input))

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(i.expected, result[0].Date); diff != "" {
				t.Errorf("date mismatch:\n%v", diff)
			}
		})
	}
}

func TestUnclosedFencedBlock(t *testing.T) {
	input := "# Note 1\n\n```\ncode block"
	_, err := lib.Parse(strings.NewReader(input))

	if err == nil {
		t.Errorf("expected error for unclosed fenced code block")
	}

	expectedMsg := "invalid onefile content, fenced code block was not closed"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestExtractIcon(t *testing.T) {
	type testcase struct {
		name     string
		input    string
		expected string
	}

	testcases := []testcase{
		{
			name:     "no emoji",
			input:    "# Note 1\n\nLine 1\n",
			expected: "",
		},
		{
			name:     "emoji in title",
			input:    "# 🎉 Note 1\n",
			expected: "🎉",
		},
		{
			name:     "emoji in body",
			input:    "# Note 1\n\nLine 1 🎯\n",
			expected: "🎯",
		},
		{
			name:     "picks first emoji when multiple are included",
			input:    "# 🎉 Note 1\n\nLine 1 🎯\n",
			expected: "🎉",
		},
		{
			name:     "child note inherits no icon from parent",
			input:    "# 01.01.2026\n\n🎉 Line 1\n\n## Child Note 1\n\nLine 2\n",
			expected: "🎉",
		},
	}

	for _, i := range testcases {
		t.Run(i.name, func(t *testing.T) {
			result, err := lib.Parse(strings.NewReader(i.input))

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diff := cmp.Diff(i.expected, result[0].Icon); diff != "" {
				t.Errorf("icon mismatch:\n%s", diff)
			}
		})
	}
}
