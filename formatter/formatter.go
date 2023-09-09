package formatter

import (
	"fmt"
	"unicode"

	"git.thomasvoss.com/gsp/parser"
)

var stringEscapes = map[rune]string{
	'"': "&quot;",
	'&': "&amp;",
	'<': "&lt;",
}

func PrintHtml(ast parser.AstNode) {
	if ast.Type == parser.Text {
		fmt.Print(ast.Text)
		return
	}

	if ast.Type == parser.Normal {
		fmt.Printf("<%s", ast.Text)

		// Classes are grouped together with ‘class="…"’, so we need
		// special handling.
		classes := []string{}
		notClasses := []parser.Attr{}

		for _, a := range ast.Attrs {
			if a.Key == "class" {
				classes = append(classes, a.Value)
			} else {
				notClasses = append(notClasses, a)
			}
		}

		if len(classes) > 0 {
			fmt.Printf(" class=\"%s", classes[0])
			for _, c := range classes[1:] {
				fmt.Printf(" %s", c)
			}
			fmt.Print("\"")
		}

		for _, a := range notClasses {
			fmt.Printf(" %s", a.Key)
			if a.Value == "" {
				continue
			}
			fmt.Print("=\"")
			for _, r := range a.Value {
				if v, ok := stringEscapes[r]; ok {
					fmt.Print(v)
				} else {
					fmt.Printf("%c", r)
				}
			}
			fmt.Print("\"")
		}

		fmt.Print(">")
	}

	if len(ast.Children) == 0 {
		return
	}

	for i, n := range ast.Children {
		if n.Type == parser.Text {
			if i == 0 {
				n.Text = trimLeftSpaces(n.Text)
			}

			if i == len(ast.Children)-1 {
				n.Text = trimRightSpaces(n.Text)
			}
		}

		PrintHtml(n)
	}

	if ast.Type == parser.Normal {
		fmt.Printf("</%s>", ast.Text)
	}
}

func trimLeftSpaces(s string) string {
	i := 0
	rs := []rune(s)
	for i < len(s) && unicode.IsSpace(rs[i]) {
		i++
	}
	return string(rs[i:])
}

func trimRightSpaces(s string) string {
	i := len(s) - 1
	rs := []rune(s)
	for i >= 0 && unicode.IsSpace(rs[i]) {
		i--
	}
	return string(rs[:i+1])
}
