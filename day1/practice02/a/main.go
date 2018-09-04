package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	var noLineBreak bool
	var split string
	flag.BoolVar(&noLineBreak, "n", false, "Not line break")
	flag.StringVar(&split, "s", " ", "Change Split Strings")
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), split))
	if !noLineBreak {
		fmt.Print("\n")
	}
}
