package ast

import "errors"

// SkipChildren is used as a return value from WalkFunc to indicate
// that the node’s children in the AST should be skipped.  It is not
// returned as an error by any function.
var SkipChildren = errors.New("skip children")

// StopTraversal is used as a return value from WalkFunc to indicate
// that the walker should cease to recursively walk the AST.  It is
// not returned as an error by any function.
//
// This is useful for avoiding unnecessary AST traversal when no more
// processing is required.
var StopTraversal = errors.New("stop traversal")

// WalkFunc is the type of the function called for each node visited
// by Walk.  If the function returns SkipChildren, Walk will not
// traverse the node’s children.
type WalkFunc func(node Node) error

// Walk traverses the provided AST nodes, calling fn for each node
// recursively.
func Walk(nodes []Node, fn WalkFunc) error {
	for _, node := range nodes {
		if err := walk(node, fn); err != nil {
			if err == StopTraversal {
				return nil
			}
			return err
		}
	}
	return nil
}

func walk(node Node, fn WalkFunc) error {
	if err := fn(node); err != nil {
		if err == SkipChildren {
			return nil
		}
		return err
	}

	for _, child := range node.Children {
		if err := walk(child, fn); err != nil {
			return err
		}
	}
	return nil
}
