package markdown

import (
	"bytes"
	"fmt"
	"io"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

// Node ---------------------------------------------------

var KindWikiLink = ast.NewNodeKind("WikiLink")

type wikiLinkNode struct {
	ast.BaseInline
	Value text.SingleLineValue
}

func newWikiLinkNode(value text.SingleLineValue) *wikiLinkNode {
	n := &wikiLinkNode{Value: value}
	n.Init(n)
	return n
}

func (n *wikiLinkNode) Kind() ast.NodeKind {
	return KindWikiLink
}

func (n *wikiLinkNode) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(n, map[string]any{"Value": n.Value})
}

// Parser -------------------------------------------------

func NewWikiLinkParser() parser.Extension {
	return &wikiLinkParserExtension{}
}

var WikiLinkParser = NewWikiLinkParser()

type wikiLinkParserExtension struct{}

func (e *wikiLinkParserExtension) ParserOptions(_ *parser.Config) []parser.Option {
	return []parser.Option{
		parser.WithInlineParsers(util.Prioritized[parser.InlineParser](&wikiLinkParser{}, 150)),
	}
}

type wikiLinkParser struct{}

func (p *wikiLinkParser) Trigger() []byte {
	return []byte{'['}
}

func (p *wikiLinkParser) Parse(parent ast.Node, block text.Reader, context parser.Context) ast.Node {
	line, segment := block.PeekLine()

	if !bytes.HasPrefix(line, []byte("[[")) {
		return nil
	}

	// Length of the target, i.e. everything between the brackets. -1 is an
	// unclosed link, 0 an empty one.
	n := bytes.Index(line[2:], []byte("]]"))
	if n <= 0 {
		return nil
	}

	value := text.NewSegment(segment.Start+2, segment.Start+2+n)
	node := newWikiLinkNode(text.NewSingleLineValueFromSegment(value, block.Decoder()))

	block.Advance(2 + n + 2)
	return node
}

// Renderer -----------------------------------------------

func NewWikiLinkHTMLRenderer(prefix string, resolve func(target string) string) html.Extension {
	return &wikiLinkHTMLRendererExtension{prefix: prefix, resolve: resolve}
}

type wikiLinkHTMLRendererExtension struct {
	prefix  string
	resolve func(target string) string
}

func (e *wikiLinkHTMLRendererExtension) RendererOptions(_ *html.Config) []html.Option {
	return []html.Option{
		html.WithNodeRenderers(map[ast.NodeKind]html.NodeRenderer{
			KindWikiLink: html.NodeRendererFunc(e.render),
		}),
	}
}

func (e *wikiLinkHTMLRendererExtension) render(
	w io.Writer, src []byte, node ast.Node, entering bool, rc renderer.Context,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n, ok := node.(*wikiLinkNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("expected wiki link node, got %v", node)
	}

	bw := w.(util.BufWriter)

	_, _ = bw.WriteString(`<a class="wikilink" href="`)
	_, _ = bw.WriteString(e.prefix)
	_, _ = html.ContextLinkURLWriter(rc).WriteString(e.resolve(n.Value.Value(src)))
	_, _ = bw.WriteString(`/">`)
	_, _ = n.Value.WriteTo(html.ContextTextWriter(rc), src)
	_, _ = bw.WriteString("</a>")

	return ast.WalkSkipChildren, nil
}
