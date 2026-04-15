package formatter

import (
	"cmp"
	"fmt"
	"html"
	"io"
	"maps"
	"strings"

	"git.thomasvoss.com/gsp/ast"
)

type Options struct {
	Comments   bool
	Doctype    bool
	SearchPath []string
}

func WriteAst(out io.Writer, ast []ast.Node, opts Options) error {
	if opts.Doctype {
		_, err := fmt.Fprint(out, "<!DOCTYPE html>")
		if err != nil {
			return err
		}
	}
	return writeNodes(out, ast, opts)
}

func writeNodes(out io.Writer, ast []ast.Node, opts Options) error {
	for _, n := range ast {
		if err := writeNode(out, n, opts); err != nil {
			return err
		}
	}
	return nil
}

func writeNode(out io.Writer, node ast.Node, opts Options) error {
	var e1, e2, e3 error

	switch node.Type {
	case ast.Comment:
		if opts.Comments {
			e1 = writeCommentStart(out)
			e2 = writeNodes(out, node.Children, opts)
			e3 = writeCommentEnd(out)
		}
	case ast.Macro, ast.VerbatimMacro:
		path, ok := findMacro(node.Name, opts.SearchPath)
		if !ok {
			return fmt.Errorf("%s: failed to find macro", node.Name)
		}
		e1 = execMacro(out, path, node, opts)
	case ast.Normal, ast.Escapable:
		e1 = writeOpenTag(out, node)
		e2 = writeNodes(out, node.Children, opts)
		e3 = writeCloseTag(out, node)
	case ast.Void:
		e1 = writeOpenTag(out, node)
	case ast.Raw:
		e1 = writeOpenTag(out, node)
		e2 = writeRawText(out, node.Children[0].Name)
		e3 = writeCloseTag(out, node)
	case ast.Text:
		e1 = writeText(out, html.EscapeString(node.Name))
	}

	return cmp.Or(e1, e2, e3)
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

func writeRawText(out io.Writer, s string) error {
	_, err := out.Write([]byte(s))
	return err
}

func writeText(out io.Writer, s string) error {
	bs := []byte(s)
	for i := 0; i < len(bs); i++ {
		if bs[i] == '\\' {
			i++
		}
		if _, err := out.Write([]byte{bs[i]}); err != nil {
			return err
		}
	}
	return nil
}
