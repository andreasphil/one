package lib

import (
	"html"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type noteTitleNode struct {
	ast.BaseBlock
	class string
}

var noteTitleNodeKind = ast.NewNodeKind("NoteTitle")

func (n *noteTitleNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Class": n.class}, nil)
}

func (n *noteTitleNode) Kind() ast.NodeKind {
	return noteTitleNodeKind
}

type noteTitleParser struct {
	class string
}

func (t *noteTitleParser) Transform(node *ast.Document, reader text.Reader, c parser.Context) {
	var h1 *ast.Heading

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := n.(*ast.Heading); ok && heading.Level == 1 && h1 == nil {
			h1 = heading
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	if h1 != nil {
		noteTitle := &noteTitleNode{class: t.class}

		for child := h1.FirstChild(); child != nil; {
			next := child.NextSibling()
			noteTitle.AppendChild(noteTitle, child)
			child = next
		}

		h1.Parent().ReplaceChild(h1.Parent(), h1, noteTitle)
	}
}

type noteTitleRenderer struct{}

func (r *noteTitleRenderer) RegisterFuncs(renderer renderer.NodeRendererFuncRegisterer) {
	renderer.Register(noteTitleNodeKind, r.renderNoteTitle)
}

func (r *noteTitleRenderer) renderNoteTitle(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*noteTitleNode)

	if entering {
		w.WriteString(`<h1><span class="` + n.class + `">`)
	} else {
		w.WriteString(`</span></h1>` + "\n")
	}

	return ast.WalkContinue, nil
}

type NoteTitleExtension struct {
	Class string
}

func NewNoteTitleExtension(class string) *NoteTitleExtension {
	return &NoteTitleExtension{Class: html.EscapeString(class)}
}

func (e *NoteTitleExtension) Extend(goldmark goldmark.Markdown) {
	goldmark.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(&noteTitleParser{class: e.Class}, 0),
	))

	goldmark.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&noteTitleRenderer{}, 0),
	))
}
