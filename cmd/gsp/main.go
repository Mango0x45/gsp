package main

import (
	"fmt"
	"os"
	"os/exec"

	"git.sr.ht/~mango/opts/v2"
	"git.thomasvoss.com/gsp/formatter"
	"git.thomasvoss.com/gsp/parser"
)

func main() {
	flags, rest, err := opts.Get(os.Args, "cdhI:")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		fmt.Fprintf(os.Stderr,
			"Usage: %s [-cd] [-I dirname] [file ...]\n"+
				"       %s -h\n",
			os.Args[0], os.Args[0])
		os.Exit(1)
	}

	fopts := formatter.Options{Doctype: true}

	for _, f := range flags {
		switch f.Key {
		case 'c':
			fopts.Comments = true
		case 'd':
			fopts.Doctype = false
		case 'h':
			openManual()
			os.Exit(0)
		case 'I':
			fopts.SearchPath = append(fopts.SearchPath, f.Value)
		}
	}

	if len(rest) == 0 {
		process("-", fopts)
	}

	for _, a := range rest {
		process(a, fopts)
	}
}

func process(path string, fopts formatter.Options) {
	var (
		file *os.File
		err  error
	)

	if path == "-" {
		file = os.Stdin
	} else {
		if file, err = os.Open(path); err != nil {
			die(err)
		} else {
			defer file.Close()
		}
	}

	ast, err := parser.Parse(file)
	if err != nil {
		die(err)
	}

	if err = formatter.WriteAst(os.Stdout, ast, fopts); err != nil {
		die(err)
	}
	fmt.Print("\n")
}

func openManual() {
	cmd := exec.Command("man", "1", "gsp")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		die(err)
	}
}

func die(e error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], e)
	os.Exit(1)
}
