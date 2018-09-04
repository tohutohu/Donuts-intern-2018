package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 0 {
		return
	}
	for i, v := range os.Args[1:] {
		fmt.Printf("%s", v)
		if i < len(os.Args[1:])-1 {
			fmt.Print(" ")
		}
	}
}
