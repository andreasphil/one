package markdown_test

import (
	"bytes"
	"testing"

	"github.com/andreasphil/one/lib/markdown"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

func TestCalloutExtension(t *testing.T) {
	p := parser.New(parser.WithExtensions(markdown.CalloutParser))
	r := html.New(html.WithExtensions(markdown.CalloutHTMLRenderer))

	type testcase struct {
		name     string
		input    string
		expected string
	}

	testcases := []testcase{
		{"note", "> [!NOTE]\n> Text", "<div class=\"callout callout-note\">\n<header>Note</header>\n<p>Text</p>\n</div>"},
		{"tip", "> [!TIP]\n> Text", "<div class=\"callout callout-tip\">\n<header>Tip</header>\n<p>Text</p>\n</div>"},
		{"important", "> [!IMPORTANT]\n> Text", "<div class=\"callout callout-important\">\n<header>Important</header>\n<p>Text</p>\n</div>"},
		{"warning", "> [!WARNING]\n> Text", "<div class=\"callout callout-warning\">\n<header>Warning</header>\n<p>Text</p>\n</div>"},
		{"caution", "> [!CAUTION]\n> Text", "<div class=\"callout callout-caution\">\n<header>Caution</header>\n<p>Text</p>\n</div>"},
		{"lowercase", "> [!note]\n> Text", "<div class=\"callout callout-note\">\n<header>Note</header>\n<p>Text</p>\n</div>"},
		{"trailing spaces", "> [!NOTE]  \n> Text", "<div class=\"callout callout-note\">\n<header>Note</header>\n<p>Text</p>\n</div>"},
		{"multiple lines", "> [!NOTE]\n> One\n> Two", "<div class=\"callout callout-note\">\n<header>Note</header>\n<p>One\nTwo</p>\n</div>"},
		{"inline markup", "> [!NOTE]\n> Some *text*", "<div class=\"callout callout-note\">\n<header>Note</header>\n<p>Some <em>text</em></p>\n</div>"},
		{"multiple blocks", "> [!NOTE]\n> One\n>\n> - Two", "<div class=\"callout callout-note\">\n<header>Note</header>\n<p>One</p>\n<ul>\n<li>Two</li>\n</ul>\n</div>"},
		{"blank line after marker", "> [!NOTE]\n>\n> Text", "<div class=\"callout callout-note\">\n<header>Note</header>\n<p>Text</p>\n</div>"},
		{"nested", "- > [!NOTE]\n  > Text", "<ul>\n<li>\n<div class=\"callout callout-note\">\n<header>Note</header>\n<p>Text</p>\n</div>\n</li>\n</ul>"},
		{"empty", "> [!NOTE]", "<blockquote>\n<p>[!NOTE]</p>\n</blockquote>"},
		{"unknown type", "> [!FOO]\n> Text", "<blockquote>\n<p>[!FOO]\nText</p>\n</blockquote>"},
		{"title after marker", "> [!NOTE] Title\n> Text", "<blockquote>\n<p>[!NOTE] Title\nText</p>\n</blockquote>"},
		{"marker not on first line", "> Text\n> [!NOTE]", "<blockquote>\n<p>Text\n[!NOTE]</p>\n</blockquote>"},
		{"plain blockquote", "> Text", "<blockquote>\n<p>Text</p>\n</blockquote>"},
		{"outside blockquote", "[!NOTE]\nText", "<p>[!NOTE]\nText</p>"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			src := []byte(tc.input)
			if err := r.Render(&buf, src, p.Parse(src)); err != nil {
				t.Fatalf("Render(%q) error = %v", tc.input, err)
			}

			if got := string(bytes.TrimSpace(buf.Bytes())); got != tc.expected {
				t.Errorf("Render(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}
