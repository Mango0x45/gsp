// Package ast defines the AST for the GSP markup language.
//
// It provides the core data structures used to represent parsed GSP
// documents, and offers utility functions that act on ASTs.
package ast

// NodeType represents the specific kind of a GSP AST node.
type NodeType int

const (
	// Normal represents a standard GSP node which may contain
	// attributes and child nodes.
	Normal NodeType = iota
	// Comment represents a comment node.
	Comment
	// Void represents an HTML void element (e.g., br, img, meta) that
	// cannot contain any child nodes.
	Void
	// Escapable represents a node whose content undergoes HTML
	// escaping, but which cannot contain child nodes (like title or
	// textarea).
	Escapable
	// Raw represents a node whose content is raw and unparsed (e.g.,
	// script, style).
	Raw
	// Text represents a text node containing string content.
	Text
	// Macro represents a GSP macro node, prefixed with ‘$’, which
	// expands dynamically during formatting.
	Macro
	// VerbatimMacro represents a verbatim GSP macro node, prefixed
	// with ‘$$’, whose output is not processed or escaped.
	VerbatimMacro
)

// Node represents a single element in a GSP AST.
type Node struct {
	// Type specifies this node’s type.
	Type NodeType
	// Name holds this node’s tag name, macro name, or text content depending on the Type.
	Name string
	// Attributes contains this node’s attributes mapping keys to slices of values.
	Attributes map[string][]string
	// Children contains this node’s descendant nodes.
	Children []Node
}
