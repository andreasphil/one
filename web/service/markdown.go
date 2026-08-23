// Package service holds the implementations the web server depends on.
package service

import (
	"bytes"
	"html/template"

	"github.com/andreasphil/one/lib/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// Markdown renders the markdown source of a note as HTML, with GFM,
// typographer and #tag support enabled.
type Markdown struct {
	renderer goldmark.Markdown
}

// NewMarkdown creates a Markdown renderer.
func NewMarkdown() Markdown {
	md := goldmark.New(goldmark.WithExtensions(
		extension.GFM,
		extension.Typographer,

		&markdown.Tag{Prefix: "/tags/"},
	))

	return Markdown{renderer: md}
}

// Render converts the markdown in input to HTML.
func (m Markdown) Render(input string) (template.HTML, error) {
	out := bytes.Buffer{}
	if err := m.renderer.Convert([]byte(input), &out); err != nil {
		return "", err
	}

	return template.HTML(out.String()), nil
}
