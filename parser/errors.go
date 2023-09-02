package parser

import "fmt"

type invalidSyntax struct {
	pos      position
	expected string
	found    string
}

func (e invalidSyntax) Error() string {
	return fmt.Sprintf("Syntax error near %v; expected %s but found %s", e.pos, e.expected, e.found)
}

type eof struct{}

func (e eof) Error() string {
	return "Hit end-of-file while parsing.  You’re probably missing a closing brace (‘}’) somewhere"
}
