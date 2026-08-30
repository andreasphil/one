// Package service holds the implementations the web server depends on.
package service

import (
	"bytes"
	"html/template"

	"github.com/andreasphil/one/lib/markdown"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

// Markdown renders the markdown source of a note as HTML, with GFM,
// typographer and #tag support enabled.
type Markdown struct {
	parser   parser.Parser
	renderer html.Renderer
}

// NewMarkdown creates a Markdown renderer.
func NewMarkdown() Markdown {
	p := parser.New(parser.WithExtensions(
		extension.GFMParser,
		extension.TypographerParser,

		markdown.TagParser,
	))

	r := html.New(html.WithExtensions(
		extension.GFMHTMLRenderer,

		markdown.NewTagHTMLRenderer("/tags/"),
	))

	return Markdown{parser: p, renderer: r}
}

// Render converts the markdown in input to HTML.
func (m Markdown) Render(input string) (template.HTML, error) {
	src := []byte(input)
	out := bytes.Buffer{}
	if err := m.renderer.Render(&out, src, m.parser.Parse(src)); err != nil {
		return "", err
	}

	return template.HTML(out.String()), nil
}
