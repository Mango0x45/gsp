package parser

import (
	"bytes"
	"fmt"

	"github.com/tdewolff/parse/v2"
)

type invalidSyntax struct {
	row      int
	col      int
	expected string
	found    string
}

func newInvalidSyntax(in *parse.Input, expected, found string) invalidSyntax {
	row, col, _ := parse.Position(bytes.NewReader(in.Bytes()), in.Offset())
	return invalidSyntax{row, col, expected, found}
}

func (e invalidSyntax) Error() string {
	return fmt.Sprintf("%d:%d: syntax error: expected %s but found %s",
		e.row, e.col-1, e.expected, e.found)
}

type invalidEscape struct {
	row  int
	col  int
	rune rune
}

func newInvalidEscape(in *parse.Input, got rune) invalidEscape {
	row, col, _ := parse.Position(bytes.NewReader(in.Bytes()), in.Offset())
	return invalidEscape{row, col, got}
}

func (e invalidEscape) Error() string {
	return fmt.Sprintf("%d:%d: invalid escape sequence: ‘\\%c’",
		e.row, e.col-1, e.rune)
}

type voidHasChildren struct {
	row int
	col int
	tag string
}

func newVoidHasChildren(in *parse.Input, tag string) voidHasChildren {
	row, col, _ := parse.Position(bytes.NewReader(in.Bytes()), in.Offset())
	return voidHasChildren{row, col, tag}
}

func (e voidHasChildren) Error() string {
	return fmt.Sprintf("%d:%d: void element ‘%s’ may not have any child nodes",
		e.row, e.col-1, e.tag)
}

type eof struct{}

func (e eof) Error() string {
	return "reached end of file while parsing"
}
