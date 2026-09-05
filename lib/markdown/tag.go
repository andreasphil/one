// Package markdown provides goldmark extensions for rendering notes.
package markdown

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/text"
	"github.com/yuin/goldmark/v2/util"
)

// Node ---------------------------------------------------

var KindTag = ast.NewNodeKind("Tag")

type tagNode struct {
	ast.BaseInline
	Value text.SingleLineValue
}

func newTagNode(value text.SingleLineValue) *tagNode {
	n := &tagNode{Value: value}
	n.Init(n)
	return n
}

func (n *tagNode) Kind() ast.NodeKind {
	return KindTag
}

func (n *tagNode) Dump(_ []byte) *ast.NodeDump {
	return ast.NewNodeDump(n, map[string]any{"Value": n.Value})
}

// Parser -------------------------------------------------

func NewTagParser() parser.Extension {
	return &tagParserExtension{}
}

var TagParser = NewTagParser()

type tagParserExtension struct{}

func (e *tagParserExtension) ParserOptions(_ *parser.Config) []parser.Option {
	return []parser.Option{
		parser.WithInlineParsers(util.Prioritized[parser.InlineParser](&tagParser{}, 999)),
	}
}

type tagParser struct{}

func (p *tagParser) Trigger() []byte {
	return []byte{'#'}
}

func (p *tagParser) Parse(parent ast.Node, block text.Reader, context parser.Context) ast.Node {
	before := block.PrecedingCharacter()
	if !unicode.IsSpace(before) {
		return nil
	}

	line, segment := block.PeekLine()

	i := 1 // skip "#"
	for i < len(line) {
		r, size := utf8.DecodeRune(line[i:])
		if !isTagRune(r) {
			break
		}
		i += size
	}

	if i == 1 {
		return nil
	}

	value := text.NewSegment(segment.Start+1, segment.Start+i)
	node := newTagNode(text.NewSingleLineValueFromSegment(value, block.Decoder()))

	block.Advance(i)
	return node
}

func isTagRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// Renderer -----------------------------------------------

func NewTagHTMLRenderer(prefix string) html.Extension {
	return &tagHTMLRendererExtension{prefix: prefix}
}

type tagHTMLRendererExtension struct {
	prefix string
}

func (e *tagHTMLRendererExtension) RendererOptions(_ *html.Config) []html.Option {
	return []html.Option{
		html.WithNodeRenderers(map[ast.NodeKind]html.NodeRenderer{
			KindTag: html.NodeRendererFunc(e.render),
		}),
	}
}

func (e *tagHTMLRendererExtension) render(
	w io.Writer, src []byte, node ast.Node, entering bool, rc renderer.Context,
) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n, ok := node.(*tagNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("expected tag node, got %v", node)
	}

	bw := w.(util.BufWriter)

	_, _ = bw.WriteString(`<a class="tag" href="`)
	_, _ = bw.WriteString(e.prefix)
	_, _ = n.Value.WriteTo(html.ContextLinkURLWriter(rc), src)
	_, _ = bw.WriteString(`/">`)
	_, _ = n.Value.WriteTo(html.ContextTextWriter(rc), src)
	_, _ = bw.WriteString("</a>")

	return ast.WalkSkipChildren, nil
}
