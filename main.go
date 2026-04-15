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
	flags, rest, err := opts.GetLong(os.Args, []opts.LongOpt{
		{Short: 'c', Long: "keep-comments", Arg: opts.None},
		{Short: 'C', Long: "clear-path", Arg: opts.None},
		{Short: 'd', Long: "no-doctype", Arg: opts.None},
		{Short: 'h', Long: "help", Arg: opts.None},
		{Short: 'I', Long: "include", Arg: opts.Required},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		fmt.Fprintf(os.Stderr,
			"Usage: %s [-cCd] [-I dirname] [file ...]\n"+
				"       %s -h\n",
			os.Args[0], os.Args[0])
		os.Exit(1)
	}

	fopts := formatter.Options{
		Doctype: true,
		SearchPath: []string{"./macros"},
	}

	for _, f := range flags {
		switch f.Key {
		case 'c':
			fopts.Comments = true
		case 'd':
			fopts.Doctype = false
		case 'I':
			fopts.SearchPath = append(fopts.SearchPath, f.Value)
		case 'h':
			openManual()
			os.Exit(0)
		}
	}

	if len(rest) == 0 {
		process("-", fopts)
	}

	for _, a := range rest {
		process(a, fopts)
	}
}

func process(filename string, fopts formatter.Options) {
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
