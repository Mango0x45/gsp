package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type nodeType uint

const (
	Normal nodeType = iota
	Tagless
	Text
)

type Attr struct {
	Key   string
	Value string
}

type AstNode struct {
	Type     nodeType
	Text     string
	Attrs    []Attr
	Children []AstNode
}

func ParseFile(file *os.File) (AstNode, error) {
	r := reader{r: bufio.NewReader(file)}
	return r.parseNode()
}

func (reader *reader) parseNode() (AstNode, error) {
	if err := reader.skipSpaces(); err != nil {
		return AstNode{}, err
	}

	if r, err := reader.peekRune(); err != nil {
		return AstNode{}, err
	} else if r == '-' {
		return reader.parseText()
	}

	node := AstNode{}
	if name, err := reader.parseNodeName(); err != nil {
		return AstNode{}, err
	} else {
		node.Type = Normal
		node.Text = name
	}

	if attrs, err := reader.parseAttrs(); err != nil {
		return AstNode{}, err
	} else {
		node.Attrs = attrs
	}

	// The above call to reader.parseAttrs() guarantees that we have the ‘{’
	// token.
	if _, err := reader.readRune(); err != nil {
		return AstNode{}, err
	}

loop:
	for {
		if err := reader.skipSpaces(); err != nil {
			return AstNode{}, err
		}

		if r, err := reader.peekRune(); err == io.EOF {
			return AstNode{}, eof{}
		} else if err != nil {
			return AstNode{}, err
		} else if r == '}' {
			break loop
		}

		if n, err := reader.parseNode(); err != nil {
			return AstNode{}, err
		} else {
			node.Children = append(node.Children, n)
		}
	}

	// The above loop guarantees that we have the ‘}’ token.
	if _, err := reader.readRune(); err != nil {
		return AstNode{}, err
	}

	return node, nil
}

func (reader *reader) parseNodeName() (string, error) {
	var r rune
	var err error

	if err = reader.skipSpaces(); err != nil {
		return "", err
	}

	sb := strings.Builder{}

	if r, err = reader.readRune(); err != nil {
		return "", err
	} else if !validNameStartChar(r) {
		return "", invalidSyntax{
			pos:      reader.pos,
			expected: "node name",
			found:    fmt.Sprintf("invalid character ‘%c’", r),
		}
	}

	for validNameChar(r) {
		sb.WriteRune(r)
		if r, err = reader.readRune(); err != nil {
			return "", err
		}
	}

	if err = reader.unreadRune(); err != nil {
		return "", err
	}
	return sb.String(), nil
}

func (reader *reader) parseText() (AstNode, error) {
	if _, err := reader.readRune(); err != nil {
		return AstNode{}, err
	}

	sb := strings.Builder{}
	node := AstNode{Type: Tagless}

loop:
	for {
		r, err := reader.readRune()
		if err != nil {
			return AstNode{}, err
		}
		switch r {
		case '}':
			if err := reader.unreadRune(); err != nil {
				return AstNode{}, err
			}
			break loop
		case '@':
			node.Children = append(node.Children, AstNode{
				Type: Text,
				Text: sb.String(),
			})
			sb = strings.Builder{}

			n, err := reader.parseNode()
			if err != nil {
				return AstNode{}, err
			}
			node.Children = append(node.Children, n)
		case '\\':
			r, err = reader.readRune()
			if err != nil {
				return AstNode{}, err
			}
			if r != '\\' && r != '@' && r != '}' {
				return AstNode{}, invalidSyntax{
					pos:      reader.pos,
					expected: "valid escape sequence (‘\\\\’, ‘\\@’, or ‘\\}’)",
					found:    fmt.Sprintf("‘\\%c’", r),
				}
			}
			fallthrough
		default:
			sb.WriteRune(r)
		}
	}

	node.Children = append(node.Children, AstNode{
		Type: Text,
		Text: sb.String(),
	})
	return node, nil
}

func (reader *reader) parseAttrs() ([]Attr, error) {
	attrs := make([]Attr, 0, 2)

loop:
	for {
		if err := reader.skipSpaces(); err != nil {
			return nil, err
		}
		r, err := reader.peekRune()
		if err != nil {
			return nil, err
		}

		attr := Attr{}
		switch r {
		case '{':
			break loop
		case '.':
			fallthrough
		case '#':
			sym := r

			// Skip ‘sym’
			if _, err := reader.readRune(); err != nil {
				return nil, err
			}

			if s, err := reader.parseNodeName(); err != nil {
				return nil, err
			} else {
				attr.Value = s
				if sym == '.' {
					attr.Key = "class"
				} else {
					attr.Key = "id"
				}
			}
		default:
			if unicode.IsSpace(r) {
				if err := reader.skipSpaces(); err != nil {
					return nil, err
				}
				continue
			}

			if s, err := reader.parseNodeName(); err != nil {
				return nil, err
			} else {
				attr.Key = s
			}

			if r, err := reader.readNonSpaceRune(); err != nil {
				return nil, err
			} else if r != '=' {
				reader.unreadRune()
				break
			}

			if s, err := reader.parseString(); err != nil {
				return nil, err
			} else {
				attr.Value = s
			}
		}
		attrs = append(attrs, attr)
	}

	return attrs, nil
}

func (reader *reader) parseString() (string, error) {
	sb := strings.Builder{}

	if r, err := reader.readNonSpaceRune(); err != nil {
		return "", err
	} else if r != '"' {
		return "", invalidSyntax{
			pos:      reader.pos,
			expected: "double-quoted string",
			found:    fmt.Sprintf("‘%c’", r),
		}
	}

	for {
		r, err := reader.readRune()
		if err != nil {
			return "", err
		}

		switch r {
		case '"':
			return sb.String(), nil
		case '\\':
			r, err := reader.readRune()
			if err != nil {
				return "", err
			}

			if r != '\\' && r != '"' {
				return "", invalidSyntax{
					pos:      reader.pos,
					expected: "valid escape sequence (‘\\\\’ or ‘\\\"’)",
					found:    fmt.Sprintf("‘\\%c’", r),
				}
			}

			sb.WriteRune(r)
		default:
			sb.WriteRune(r)
		}
	}
}

func validNameStartChar(r rune) bool {
	return r == ':' || r == '_' ||
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
