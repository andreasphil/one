package markdown_test

import (
	"bytes"
	"testing"

	"github.com/andreasphil/one/lib/markdown"
	"github.com/andreasphil/one/lib/note"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

func TestWikiLinkExtension(t *testing.T) {
	// Stands in for the note lookup the web server passes in, falling back to
	// the plain slug for notes that don't exist.
	resolve := func(target string) (string, bool) {
		slugs := map[string]string{
			"Rehearsal":  "2026-06-27-rehearsal",
			"01.02.2026": "2026-02-01",
		}

		if slug, found := slugs[target]; found {
			return slug, true
		}

		return note.Slug(target), false
	}

	p := parser.New(parser.WithExtensions(markdown.WikiLinkParser))
	r := html.New(html.WithExtensions(markdown.NewWikiLinkHTMLRenderer("/notes/", resolve)))

	type testcase struct {
		name     string
		in       string
		expected string
	}

	tests := []testcase{
		{"simple", "[[a thing]]", `<p><a class="wikilink unresolved" href="/notes/a-thing/">a thing</a></p>`},
		{"title case", "[[A Thing]]", `<p><a class="wikilink unresolved" href="/notes/a-thing/">A Thing</a></p>`},
		{"unicode", "[[Äpfel]]", `<p><a class="wikilink unresolved" href="/notes/%C3%A4pfel/">Äpfel</a></p>`},
		{"punctuation", "[[Hello! World?]]", `<p><a class="wikilink unresolved" href="/notes/hello-world/">Hello! World?</a></p>`},
		{"surrounded by text", "see [[Rehearsal]] here", `<p>see <a class="wikilink" href="/notes/2026-06-27-rehearsal/">Rehearsal</a> here</p>`},
		{"two links", "[[Rehearsal]] and [[b]]", `<p><a class="wikilink" href="/notes/2026-06-27-rehearsal/">Rehearsal</a> and <a class="wikilink unresolved" href="/notes/b/">b</a></p>`},
		{"resolved child note", "[[Rehearsal]]", `<p><a class="wikilink" href="/notes/2026-06-27-rehearsal/">Rehearsal</a></p>`},
		{"resolved daily note", "[[01.02.2026]]", `<p><a class="wikilink" href="/notes/2026-02-01/">01.02.2026</a></p>`},
		{"empty", "[[]]", `<p>[[]]</p>`},
		{"unclosed", "[[unclosed", `<p>[[unclosed</p>`},
		{"single brackets", "[a thing]", `<p>[a thing]</p>`},
		{"markdown link", "[a thing](/url)", `<p><a href="/url">a thing</a></p>`},
		{"code span", "`[[a thing]]`", "<p><code>[[a thing]]</code></p>"},
		{"across lines", "[[a thing\n]]", "<p>[[a thing\n]]</p>"},
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
