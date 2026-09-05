package markdown_test

import (
	"bytes"
	"testing"

	"github.com/andreasphil/one/lib/markdown"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

func TestTagExtension(t *testing.T) {
	p := parser.New(parser.WithExtensions(markdown.TagParser))
	r := html.New(html.WithExtensions(markdown.NewTagHTMLRenderer("/tags/")))

	type testcase struct {
		name     string
		in       string
		expected string
	}

	tests := []testcase{
		{"simple", "#world", `<p><a class="tag" href="/tags/world/">#world</a></p>`},
		{"underscore", "#my_tag", `<p><a class="tag" href="/tags/my_tag/">#my_tag</a></p>`},
		{"unicode", "#grüße", `<p><a class="tag" href="/tags/gr%C3%BC%C3%9Fe/">#grüße</a></p>`},
		{"mid-word", "foo#bar", `<p>foo#bar</p>`},
		{"url anchor", "/page#anchor", `<p>/page#anchor</p>`},
		{"code span", "`#code`", "<p><code>#code</code></p>"},
		{"heading", "# Heading", `<h1>Heading</h1>`},
		{"bare hash", "foo #", `<p>foo #</p>`},
		{"in parens", "(#foo)", `<p>(#foo)</p>`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			src := []byte(tc.in)
			if err := r.Render(&buf, src, p.Parse(src)); err != nil {
				t.Fatal(err)
			}

			result := bytes.TrimSpace(buf.Bytes())
			if string(result) != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}
