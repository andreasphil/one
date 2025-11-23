package service

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type MarkdownService struct {
	renderer goldmark.Markdown
}

func NewMarkdownService() MarkdownService {
	markdown := goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Typographer))
	return MarkdownService{renderer: markdown}
}

func (m MarkdownService) RenderMarkdown(input string) (string, error) {
	out := bytes.Buffer{}
	if err := m.renderer.Convert([]byte(input), &out); err != nil {
		return "", err
	}

	return out.String(), nil
}
