package web

import (
	"bytes"
	"html/template"

	"github.com/andreasphil/one/lib/markdown"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

type markdownRenderer struct {
	parser   parser.Parser
	renderer html.Renderer
}

func newMarkdownRenderer(resolveNote func(target string) (string, bool)) markdownRenderer {
	p := parser.New(parser.WithExtensions(
		extension.GFMParser,
		extension.TypographerParser,
		extension.DefinitionListParser,

		markdown.TagParser,
		markdown.WikiLinkParser,
	))

	r := html.New(html.WithExtensions(
		extension.GFMHTMLRenderer,
		extension.DefinitionListHTMLRenderer,

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
