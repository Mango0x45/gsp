package main

import (
	"fmt"
	"os"

	"git.thomasvoss.com/gsp/formatter"
	"git.thomasvoss.com/gsp/parser"
)

func main() {
	if len(os.Args) == 1 {
		process("-")
	}

	for _, arg := range os.Args[1:] {
		process(arg)
	}
}

func process(filename string) {
	var file *os.File
	var err error

	if filename == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(filename)
		if err != nil {
			die(err)
		}
		defer file.Close()
	}

	ast, err := parser.ParseFile(file)
	if err != nil {
		die(err)
	}

	formatter.PrintAst(ast)
	fmt.Print("\n")
}

func die(e error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], e)
	os.Exit(1)
}
