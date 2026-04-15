package main

import (
	"fmt"
	"os"

	"git.sr.ht/~mango/opts/v2"
	"git.thomasvoss.com/gsp/formatter"
	"git.thomasvoss.com/gsp/parser"
)

var cflag, dflag bool

func main() {
	flags, rest, err := opts.GetLong(os.Args, []opts.LongOpt{
		{Short: 'c', Long: "keep-comments", Arg: opts.None},
		{Short: 'd', Long: "no-doctype", Arg: opts.None},
		{Short: 'h', Long: "help", Arg: opts.None},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		fmt.Fprintf(os.Stderr,
			"Usage: %s [-cd] [file ...]\n"+
				"       %s -h\n",
			os.Args[0], os.Args[0])
		os.Exit(1)
	}

	for _, f := range flags {
		switch f.Key {
		case 'c':
			cflag = true
		case 'd':
			dflag = true
		case 'h':
			panic("TODO")
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

	ast, err := parser.Parse(file)
	if err != nil {
		die(err)
	}

	if !dflag {
		fmt.Print("<!DOCTYPE html>")
	}
	if err = formatter.WriteAst(os.Stdout, ast); err != nil {
		die(err)
	}
	fmt.Print("\n")
}

func die(e error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], e)
	os.Exit(1)
}
