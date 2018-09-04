package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	flag.Parse()

	url := flag.Arg(0)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error occured in request to %s\n%s\n", url, err.Error())
		os.Exit(1)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error occured in reading body\n%s\n", err.Error())
		os.Exit(1)
		return
	}
	fmt.Println(string(body))
}
