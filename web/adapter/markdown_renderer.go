package adapter

import "html/template"

type MarkdownRenderer interface {
	Render(input string) (template.HTML, error)
}
