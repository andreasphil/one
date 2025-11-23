package adapter

type MarkdownRenderer interface {
	RenderMarkdown(input string) (string, error)
}
