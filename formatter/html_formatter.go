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

func WriteAst(w io.Writer, path string, ast []ast.Node, opts Options) error {
	if opts.Doctype {
		_, err := fmt.Fprint(w, "<!DOCTYPE html>")
		if err != nil {
			return err
		}
	}
	return writeNodes(w, path, ast, opts)
}

func writeNodes(w io.Writer, path string, ast []ast.Node, opts Options) error {
	for _, n := range ast {
		if err := writeNode(w, path, n, opts); err != nil {
			return err
		}
	}
	return nil
}

func writeNode(w io.Writer, path string, node ast.Node, opts Options) error {
	var e1, e2, e3 error

	switch node.Type {
	case ast.Comment:
		if opts.Comments {
			e1 = writeCommentStart(w)
			e2 = writeNodes(w, path, node.Children, opts)
			e3 = writeCommentEnd(w)
		}
	case ast.Macro, ast.VerbatimMacro:
		mpath, ok := findMacro(node.Name, opts.SearchPath)
		if !ok {
			return fmt.Errorf("%s: failed to find macro", node.Name)
		}
		e1 = execMacro(w, mpath, path, node, opts)
	case ast.Normal, ast.Escapable:
		e1 = writeOpenTag(w, node)
		e2 = writeNodes(w, path, node.Children, opts)
		e3 = writeCloseTag(w, node)
	case ast.Void:
		e1 = writeOpenTag(w, node)
	case ast.Raw:
		e1 = writeOpenTag(w, node)
		e2 = writeRawText(w, node.Children[0].Name)
		e3 = writeCloseTag(w, node)
	case ast.Text:
		e1 = writeText(w, html.EscapeString(node.Name))
	}

	return cmp.Or(e1, e2, e3)
}

func writeOpenTag(w io.Writer, node ast.Node) error {
	_, err := fmt.Fprintf(w, "<%s", node.Name)
	if err != nil {
		return err
	}

	for k, vs := range maps.All(node.Attributes) {
		v := html.EscapeString(strings.Join(vs, " "))
		if len(v) == 0 {
			_, err = fmt.Fprintf(w, ` %s`, k)
		} else {
			_, err = fmt.Fprintf(w, ` %s="%s"`, k, v)
		}

		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, ">")
	return err
}

func writeCloseTag(w io.Writer, node ast.Node) error {
	_, err := fmt.Fprintf(w, "</%s>", node.Name)
	return err
}

func writeCommentStart(w io.Writer) error {
	_, err := fmt.Fprint(w, "<!-- ")
	return err
}

func writeCommentEnd(w io.Writer) error {
	_, err := fmt.Fprint(w, " -->")
	return err
}

func writeRawText(w io.Writer, s string) error {
	_, err := w.Write([]byte(s))
	return err
}

func writeText(w io.Writer, s string) error {
	bs := []byte(s)
	for i := 0; i < len(bs); i++ {
		if bs[i] == '\\' {
			i++
		}
		if _, err := w.Write([]byte{bs[i]}); err != nil {
			return err
		}
	}
	return nil
}
