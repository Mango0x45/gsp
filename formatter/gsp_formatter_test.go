package formatter

import (
	"strings"
	"testing"

	"git.thomasvoss.com/gsp/ast"
)

func TestWriteUntranslatedAST(t *testing.T) {
	tests := []struct {
		name    string
		nodes   []ast.Node
		want    string
		wantErr bool
	}{
		{
			name: "Basic normal node",
			nodes: []ast.Node{
				{
					Type: ast.Normal,
					Name: "div",
				},
			},
			want: `div {}`,
		},
		{
			name: "Node with single attribute",
			nodes: []ast.Node{
				{
					Type: ast.Normal,
					Name: "a",
					Attributes: map[string][]string{
						"href": {"https://example.com"},
					},
				},
			},
			want: `a href="https://example.com" {}`,
		},
		{
			name: "Node with escaped attribute",
			nodes: []ast.Node{
				{
					Type: ast.Normal,
					Name: "span",
					Attributes: map[string][]string{
						"data-text": {`some "quoted" \ text`},
					},
				},
			},
			want: `span data-text="some \"quoted\" \\ text" {}`,
		},
		{
			name: "Text node sequence triggering {= block",
			nodes: []ast.Node{
				{
					Type: ast.Normal,
					Name: "p",
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: "Hello ",
						},
						{
							Type: ast.Normal,
							Name: "em",
							Children: []ast.Node{
								{
									Type: ast.Text,
									Name: "world",
								},
							},
						},
						{
							Type: ast.Text,
							Name: "!",
						},
					},
				},
			},
			want: `p {=Hello @em {=world}!}`,
		},
		{
			name: "Raw node",
			nodes: []ast.Node{
				{
					Type: ast.Raw,
					Name: "script",
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: "\n  const x = 1;\n",
						},
					},
				},
			},
			want: "script {\n  const x = 1;\n}",
		},
		{
			name: "Comment node",
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
			want: `/ div {}`,
		},
		{
			name: "Macro node",
			nodes: []ast.Node{
				{
					Type: ast.Macro,
					Name: "date",
					Attributes: map[string][]string{
						"format": {"%Y"},
					},
					Children: []ast.Node{},
				},
			},
			want: `$date format="%Y" {}`,
		},
		{
			name: "Verbatim macro node",
			nodes: []ast.Node{
				{
					Type: ast.VerbatimMacro,
					Name: "syntax_highlight",
					Attributes: map[string][]string{
						"lang": {"c"},
					},
					Children: []ast.Node{},
				},
			},
			want: `$$syntax_highlight lang="c" {}`,
		},
		{
			name: "Nested regular block",
			nodes: []ast.Node{
				{
					Type: ast.Normal,
					Name: "body",
					Children: []ast.Node{
						{
							Type: ast.Normal,
							Name: "div",
						},
					},
				},
			},
			want: `body {div {}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			err := WriteUntranslatedAST(&buf, tt.nodes)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteUntranslatedAST() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("WriteUntranslatedAST() = %q, want %q", got, tt.want)
			}
		})
	}
}
