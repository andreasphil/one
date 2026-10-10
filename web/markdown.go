package web

import (
	"bytes"
	gohtml "html"
	"html/template"
	"strings"

	"github.com/andreasphil/one/lib/markdown"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

type markdownRenderer struct {
	parser   parser.Parser
	renderer html.Renderer
}

func newMarkdownRenderer(resolveNote func(target string) (string, bool)) markdownRenderer {
	p := parser.New(
		parser.WithAutoHeadingID(),
		parser.WithExtensions(
			extension.GFMParser,
			extension.TypographerParser,
			extension.DefinitionListParser,

			markdown.CalloutParser,
			markdown.TagParser,
			markdown.WikiLinkParser,
		))

	r := html.New(html.WithExtensions(
		extension.GFMHTMLRenderer,
		extension.DefinitionListHTMLRenderer,

		markdown.CalloutHTMLRenderer,
		markdown.NewTagHTMLRenderer("/tags/"),
		markdown.NewWikiLinkHTMLRenderer("/notes/", resolveNote),
	))

	return markdownRenderer{parser: p, renderer: r}
}

func (m markdownRenderer) render(input string) (template.HTML, error) {
	src := []byte(input)
	out := bytes.Buffer{}
	if err := m.renderer.Render(&out, src, m.parser.Parse(src)); err != nil {
		return "", err
	}

	return template.HTML(out.String()), nil
}

const excerptWords = 30

// excerpt returns the first 30 words of the plain text of input. Longer text is
// truncated and ends with "...". Code blocks are skipped.
func (m markdownRenderer) excerpt(input string) string {
	src := []byte(input)
	var text strings.Builder

	ast.Walk(m.parser.Parse(src), func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			if _, ok := n.(ast.BlockNode); ok {
				text.WriteString(" ")
			}
			return ast.WalkContinue, nil
		}

		switch n := n.(type) {
		case *ast.Text:
			text.WriteString(gohtml.UnescapeString(n.Value.Value(src)))
			if n.SoftLineBreak() || n.HardLineBreak() {
				text.WriteString(" ")
			}
		case *ast.CodeSpan:
			text.WriteString(n.Value.Value(src))
		case *ast.AutoLink:
			text.WriteString(n.Label.Value(src))
		case *markdown.WikiLinkNode:
			text.WriteString(n.Value.Value(src))
		case *markdown.TagNode:
			text.WriteString("#")
			text.WriteString(n.Value.Value(src))
		}

		return ast.WalkContinue, nil
	})

	words := strings.Fields(text.String())
	if len(words) > excerptWords {
		return strings.Join(words[:excerptWords], " ") + "..."
	}

	return strings.Join(words, " ")
}
