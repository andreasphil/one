package lib_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/andreasphil/one/lib"
	"github.com/yuin/goldmark"
)

func TestNoteTitleExtension(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		class    string
		want     string
	}{
		{
			name:     "converts first H1 to note title",
			markdown: "# Hello World\n\nSome content",
			class:    "note-title",
			want:     `<h1><span class="note-title">Hello World</span></h1>`,
		},
		{
			name:     "custom class name",
			markdown: "# My Note",
			class:    "custom-class",
			want:     `<h1><span class="custom-class">My Note</span></h1>`,
		},
		{
			name:     "only first H1 is converted",
			markdown: "# First\n\n## Second\n\n# Third",
			class:    "note-title",
			want:     `<h1><span class="note-title">First</span></h1>`,
		},
		{
			name:     "handles H1 with inline formatting",
			markdown: "# **Bold** and *italic*",
			class:    "note-title",
			want:     `<h1><span class="note-title"><strong>Bold</strong> and <em>italic</em></span></h1>`,
		},
		{
			name:     "no H1 present",
			markdown: "## Just H2\n\nSome text",
			class:    "note-title",
			want:     `<h2>Just H2</h2>`,
		},
		{
			name:     "empty H1",
			markdown: "#\n\nContent",
			class:    "note-title",
			want:     `<h1><span class="note-title"></span></h1>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithExtensions(lib.NewNoteTitleExtension(tt.class)),
			)

			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.markdown), &buf); err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			got := strings.TrimSpace(buf.String())
			if !strings.Contains(got, tt.want) {
				t.Errorf("NoteTitleExtension output mismatch:\ngot:\n%s\n\nwant to contain:\n%s", got, tt.want)
			}
		})
	}
}
