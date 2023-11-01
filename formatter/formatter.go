package formatter

import (
	"fmt"
	"unicode"

	"git.thomasvoss.com/gsp/parser"
)

var (
	attrValueEscapes = map[rune]string{
		'"': "&quot;",
		'&': "&amp;",
		'<': "&lt;",
	}
	stringEscapes = map[rune]string{
		'"':  "&quot;",
		'&':  "&amp;",
		'<':  "&lt;",
		'>':  "&gt;",
		'\'': "&apos;",
	}
)

func PrintAst(ast parser.AstNode) {
	switch ast.Type {
	case parser.Text:
		printText(ast.Text)
	case parser.Normal:
		fmt.Printf("<%s", ast.Text)
		printAttrs(ast.Attrs)
		fmt.Print(">")

		if len(ast.Children) > 0 {
			printChildren(ast.Children)
			fmt.Printf("</%s>", ast.Text)
		}
	case parser.Tagless:
		printChildren(ast.Children)
	case parser.TaglessTrim:
		printChildrenTrim(ast.Children)
	}
}

func printAttrs(attrs []parser.Attr) {
	classes := make([]parser.Attr, 0, cap(attrs))
	nClasses := make([]parser.Attr, 0, cap(attrs))

	for _, a := range attrs {
		if a.Key == "class" {
			classes = append(classes, a)
		} else {
			nClasses = append(nClasses, a)
		}
	}

	if len(classes) > 0 {
		fmt.Print(" class=\"")
		for i, a := range classes {
			printAttrVal(a.Value)
			if i != len(classes)-1 {
				fmt.Print(" ")
			} else {
				fmt.Print("\"")
			}
		}
	}

	for _, a := range nClasses {
		fmt.Printf(" %s", a.Key)
		if a.Value != "" {
			fmt.Print("=\"")
			printAttrVal(a.Value)
			fmt.Print("\"")
		}
	}
}

func printAttrVal(s string) {
	for _, r := range s {
		if v, ok := attrValueEscapes[r]; ok {
			fmt.Print(v)
		} else {
			fmt.Printf("%c", r)
		}
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
	for _, n := range nodes {
		PrintAst(n)
	}
}

func printChildrenTrim(nodes []parser.AstNode) {
	for i, n := range nodes {
		if i == 0 && n.Type == parser.Text {
			n.Text = trimLeftSpaces(n.Text)
		}
		if i == len(nodes)-1 && n.Type == parser.Text {
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
	rs := []rune(s)
	i := len(rs) - 1
	for i >= 0 && unicode.IsSpace(rs[i]) {
		i--
	}
	return string(rs[:i+1])
}
