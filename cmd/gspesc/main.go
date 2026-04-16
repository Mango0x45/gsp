package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"git.sr.ht/~mango/opts/v2"
)

var (
	attrchars = [256]bool{
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, true, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, true, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
	}
	descchars = [256]bool{
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		true, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, true, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, true, false, true, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
		false, false, false, false, false, false, false, false,
	}
	chars = descchars
)

func main() {
	flags, rest, err := opts.Get(os.Args, "ah")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		fmt.Fprintf(os.Stderr,
			"Usage: %s [-a] [file ...]\n"+
				"       %s -h\n",
			os.Args[0], os.Args[0])
		os.Exit(1)
	}

	for _, f := range flags {
		switch f.Key {
		case 'a':
			chars = attrchars
		case 'h':
			openManual()
			os.Exit(0)
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
	var (
		file *os.File
		err  error
	)

	if filename == "-" {
		file = os.Stdin
	} else {
		if file, err = os.Open(filename); err != nil {
			die(err)
		} else {
			defer file.Close()
		}
	}

	for {
		var buf [1]byte
		_, err = file.Read(buf[:])
		switch err {
		case io.EOF:
			return
		case nil:
			ch := buf[0]
			if chars[ch] {
				os.Stdout.Write([]byte{'\\', ch})
			} else {
				os.Stdout.Write([]byte{ch})
			}
		default:
			die(err)
		}
	}
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
