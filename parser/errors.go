package parser

import (
	"bytes"
	"fmt"

	"github.com/tdewolff/parse/v2"
)

// Location represents the location at which an error occured.
type Location struct {
	Path string
	Row  int
	Col  int
}

func (l Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.Path, l.Row, l.Col-1)
}

func locFromInput(in *parse.Input) Location {
	row, col, _ := parse.Position(bytes.NewReader(in.Bytes()), in.Offset())
	return Location{"", row, col}
}

// InvalidSyntaxError indicates that the parser encountered an
// unexpected token or character while evaluating the GSP document.
type InvalidSyntaxError struct {
	Where    Location
	Expected string
	Found    string
}

func newInvalidSyntaxError(in *parse.Input, expected, found string) InvalidSyntaxError {
	return InvalidSyntaxError{locFromInput(in), expected, found}
}

func (e InvalidSyntaxError) Error() string {
	return fmt.Sprintf("%s: syntax error: expected %s but found %s",
		e.Where, e.Expected, e.Found)
}

// InvalidEscapeError indicates that an invalid escape sequence was
// encountered within a text block or string literal.
type InvalidEscapeError struct {
	Where Location
	Rune  rune
}

func newInvalidEscapeError(in *parse.Input, got rune) InvalidEscapeError {
	return InvalidEscapeError{locFromInput(in), got}
}

func (e InvalidEscapeError) Error() string {
	return fmt.Sprintf("%s: invalid escape sequence: ‘\\%c’",
		e.Where, e.Rune)
}

// VoidHasChildrenError indicates that an HTML void element (such as
// img or br) was incorrectly given child nodes.
type VoidHasChildrenError struct {
	Where Location
	Tag   string
}

func newVoidHasChildrenError(in *parse.Input, tag string) VoidHasChildrenError {
	return VoidHasChildrenError{locFromInput(in), tag}
}

func (e VoidHasChildrenError) Error() string {
	return fmt.Sprintf("%s: void element ‘%s’ may not have any child nodes",
		e.Where, e.Tag)
}

// EOFError indicates that the parser reached the end of the file
// unexpectedly while parsing a construct.
type EOFError struct{}

func (e EOFError) Error() string {
	return "reached end of file while parsing"
}
