package formatter

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"git.thomasvoss.com/gsp/ast"
	"git.thomasvoss.com/gsp/parser"
)

func findMacro(name string, dirs []string) (string, bool) {
	for _, dir := range dirs {
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)

		if err == nil && !info.IsDir() {
			return path, true
		}
	}

	return "", false
}

func execMacro(path string, node ast.Node) ([]ast.Node, error) {
	env := os.Environ()
	for k, v := range maps.All(node.Attributes) {
		env = append(env, fmt.Sprintf("GSP_%s=%s",
			strings.ToUpper(k),
			strings.Join(v, " ")))
	}

	cmd := exec.Cmd{
		Path:   path,
		Env:    env,
		Stderr: os.Stderr,
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return []ast.Node{}, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return []ast.Node{}, err
	}
	if err = cmd.Start(); err != nil {
		return []ast.Node{}, err
	}
	if err = WriteUntranslatedAST(stdin, node.Children); err != nil {
		return []ast.Node{}, err
	}
	stdin.Close()
	nodes, err := parser.Parse(stdout)
	if err != nil {
		return []ast.Node{}, err
	}
	if err = cmd.Wait(); err != nil {
		return []ast.Node{}, err
	}

	return nodes, nil
}
