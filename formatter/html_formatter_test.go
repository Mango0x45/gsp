package formatter

import (
	"strings"
	"testing"

	"git.thomasvoss.com/gsp/ast"
)

func TestWriteAst(t *testing.T) {
	tests := []struct {
		name    string
		nodes   []ast.Node
		opts    Options
		want    string
		wantErr bool
	}{
		{
			name: "Doctype generation",
			nodes: []ast.Node{
				{
					Type: ast.Normal,
					Name: "html",
				},
			},
			opts: Options{Doctype: true},
			want: "<!DOCTYPE html><html></html>",
		},
		{
			name: "Normal node with an attribute and escapable characters",
			nodes: []ast.Node{
				{
					Type: ast.Normal,
					Name: "p",
					Attributes: map[string][]string{
						"id": {"x-p"},
					},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: "Hello <world> & \"friends\"",
						},
					},
				},
			},
			opts: Options{},
			want: `<p id="x-p">Hello &lt;world&gt; &amp; &#34;friends&#34;</p>`,
		},
		{
			name: "Void node",
			nodes: []ast.Node{
				{
					Type: ast.Void,
					Name: "img",
					Attributes: map[string][]string{
						"src": {"image.png"},
					},
				},
			},
			opts: Options{},
			want: `<img src="image.png">`,
		},
		{
			name: "Raw node (script)",
			nodes: []ast.Node{
				{
					Type: ast.Raw,
					Name: "script",
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: `const a = "a" < "b";`,
						},
					},
				},
			},
			opts: Options{},
			want: `<script>const a = "a" < "b";</script>`,
		},
		{
			name: "Comment node with comments enabled",
			nodes: []ast.Node{
				{
					Type: ast.Comment,
					Name: "/",
					Children: []ast.Node{
						{
							Type: ast.Normal,
							Name: "div",
						},
					},
				},
			},
			opts: Options{Comments: true},
			want: `<!-- <div></div> -->`,
		},
		{
			name: "Comment node with comments disabled",
			nodes: []ast.Node{
				{
					Type: ast.Comment,
					Name: "/",
					Children: []ast.Node{
						{
							Type: ast.Normal,
							Name: "div",
						},
					},
				},
			},
			opts: Options{Comments: false},
			want: ``,
		},
		{
			name: "Empty attribute",
			nodes: []ast.Node{
				{
					Type: ast.Normal,
					Name: "input",
					Attributes: map[string][]string{
						"disabled": {""},
					},
				},
			},
			opts: Options{},
			want: `<input disabled></input>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			err := WriteAst(&buf, tt.nodes, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteAst() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("WriteAst() = %q, want %q", got, tt.want)
			}
		})
	}
}
