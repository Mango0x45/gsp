package parser

import (
	"reflect"
	"strings"
	"testing"

	"git.thomasvoss.com/gsp/v4/ast"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []ast.Node
		wantErr bool
	}{
		{
			name:  "Basic node",
			input: `div {}`,
			want: []ast.Node{
				{
					Type:       ast.Normal,
					Name:       "div",
					Attributes: map[string][]string{},
				},
			},
		},
		{
			name:  "Attributes and shorthands",
			input: `p #my-id .class1 .class2 key="value" noval {}`,
			want: []ast.Node{
				{
					Type: ast.Normal,
					Name: "p",
					Attributes: map[string][]string{
						"id":    {"my-id"},
						"class": {"class1", "class2"},
						"key":   {"value"},
						"noval": {""},
					},
				},
			},
		},
		{
			name:  "Void element",
			input: `img src="test.png" {}`,
			want: []ast.Node{
				{
					Type: ast.Void,
					Name: "img",
					Attributes: map[string][]string{
						"src": {"test.png"},
					},
				},
			},
		},
		{
			name:  "Nested elements",
			input: `html { body { p {} } }`,
			want: []ast.Node{
				{
					Type:       ast.Normal,
					Name:       "html",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type:       ast.Normal,
							Name:       "body",
							Attributes: map[string][]string{},
							Children: []ast.Node{
								{
									Type:       ast.Normal,
									Name:       "p",
									Attributes: map[string][]string{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Trimmed text block",
			input: `p {-   trimmed text   }`,
			want: []ast.Node{
				{
					Type:       ast.Normal,
					Name:       "p",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: "trimmed text",
						},
					},
				},
			},
		},
		{
			name:  "Untrimmed text block",
			input: `pre {=   untrimmed text   }`,
			want: []ast.Node{
				{
					Type:       ast.Normal,
					Name:       "pre",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: "   untrimmed text   ",
						},
					},
				},
			},
		},
		{
			name:  "Embedded nodes in text",
			input: `p {- Hello @em{-world}! }`,
			want: []ast.Node{
				{
					Type:       ast.Normal,
					Name:       "p",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: "Hello ",
						},
						{
							Type:       ast.Normal,
							Name:       "em",
							Attributes: map[string][]string{},
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
		},
		{
			name:  "Macro node",
			input: `$date format="%Y" {}`,
			want: []ast.Node{
				{
					Type: ast.Macro,
					Name: "date",
					Attributes: map[string][]string{
						"format": {"%Y"},
					},
				},
			},
		},
		{
			name:  "Verbatim macro node",
			input: `$$syntax_highlight lang="c" {-extern int optind;}`,
			want: []ast.Node{
				{
					Type: ast.VerbatimMacro,
					Name: "syntax_highlight",
					Attributes: map[string][]string{
						"lang": {"c"},
					},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: "extern int optind;",
						},
					},
				},
			},
		},
		{
			name:  "Comment node",
			input: `/ div { p {} }`,
			want: []ast.Node{
				{
					Type:       ast.Comment,
					Name:       "/",
					Attributes: nil,
					Children: []ast.Node{
						{
							Type:       ast.Normal,
							Name:       "div",
							Attributes: map[string][]string{},
							Children: []ast.Node{
								{
									Type:       ast.Normal,
									Name:       "p",
									Attributes: map[string][]string{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "CSS style block",
			input: `style { body { /* foo */ color: red; } }`,
			want: []ast.Node{
				{
					Type:       ast.Raw,
					Name:       "style",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: " body { /* foo */ color: red; } ",
						},
					},
				},
			},
		},
		{
			name:  "JS script block",
			input: `script { /* foo */ const a = { b: 1 }; }`,
			want: []ast.Node{
				{
					Type:       ast.Raw,
					Name:       "script",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: " /* foo */ const a = { b: 1 }; ",
						},
					},
				},
			},
		},
		{
			name:  "Escaped characters in text block",
			input: `p {- \@ \{ \} \\ }`,
			want: []ast.Node{
				{
					Type:       ast.Normal,
					Name:       "p",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: `\@ \{ \} \\`,
						},
					},
				},
			},
		},
		{
			name:    "Invalid syntax – missing brace",
			input:   `div {`,
			wantErr: true,
		},
		{
			name:    "Invalid syntax – malformed attribute",
			input:   `div key= {}`,
			wantErr: true,
		},
		{
			name:    "Invalid syntax – invalid escape",
			input:   `p {- \n }`,
			wantErr: true,
		},
		{
			name:  "Complex JS block with braces in strings and comments",
			input: `script { const obj = { a: "}" }; /* } */ }`,
			want: []ast.Node{
				{
					Type:       ast.Raw,
					Name:       "script",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: ` const obj = { a: "}" }; /* } */ `,
						},
					},
				},
			},
		},
		{
			name:  "Boolean and empty attributes",
			input: `input disabled readonly="" {}`,
			want: []ast.Node{
				{
					Type: ast.Void,
					Name: "input",
					Attributes: map[string][]string{
						"disabled": {""},
						"readonly": {""},
					},
				},
			},
		},
		{
			name:  "Embedded macro in text block",
			input: `p {- Today is @$date{} }`,
			want: []ast.Node{
				{
					Type:       ast.Normal,
					Name:       "p",
					Attributes: map[string][]string{},
					Children: []ast.Node{
						{
							Type: ast.Text,
							Name: "Today is ",
						},
						{
							Type:       ast.Macro,
							Name:       "date",
							Attributes: map[string][]string{},
						},
						{
							Type: ast.Text,
							Name: "",
						},
					},
				},
			},
		},
		{
			name:  "Multiple IDs and classes shorthands",
			input: `div #primary #secondary .btn .btn-large {}`,
			want: []ast.Node{
				{
					Type: ast.Normal,
					Name: "div",
					Attributes: map[string][]string{
						"id":    {"primary", "secondary"},
						"class": {"btn", "btn-large"},
					},
				},
			},
		},
		{
			name:  "Empty text blocks",
			input: `div {-} p {=}`,
			want: []ast.Node{
				{
					Type:       ast.Normal,
					Name:       "div",
					Attributes: map[string][]string{},
					Children:   []ast.Node{{Type: ast.Text, Name: ""}},
				},
				{
					Type:       ast.Normal,
					Name:       "p",
					Attributes: map[string][]string{},
					Children:   []ast.Node{{Type: ast.Text, Name: ""}},
				},
			},
		},
		{
			name:  "Nested comments",
			input: `/ / div {}`,
			want: []ast.Node{
				{
					Type:       ast.Comment,
					Name:       "/",
					Attributes: nil,
					Children: []ast.Node{
						{
							Type:       ast.Comment,
							Name:       "/",
							Attributes: nil,
							Children: []ast.Node{
								{
									Type:       ast.Normal,
									Name:       "div",
									Attributes: map[string][]string{},
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "Custom tag with hyphens",
			input: `custom-element custom-attr="val" {}`,
			want: []ast.Node{
				{
					Type: ast.Normal,
					Name: "custom-element",
					Attributes: map[string][]string{
						"custom-attr": {"val"},
					},
				},
			},
		},
		{
			name:    "Invalid macro name syntax",
			input:   `$ {}`,
			wantErr: true,
		},
		{
			name:    "Invalid verbatim macro name syntax",
			input:   `$$ {}`,
			wantErr: true,
		},
		{
			name:    "Invalid shorthand syntax",
			input:   `div . {}`,
			wantErr: true,
		},
		{
			name:    "Invalid children in void element",
			input:   `meta charset="UTF-8" { title {- Hello! } }`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(tt.input), "<string>")
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() \ngot  = %v\nwant = %v", got, tt.want)
			}
		})
	}
}
