package main

import (
	"fmt"
	"os"
)

func main() {
	cur, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println(cur)
}
