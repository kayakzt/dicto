package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	return runWithDeps(args, stdout, stderr, getWeblioDict, getWeblioEijiDict)
}

type dictFunc func(word string, out io.Writer) error

func runWithDeps(args []string, stdout, stderr io.Writer, dict, eijiDict dictFunc) error {
	fs := flag.NewFlagSet("dicto", flag.ContinueOnError)
	fs.SetOutput(stderr)

	eiji := fs.Bool("eiji", false, "Search Japanese-English dictionary on Weblio.")
	fs.BoolVar(eiji, "e", false, "Alias for --eiji")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "Usage: dicto [--eiji|-e] <word>")
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return errors.New("invalid arguments")
	}

	word := fs.Arg(0)
	if *eiji {
		return eijiDict(word, stdout)
	}
	return dict(word, stdout)
}
