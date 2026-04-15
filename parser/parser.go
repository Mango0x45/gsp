package parser

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
	"github.com/tdewolff/parse/v2/js"

	"git.thomasvoss.com/gsp/ast"
)

var tagToType = map[string]ast.NodeType{
	"area":     ast.Void,
	"base":     ast.Void,
	"br":       ast.Void,
	"col":      ast.Void,
	"embed":    ast.Void,
	"hr":       ast.Void,
	"img":      ast.Void,
	"input":    ast.Void,
	"link":     ast.Void,
	"meta":     ast.Void,
	"param":    ast.Void,
	"script":   ast.Raw,
	"source":   ast.Void,
	"style":    ast.Raw,
	"textarea": ast.Escapable,
	"title":    ast.Escapable,
	"track":    ast.Void,
	"wbr":      ast.Void,
}

func Parse(r io.Reader) ([]ast.Node, error) {
	in := parse.NewInput(r)
	var nodes []ast.Node

	for {
		n, err := parseNode(in)
		switch {
		case err == io.EOF:
			return []ast.Node{}, eof{}
		case err != nil:
			return []ast.Node{}, err
		}
		nodes = append(nodes, n)

		switch err = skipSpaces(in); {
		case err == io.EOF:
			return nodes, nil
		case err != nil:
			return []ast.Node{}, err
		}
	}
}

func parseNode(in *parse.Input) (ast.Node, error) {
	if in.Peek(0) == '/' {
		in.Move(1)
		if err := skipSpaces(in); err != nil {
			return ast.Node{}, err
		}
		n, err := parseNode(in)
		if err != nil {
			return ast.Node{}, err
		}
		return ast.Node{
			Type:     ast.Comment,
			Name:     "/",
			Children: []ast.Node{n},
		}, nil
	}

	name, err := parseIdent(in, false)
	if err != nil {
		return ast.Node{}, err
	}

	ty, ok := tagToType[name]
	if !ok {
		ty = ast.Normal
	}

	var kids []ast.Node
	attrs := make(map[string][]string)

outer:
	for {
		if err := skipSpaces(in); err != nil {
			return ast.Node{}, err
		}

		ch, n := in.PeekRune(0)
		if ch == 0 && in.Err() != nil {
			return ast.Node{}, in.Err()
		}

		switch {
		case ch == '#':
			in.Move(n)
			in.Skip()
			sh, err := parseShorthand(in)
			if err != nil {
				return ast.Node{}, err
			}
			attrs["id"] = append(attrs["id"], sh)
		case ch == '.':
			in.Move(n)
			in.Skip()
			sh, err := parseShorthand(in)
			if err != nil {
				return ast.Node{}, err
			}
			attrs["class"] = append(attrs["class"], sh)
		case validNameStartChar(ch):
			k, v, err := parseAttribute(in)
			if err != nil {
				return ast.Node{}, err
			}
			attrs[k] = append(attrs[k], v)
		case ch == '{':
			in.Move(n)
			in.Skip()
			switch name {
			case "style":
				s, err := parseCSSBody(in)
				if err != nil {
					return ast.Node{}, err
				}
				kids = []ast.Node{ast.Node{
					Type: ast.Text,
					Name: s,
				}}
				break outer
			case "script":
				s, err := parseJSBody(in)
				if err != nil {
					return ast.Node{}, err
				}
				kids = []ast.Node{ast.Node{
					Type: ast.Text,
					Name: s,
				}}
				break outer
			default:
				if ch := in.Peek(0); ch == '-' || ch == '=' {
					in.Move(1)
					in.Skip()
					kids, err = parseTextBlock(in, ch == '=')
					if err != nil {
						return ast.Node{}, err
					}
					break outer
				}

				for {
					if err := skipSpaces(in); err != nil {
						return ast.Node{}, err
					}

					ch, n := in.PeekRune(0)
					if ch == 0 && in.Err() != nil {
						return ast.Node{}, in.Err()
					}

					if ch == '}' {
						in.Move(n)
						in.Skip()
						break outer
					}

					node, err := parseNode(in)
					if err != nil {
						return ast.Node{}, err
					}
					kids = append(kids, node)
				}
			}
		default:
			return ast.Node{}, newInvalidSyntax(in,
				"node attributes or braces",
				fmt.Sprintf("invalid character ‘%c’", ch))
		}
	}

	if name[0] == '$' {
		ty = ast.Macro
		name = name[1:]
	}

	return ast.Node{
		Type:       ty,
		Name:       name,
		Attributes: attrs,
		Children:   kids,
	}, nil
}

func parseIdent(in *parse.Input, attr bool) (string, error) {
	in.Skip()
	r, n := in.PeekRune(0)
	if r == 0 && in.Err() != nil {
		return "", in.Err()
	}

	if !validNameStartChar(r) {
		expected := "node name"
		if attr {
			expected = "attribute name"
		}
		return "", newInvalidSyntax(in,
			expected,
			fmt.Sprintf("invalid character ‘%c’", r))
	} else if !attr && !validNameChar(r) {
		return "", newInvalidSyntax(in,
			"class/id shorthand",
			fmt.Sprintf("invalid character ‘%c’", r))
	}

	in.Move(n)
	for {
		r, n = in.PeekRune(0)
		if r == 0 && in.Err() != nil {
			break
		}
		if !validNameChar(r) {
			break
		}
		in.Move(n)
	}

	s := string(in.Shift())
	if !attr && s == "$" {
		return "", newInvalidSyntax(in,
			"macro name", "nothing")
	}
	return s, nil
}

func parseShorthand(in *parse.Input) (string, error) {
	r, n := in.PeekRune(0)
	if r == 0 && in.Err() != nil {
		return "", in.Err()
	}

	if !validNameChar(r) {
		return "", newInvalidSyntax(in,
			"id/class identifier",
			fmt.Sprintf("invalid character ‘%c’", r))
	}

	in.Move(n)
	for {
		r, n = in.PeekRune(0)
		if r == 0 && in.Err() != nil {
			break
		}
		if !validNameChar(r) {
			break
		}
		in.Move(n)
	}

	return string(in.Shift()), nil
}

func parseAttribute(in *parse.Input) (string, string, error) {
	k, err := parseIdent(in, true)
	if err != nil {
		return "", "", err
	}

	r, n := in.PeekRune(0)
	if r == 0 && in.Err() != nil {
		return "", "", in.Err()
	} else if r != '=' {
		return k, "", nil
	}

	in.Move(n)
	in.Skip()

	v, err := parseString(in)
	if err != nil {
		return "", "", err
	}

	return k, v, nil
}

func parseCSSBody(in *parse.Input) (string, error) {
	off := in.Offset()
	depth := 1
	l := css.NewLexer(in)

	for {
		tt, _ := l.Next()
		if tt == css.ErrorToken {
			if l.Err() == io.EOF {
				return "", eof{}
			}
			return "", l.Err()
		}

		if tt == css.LeftBraceToken {
			depth++
		} else if tt == css.RightBraceToken {
			depth--
			if depth == 0 {
				s := string(in.Bytes()[off : in.Offset()-1])
				return s, nil
			}
		}
	}
}

func parseJSBody(in *parse.Input) (string, error) {
	off := in.Offset()
	depth := 1
	l := js.NewLexer(in)

	for {
		tt, _ := l.Next()
		if tt == js.ErrorToken {
			if l.Err() == io.EOF {
				return "", eof{}
			}
			return "", l.Err()
		}

		if tt == js.OpenBraceToken {
			depth++
		} else if tt == js.CloseBraceToken {
			depth--
			if depth == 0 {
				s := string(in.Bytes()[off : in.Offset()-1])
				return s, nil
			}
		}
	}
}

func parseTextBlock(in *parse.Input, untrimmed bool) ([]ast.Node, error) {
	depth := 1
	nodes := make([]ast.Node, 0, 8)

	in.Skip()
outer:
	for {
		ch := in.Peek(0)
		in.Move(1)

		switch ch {
		case 0:
			if in.Err() == io.EOF {
				return []ast.Node{}, eof{}
			}
			return []ast.Node{}, in.Err()
		case '@':
			in.Move(-1)
			nodes = append(nodes, ast.Node{
				Type: ast.Text,
				Name: string(in.Shift()),
			})
			in.Move(1)
			n, err := parseNode(in)
			if err != nil {
				return []ast.Node{}, err
			}
			nodes = append(nodes, n)
			in.Skip()
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				in.Move(-1)
				nodes = append(nodes, ast.Node{
					Type: ast.Text,
					Name: string(in.Shift()),
				})
				in.Move(1)
				break outer
			}
		case '\\':
			switch ch := in.Peek(0); ch {
			/* Ignore escaping EOF so that we throw a syntax error instead */
			case 0, '@', '{', '}', '\\':
			default:
				return []ast.Node{}, newInvalidEscape(in, rune(ch))
			}
			in.Move(1)
		}
	}

	if !untrimmed {
		l := len(nodes) - 1
		nodes[0].Name = strings.TrimLeftFunc(nodes[0].Name, unicode.IsSpace)
		nodes[l].Name = strings.TrimRightFunc(nodes[l].Name, unicode.IsSpace)
	}

	return nodes, nil
}

func parseString(in *parse.Input) (string, error) {
	r, n := in.PeekRune(0)
	if r == 0 && in.Err() != nil {
		return "", in.Err()
	} else if r != '"' {
		return "", newInvalidSyntax(in,
			"double-quoted string",
			fmt.Sprintf("‘%c’", r))
	}
	in.Move(n)
	in.Skip()

	var sb strings.Builder
	for {
		r, n := in.PeekRune(0)
		if r == 0 && in.Err() != nil {
			return "", in.Err()
		}
		in.Move(n)

		switch r {
		case '"':
			in.Skip()
			return sb.String(), nil
		case '\\':
			r2, n2 := in.PeekRune(0)
			if r2 == 0 && in.Err() != nil {
				return "", in.Err()
			}
			in.Move(n2)

			if r2 != '\\' && r2 != '"' {
				return "", newInvalidEscape(in, r2)
			}
			sb.WriteRune(r2)
		default:
			sb.WriteRune(r)
		}
		in.Skip()
	}
}

func skipSpaces(in *parse.Input) error {
	for {
		r, n := in.PeekRune(0)
		if r == 0 && in.Err() != nil {
			return in.Err()
		}
		if unicode.IsSpace(r) {
			in.Move(n)
		} else {
			in.Skip()
			return nil
		}
	}
}

func validNameStartChar(r rune) bool {
	return r == '$' || r == ':' || r == '_' ||
		(r >= 'A' && r <= 'Z') ||
		(r >= 'a' && r <= 'z') ||
		(r >= 0x000C0 && r <= 0x000D6) ||
		(r >= 0x000D8 && r <= 0x000F6) ||
		(r >= 0x000F8 && r <= 0x002FF) ||
		(r >= 0x00370 && r <= 0x0037D) ||
		(r >= 0x0037F && r <= 0x01FFF) ||
		(r >= 0x0200C && r <= 0x0200D) ||
		(r >= 0x02070 && r <= 0x0218F) ||
		(r >= 0x02C00 && r <= 0x02FEF) ||
		(r >= 0x03001 && r <= 0x0D7FF) ||
		(r >= 0x0F900 && r <= 0x0FDCF) ||
		(r >= 0x0FDF0 && r <= 0x0FFFD) ||
		(r >= 0x10000 && r <= 0xEFFFF)
}

func validNameChar(r rune) bool {
	return validNameStartChar(r) ||
		r == '-' || r == '.' || r == '·' ||
		(r >= '0' && r <= '9') ||
		(r >= 0x0300 && r <= 0x036F) ||
		(r >= 0x203F && r <= 0x2040)
}
