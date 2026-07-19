package service

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type MarkdownService struct {
	renderer goldmark.Markdown
}

func NewMarkdownService() MarkdownService {
	markdown := goldmark.New(goldmark.WithExtensions(
		extension.GFM,
		extension.Typographer,
	))

	return MarkdownService{renderer: markdown}
}

func (m MarkdownService) Render(input string) (template.HTML, error) {
	out := bytes.Buffer{}
	if err := m.renderer.Convert([]byte(input), &out); err != nil {
		return "", err
	}

	return template.HTML(out.String()), nil
}
