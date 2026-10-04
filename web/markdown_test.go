package web_test

import (
	"strings"
	"testing"
)

func TestExcerpt(t *testing.T) {
	type testcase struct {
		name     string
		content  string
		expected string
	}

	testcases := []testcase{
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "plain content",
			content:  "Just some plain text",
			expected: "Just some plain text",
		},
		{
			name:     "emphasis and strong",
			content:  "Some *emphasized* and **strong** and _underscored_ and __double__ words",
			expected: "Some emphasized and strong and underscored and double words",
		},
		{
			name:     "inline code and strikethrough",
			content:  "Run `go test` and ~~forget~~ it",
			expected: "Run go test and forget it",
		},
		{
			name:     "links and images",
			content:  "See [the docs](https://example.com) and ![a picture](pic.png) here",
			expected: "See the docs and a picture here",
		},
		{
			name:     "autolinks",
			content:  "Found at <https://example.com/page> today",
			expected: "Found at https://example.com/page today",
		},
		{
			name:     "headings, quotes and lists",
			content:  "## Section\n\n> A quote\n\n- First item\n* Second item\n+ Third item\n1. Fourth item",
			expected: "Section A quote First item Second item Third item Fourth item",
		},
		{
			name:     "nested list markers and checkboxes",
			content:  "- Outer\n  - Inner\n- [ ] Open\n- [x] Done",
			expected: "Outer Inner Open Done",
		},
		{
			name:     "callout with empty quote lines",
			content:  "> [!NOTE]\n>\n> Not responsible for singed eyebrows.\n>\n> Especially Beaker.",
			expected: "Not responsible for singed eyebrows. Especially Beaker.",
		},
		{
			name:     "snake case preserved across lines",
			content:  "previous_password: BorkBorkBork1\nnext_scheduled_change: soon",
			expected: "previous_password: BorkBorkBork1 next_scheduled_change: soon",
		},
		{
			name:     "underscore emphasis next to snake case",
			content:  "The _old_ value of previous_password is gone",
			expected: "The old value of previous_password is gone",
		},
		{
			name:     "hashtag preserved at start of line",
			content:  "Some content\n#tech 💻",
			expected: "Some content #tech 💻",
		},
		{
			name:     "tags",
			content:  "Filed under #work and #ideas",
			expected: "Filed under #work and #ideas",
		},
		{
			name:     "wiki links",
			content:  "See [[Reading list]] and [[01.02.2026]] for more",
			expected: "See Reading list and 01.02.2026 for more",
		},
		{
			name:     "unclosed wiki link left alone",
			content:  "See [[unclosed here",
			expected: "See [[unclosed here",
		},
		{
			name:     "table",
			content:  "| Fruit | Color |\n| --- | --- |\n| Apple | Green |\n| Plum | Purple |",
			expected: "Fruit Color Apple Green Plum Purple",
		},
		{
			name:     "table with alignment and no outer pipes",
			content:  "Fruit | Color\n:--- | ---:\nApple | Green",
			expected: "Fruit Color Apple Green",
		},
		{
			name:     "code blocks skipped and horizontal rules",
			content:  "Before\n\n```go\nfmt.Println()\n```\n\n---\n\nAfter",
			expected: "Before After",
		},
		{
			name:     "typographer substitutions",
			content:  `"Quoted" -- and... done`,
			expected: "“Quoted” – and… done",
		},
		{
			name:     "line breaks and whitespace normalized",
			content:  "One\n\n\nTwo\t\tthree   four\n",
			expected: "One Two three four",
		},
		{
			name:     "truncated to 30 words",
			content:  strings.Repeat("word ", 50),
			expected: strings.TrimSpace(strings.Repeat("word ", 30)) + "...",
		},
		{
			name:     "exactly 30 words kept",
			content:  strings.Repeat("word ", 30),
			expected: strings.TrimSpace(strings.Repeat("word ", 30)),
		},
	}

	renderer := newTestMarkdownRenderer(nil)

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := renderer.Excerpt(tc.content); got != tc.expected {
				t.Errorf("Excerpt(%q) = %q, want %q", tc.content, got, tc.expected)
			}
		})
	}
}
