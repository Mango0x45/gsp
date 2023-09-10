package formatter

import (
	"fmt"
	"slices"
	"unicode"

	"git.thomasvoss.com/gsp/parser"
)

var (
	attrValueEscapes = map[rune]string{
		'"': "&quot;",
		'&': "&amp;",
		'<': "&lt;",
	}
	stringEscapes = map[rune]string {
		'"': "&quot;",
		'&': "&amp;",
		'<': "&lt;",
		'>': "&gt;",
		'\'': "&apos;",
	}
)

func PrintAst(ast parser.AstNode) {
	switch ast.Type {
	case parser.Text:
		printText(ast.Text)
	case parser.DocType:
		printDocType(ast)
	case parser.Normal:
		fmt.Printf("<%s", ast.Text)
		printAttrs(ast.Attrs)

		if len(ast.Children) == 0 {
			if parser.Xml {
				fmt.Print("/>")
			} else {
				fmt.Print(">")
			}
		} else {
			fmt.Print(">")
			printChildren(ast.Children)
			fmt.Printf("</%s>", ast.Text)
		}
	case parser.Tagless:
		printChildren(ast.Children)
	}
}

func printAttrs(attrs []parser.Attr) {
	classes := attrs
	classes = slices.DeleteFunc(classes, func (a parser.Attr) bool {
		return a.Key != "class"
	})
	attrs = slices.DeleteFunc(attrs, func (a parser.Attr) bool {
		return a.Key == "class"
	})

	if len(classes) > 0 {
		fmt.Print(" class=\"")
		for i, a := range classes {
			fmt.Print(a.Value)
			if i != len(classes) - 1 {
				fmt.Print(" ")
			} else {
				fmt.Print("\"")
			}
		}
	}

	for _, a := range attrs {
		fmt.Printf(" %s", a.Key)
		if a.Value != "" {
			fmt.Print("=\"")
			for _, r := range a.Value {
				if v, ok := attrValueEscapes[r]; ok {
					fmt.Print(v)
				} else {
					fmt.Printf("%c", r)
				}
			}
			fmt.Print("\"")
		}
	}
}

func printDocType(node parser.AstNode) {
	if parser.Xml {
		fmt.Print("<?xml")
		printAttrs(node.Attrs)
		fmt.Print("?>")
	} else {
		fmt.Print("<!DOCTYPE")
		printAttrs(node.Attrs)
		fmt.Print(">")
	}
}

func printText(s string) {
	for _, r := range s {
		if v, ok := stringEscapes[r]; ok {
			fmt.Print(v)
		} else {
			fmt.Printf("%c", r)
		}
	}
}

func printChildren(nodes []parser.AstNode) {
	for i, n := range nodes {
		if i == 0 && n.Type == parser.Text {
			n.Text = trimLeftSpaces(n.Text)
		}
		if i == len(nodes) - 1 && n.Type == parser.Text {
			n.Text = trimRightSpaces(n.Text)
		}

		PrintAst(n)
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
