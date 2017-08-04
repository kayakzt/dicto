package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	eiji bool
)

func main() {
	// parse flag & args
	flag.BoolVar(&eiji, "eiji", false, "Search Japanese-English dictionary on Weblio.")
	flag.BoolVar(&eiji, "e", false, "Alias for --eiji")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Invalid arguments!")
		os.Exit(1)
	}

	word := flag.Args()[0]

	if eiji {
		getWeblioEijiDict(word)
	} else {
		getWeblioDict(word)
	}
}
