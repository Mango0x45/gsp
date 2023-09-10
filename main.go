package main

import (
	"fmt"
	"os"

	"git.thomasvoss.com/getgopt"
	"git.thomasvoss.com/gsp/formatter"
	"git.thomasvoss.com/gsp/parser"
)

var dflag bool

func main() {
	for opt := byte(0); getgopt.Getopt(len(os.Args), os.Args, "d", &opt); {
		switch opt {
		case 'd':
			dflag = true
		default:
			fmt.Fprintf(os.Stderr, "Usage: %s [-d] [file ...]\n", os.Args[0])
			os.Exit(1)
		}
	}

	args := os.Args[getgopt.Optind:]

	if len(args) == 0 {
		process("-")
	}

	for _, a := range args {
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
