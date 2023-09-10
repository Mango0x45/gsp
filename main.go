package main

import (
	"fmt"
	"os"

	"git.thomasvoss.com/getgopt"
	"git.thomasvoss.com/gsp/formatter"
	"git.thomasvoss.com/gsp/parser"
)

func main() {
	for opt := byte(0); getgopt.Getopt(len(os.Args), os.Args, "x", &opt); {
		switch opt {
		case 'x':
			parser.Xml = true
		}
	}

	os.Args = os.Args[getgopt.Optind:]

	if len(os.Args) == 0 {
		process("-")
	}

	for _, arg := range os.Args {
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
