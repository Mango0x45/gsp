package formatter

import (
	"fmt"
	"html"
	"io"
	"maps"
	"strings"

	"git.thomasvoss.com/gsp/ast"
)

func WriteAst(out io.Writer, ast []ast.Node) error {
	for _, n := range ast {
		if err := writeNode(out, n); err != nil {
			return err
		}
	}
	return nil
}

func writeNode(out io.Writer, node ast.Node) error {
	var e1, e2, e3 error

	switch node.Type {
	case ast.Comment:
		e1 = writeCommentStart(out)
		e2 = WriteAst(out, node.Children)
		e3 = writeCommentEnd(out)
	case ast.Normal, ast.Escapable:
		e1 = writeOpenTag(out, node)
		e2 = WriteAst(out, node.Children)
		e3 = writeCloseTag(out, node)
	case ast.Void:
		e1 = writeOpenTag(out, node)
	case ast.Raw:
		e1 = writeOpenTag(out, node)

		/* This is only reached by <script> and <style> tags, where we
		   assume either 0 or 1 children, and where the child (if it
		   exists) is a text node */
		if len(node.Children) == 1 {
			e2 = writeText(out, node.Children[0].Name)
		}

		e3 = writeCloseTag(out, node)
	case ast.Text:
		e1 = writeText(out, html.EscapeString(node.Name))
	}

	if e1 != nil {
		return e1
	}
	if e2 != nil {
		return e2
	}
	return e3
}

func writeOpenTag(out io.Writer, node ast.Node) error {
	_, err := fmt.Fprintf(out, "<%s", node.Name)
	if err != nil {
		return err
	}

	for k, vs := range maps.All(node.Attributes) {
		v := html.EscapeString(strings.Join(vs, " "))
		if len(v) == 0 {
			_, err = fmt.Fprintf(out, ` %s`, k)
		} else {
			_, err = fmt.Fprintf(out, ` %s="%s"`, k, v)
		}

		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(out, ">")
	return err
}

func writeCloseTag(out io.Writer, node ast.Node) error {
	_, err := fmt.Fprintf(out, "</%s>", node.Name)
	return err
}

func writeCommentStart(out io.Writer) error {
	_, err := fmt.Fprint(out, "<!-- ")
	return err
}

func writeCommentEnd(out io.Writer) error {
	_, err := fmt.Fprint(out, " -->")
	return err
}

func writeText(out io.Writer, s string) error {
	_, err := out.Write([]byte(s))
	return err
}
