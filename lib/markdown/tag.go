package markdown

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Tag is a goldmark extension that turns hashtags such as "#example" into
// links. A tag starts with "#" preceded by whitespace or the start of a line,
// followed by one or more letters, digits, or underscores. It renders as an
// anchor with the class "tag", e.g. `<a class="tag" href="/tags/example">#example</a>`.
type Tag struct {
	// Prefix is prepended to the tag to build the href, e.g. "/tags/"
	Prefix string
}

// Extend registers the tag parser and renderer on m. It implements
// goldmark.Extender, so a Tag can be passed to goldmark.WithExtensions.
func (e *Tag) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(util.Prioritized(newTagParser(), 999)))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(util.Prioritized(newTagRenderer(e.Prefix), 999)))
}

// Node ---------------------------------------------------

var tagKind = ast.NewNodeKind("Tag")

type tagNode struct {
	ast.BaseInline
	// Segment covers the tag text without the leading "#"
	Segment text.Segment
}

// Value returns the tag text for the given source
func (n *tagNode) Value(src []byte) []byte {
	return n.Segment.Value(src)
}

func (n *tagNode) Kind() ast.NodeKind {
	return tagKind
}

func (n *tagNode) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"Value": fmt.Sprintf(`"%s"`, n.Value(src)),
	}, nil)
}

// Parser -------------------------------------------------

type tagParser struct{}

func newTagParser() parser.InlineParser {
	return &tagParser{}
}

func (p *tagParser) Trigger() []byte {
	return []byte{'#'}
}

func (p *tagParser) Parse(parent ast.Node, block text.Reader, context parser.Context) ast.Node {
	before := block.PrecendingCharacter()
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

	node := &tagNode{
		Segment: text.NewSegment(segment.Start+1, segment.Start+i),
	}

	block.Advance(i)
	return node
}

func isTagRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// Renderer -----------------------------------------------

type tagRenderer struct {
	Prefix string
}

func newTagRenderer(prefix string) renderer.NodeRenderer {
	return &tagRenderer{Prefix: prefix}
}

func (r *tagRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(tagKind, r.render)
}

func (r *tagRenderer) render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n, ok := node.(*tagNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("expected tag node, got %v", node)
	}

	value := n.Value(src)

	_, _ = w.WriteString(`<a class="tag" href="`)
	_, _ = w.WriteString(r.Prefix)
	_, _ = w.Write(util.URLEscape(value, false))
	_, _ = w.WriteString(`">#`)
	_, _ = w.Write(util.EscapeHTML(value))
	_, _ = w.WriteString("</a>")

	return ast.WalkSkipChildren, nil
}
