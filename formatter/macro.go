package formatter

import (
	"fmt"
	"io"
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

func execMacro(out io.Writer, path string, node ast.Node, opts Options) error {
	verbatim := node.Type == ast.VerbatimMacro

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
		return err
	}

	var stdout io.ReadCloser
	if verbatim {
		cmd.Stdout = out
	} else {
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			return err
		}
	}

	if err = cmd.Start(); err != nil {
		return err
	}
	if err = WriteUntranslatedAST(stdin, node.Children); err != nil {
		return err
	}
	stdin.Close()

	if !verbatim {
		nodes, err := parser.Parse(stdout)
		if err != nil {
			return err
		}
		if err = writeNodes(out, nodes, opts); err != nil {
			return err
		}
	}
	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
