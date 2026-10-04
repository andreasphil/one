package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

var calloutTitles = map[string]string{
	"note":      "Note",
	"tip":       "Tip",
	"important": "Important",
	"warning":   "Warning",
	"caution":   "Caution",
}

// Node ---------------------------------------------------

var KindCallout = ast.NewNodeKind("Callout")

type CalloutNode struct {
	ast.BaseBlock
	Type string
}

func newCalloutNode(calloutType string) *CalloutNode {
	n := &CalloutNode{Type: calloutType}
	n.Init(n)
	return n
}

func (n *CalloutNode) Kind() ast.NodeKind {
	return KindCallout
}

func (n *CalloutNode) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(n, map[string]any{"Type": n.Type})
}

// Parser -------------------------------------------------

func NewCalloutParser() parser.Extension {
	return &calloutParserExtension{}
}

var CalloutParser = NewCalloutParser()

type calloutParserExtension struct{}

func (e *calloutParserExtension) ParserOptions(_ *parser.Config) []parser.Option {
	return []parser.Option{
		parser.WithASTTransformers(util.Prioritized[parser.ASTTransformer](&calloutTransformer{}, 999)),
	}
}

type calloutTransformer struct{}

func (t *calloutTransformer) Transform(doc *ast.Document, reader text.Reader, _ parser.Context) {
	src := reader.Source()

	var blockquotes []*ast.Blockquote
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if bq, ok := n.(*ast.Blockquote); ok && entering {
			blockquotes = append(blockquotes, bq)
		}
		return ast.WalkContinue, nil
	})

	for _, bq := range blockquotes {
		p, ok := bq.FirstChild().(*ast.Paragraph)
		if !ok {
			continue
		}

		calloutType, marker, ok := parseCalloutMarker(p, src)
		if !ok {
			continue
		}

		// GitHub keeps blockquotes without content after the marker as they are
		if len(marker) == p.ChildCount() && bq.ChildCount() == 1 {
			continue
		}

		for _, n := range marker {
			p.RemoveChild(n)
		}
		if !p.HasChildren() {
			bq.RemoveChild(p)
		}

		callout := newCalloutNode(calloutType)
		for bq.HasChildren() {
			callout.AppendChild(bq.FirstChild())
		}
		bq.Parent().ReplaceChild(bq, callout)
	}
}

// parseCalloutMarker checks whether the first line of p is a callout marker
// like "[!NOTE]". If it is, it returns the callout type and the inline nodes
// that make up the marker.
func parseCalloutMarker(p *ast.Paragraph, src []byte) (string, []ast.Node, bool) {
	var line strings.Builder
	var marker []ast.Node

	for n := range p.Children() {
		t, ok := n.(*ast.Text)
		if !ok {
			return "", nil, false
		}

		line.WriteString(t.Value.Value(src))
		marker = append(marker, n)

		if t.SoftLineBreak() || t.HardLineBreak() {
			break
		}
	}

	s := strings.TrimSpace(line.String())
	if !strings.HasPrefix(s, "[!") || !strings.HasSuffix(s, "]") {
		return "", nil, false
	}

	calloutType := strings.ToLower(s[2 : len(s)-1])
	if _, ok := calloutTitles[calloutType]; !ok {
		return "", nil, false
	}

	return calloutType, marker, true
}

// Renderer -----------------------------------------------

func NewCalloutHTMLRenderer() html.Extension {
	return &calloutHTMLRendererExtension{}
}

var CalloutHTMLRenderer = NewCalloutHTMLRenderer()

type calloutHTMLRendererExtension struct{}

func (e *calloutHTMLRendererExtension) RendererOptions(_ *html.Config) []html.Option {
	return []html.Option{
		html.WithNodeRenderers(map[ast.NodeKind]html.NodeRenderer{
			KindCallout: html.NodeRendererFunc(e.render),
		}),
	}
}

func (e *calloutHTMLRendererExtension) render(
	w io.Writer, _ []byte, node ast.Node, entering bool, _ renderer.Context,
) (ast.WalkStatus, error) {
	n, ok := node.(*CalloutNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("expected callout node, got %v", node)
	}

	bw := w.(util.BufWriter)

	if !entering {
		_, _ = bw.WriteString("</div>\n")
		return ast.WalkContinue, nil
	}

	_, _ = bw.WriteString(`<div class="callout callout-`)
	_, _ = bw.WriteString(n.Type)
	_, _ = bw.WriteString("\">\n<header>")
	_, _ = bw.WriteString(calloutTitles[n.Type])
	_, _ = bw.WriteString("</header>\n")

	return ast.WalkContinue, nil
}
