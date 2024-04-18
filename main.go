package main

import (
	"fmt"
	"os"

	"git.sr.ht/~mango/opts/v2"
	"git.thomasvoss.com/gsp/formatter"
	"git.thomasvoss.com/gsp/parser"
)

var dflag bool

func main() {
	flags, rest, err := opts.GetLong(os.Args, []opts.LongOpt{
		{Short: 'd', Long: "no-doctype", Arg: opts.None},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		fmt.Fprintf(os.Stderr, "Usage: %s [-d] [file ...]\n", os.Args[0])
		os.Exit(1)
	}

	for _, f := range flags {
		switch f.Key {
		case 'd':
			dflag = true
		}
	}

	if len(rest) == 0 {
		process("-")
	}

	for _, a := range rest {
		process(a)
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

	if !dflag {
		fmt.Print("<!DOCTYPE html>")
	}
	formatter.PrintAst(ast)
	fmt.Print("\n")
}

func die(e error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], e)
	os.Exit(1)
}
