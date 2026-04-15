package ast

type NodeType int

const (
	Normal NodeType = iota
	Comment
	Void
	Escapable
	Raw
	Text
	Macro
	VerbatimMacro
)

type Node struct {
	Type       NodeType
	Name       string
	Attributes map[string][]string
	Children   []Node
}
