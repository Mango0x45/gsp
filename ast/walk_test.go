package ast

import (
	"errors"
	"testing"
)

func TestWalk_FullTraversal(t *testing.T) {
	nodes := []Node{
		{
			Type: Normal,
			Name: "html",
			Children: []Node{
				{
					Type: Normal,
					Name: "head",
				},
				{
					Type: Normal,
					Name: "body",
					Children: []Node{
						{Type: Text, Name: "Hello"},
					},
				},
			},
		},
	}

	var visited []string
	err := Walk(nodes, func(node *Node) error {
		visited = append(visited, node.Name)
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"html", "head", "body", "Hello"}
	if len(visited) != len(expected) {
		t.Fatalf("expected %d visited nodes, got %d",
			len(expected), len(visited))
	}

	for i, name := range expected {
		if visited[i] != name {
			t.Errorf("expected node %d to be %q, got %q", i, name, visited[i])
		}
	}
}

func TestWalk_SkipChildren(t *testing.T) {
	nodes := []Node{
		{
			Type: Normal,
			Name: "html",
			Children: []Node{
				{
					Type: Comment,
					Name: "/",
					Children: []Node{
						{Type: Normal, Name: "hidden"},
					},
				},
				{
					Type: Normal,
					Name: "body",
				},
			},
		},
	}

	var visited []string
	err := Walk(nodes, func(node *Node) error {
		visited = append(visited, node.Name)
		if node.Type == Comment {
			return SkipChildren
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"html", "/", "body"}
	if len(visited) != len(expected) {
		t.Fatalf("expected %d visited nodes, got %d",
			len(expected), len(visited))
	}

	for i, name := range expected {
		if visited[i] != name {
			t.Errorf("expected node %d to be %q, got %q", i, name, visited[i])
		}
	}
}

func TestWalk_StopTraversal(t *testing.T) {
	nodes := []Node{
		{
			Type: Normal,
			Name: "html",
			Children: []Node{
				{Type: Normal, Name: "head"},
				{
					Type: Normal,
					Name: "body",
					Children: []Node{
						{Type: Text, Name: "Hello"},
					},
				},
			},
		},
		{
			Type: Normal,
			Name: "footer",
		},
	}

	var visited []string
	err := Walk(nodes, func(node *Node) error {
		visited = append(visited, node.Name)
		if node.Name == "head" {
			return StopTraversal
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"html", "head"}
	if len(visited) != len(expected) {
		t.Fatalf("expected %d visited nodes, got %d",
			len(expected), len(visited))
	}

	for i, name := range expected {
		if visited[i] != name {
			t.Errorf("expected node %d to be %q, got %q", i, name, visited[i])
		}
	}
}

func TestWalk_ErrorPropagation(t *testing.T) {
	nodes := []Node{
		{
			Type: Normal,
			Name: "html",
			Children: []Node{
				{Type: Normal, Name: "head"},
				{Type: Normal, Name: "body"},
			},
		},
	}

	expectedErr := errors.New("halt execution")
	var visited []string

	err := Walk(nodes, func(node *Node) error {
		visited = append(visited, node.Name)
		if node.Name == "head" {
			return expectedErr
		}
		return nil
	})

	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	expected := []string{"html", "head"}
	if len(visited) != len(expected) {
		t.Fatalf("expected %d visited nodes, got %d",
			len(expected), len(visited))
	}

	for i, name := range expected {
		if visited[i] != name {
			t.Errorf("expected node %d to be %q, got %q", i, name, visited[i])
		}
	}
}
