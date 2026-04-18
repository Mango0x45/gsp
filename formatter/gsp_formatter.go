package formatter

import (
	"cmp"
	"fmt"
	"io"
	"maps"
	"strings"

	"git.thomasvoss.com/gsp/v4/ast"
	g_strconv "git.thomasvoss.com/gsp/v4/strconv"
)

// WriteUntranslatedAST serializes a GSP abstract syntax tree back
// into valid GSP markup and writes it to the provided io.Writer.
// This is primarily used to stream node bodies to external macro
// executables.
//
// Node that the output markup may not be 1:1 identical to the
// original input from which the AST was parsed.  The only guarantee
// is that both the original source and the output of this function
// are semantically equivalant.
func WriteUntranslatedAST(out io.Writer, ast []ast.Node) error {
	for _, n := range ast {
		if err := writeUntranslatedNode(out, n); err != nil {
			return err
		}
	}
	return nil
}

func writeUntranslatedNode(out io.Writer, node ast.Node) error {
	var e1, e2, e3 error

	switch node.Type {
	case ast.Comment:
		_, e1 = fmt.Fprint(out, "/ ")
		e2 = writeUntranslatedNode(out, node.Children[0])
	case ast.Macro:
		_, e1 = fmt.Fprint(out, "$")
		e2 = writeUntranslatedTag(out, node)
		e3 = writeUntranslatedBody(out, node)
	case ast.VerbatimMacro:
		_, e1 = fmt.Fprint(out, "$$")
		e2 = writeUntranslatedTag(out, node)
		e3 = writeUntranslatedBody(out, node)
	case ast.Normal, ast.Escapable, ast.Void:
		e1 = writeUntranslatedTag(out, node)
		e2 = writeUntranslatedBody(out, node)
	case ast.Raw:
		e1 = writeUntranslatedTag(out, node)
		e2 = writeUntranslatedRawBody(out, node)
	case ast.Text:
		e1 = writeUntranslatedText(out, node.Name)
	}

	return cmp.Or(e1, e2, e3)
}

func writeUntranslatedTag(out io.Writer, node ast.Node) error {
	if _, err := fmt.Fprintf(out, "%s ", node.Name); err != nil {
		return err
	}
	for k, vs := range maps.All(node.Attributes) {
		v := g_strconv.EscapeString(strings.Join(vs, " "))
		if _, err := fmt.Fprintf(out, `%s="%s" `, k, v); err != nil {
			return err
		}
	}
	return nil
}

func writeUntranslatedBody(out io.Writer, node ast.Node) error {
	if len(node.Children) != 0 && node.Children[0].Type == ast.Text {
		if _, err := fmt.Fprint(out, "{="); err != nil {
			return err
		}

		for i, n := range node.Children {
			/* In a text block, ‘real’ nodes are always at odd indices */
			if i&1 == 1 {
				if _, err := fmt.Fprint(out, "@"); err != nil {
					return err
				}
			}
			if err := writeUntranslatedNode(out, n); err != nil {
				return err
			}
		}
	} else {
		if _, err := fmt.Fprint(out, "{"); err != nil {
			return err
		}

		if err := WriteUntranslatedAST(out, node.Children); err != nil {
			return err
		}
	}

	_, err := fmt.Fprint(out, "}")
	return err
}

func writeUntranslatedRawBody(out io.Writer, node ast.Node) error {
	_, err := fmt.Fprintf(out, "{%s}", node.Children[0].Name)
	return err
}

func writeUntranslatedText(out io.Writer, s string) error {
	_, err := fmt.Fprint(out, s)
	return err
}
