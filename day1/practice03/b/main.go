package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
)

func main() {
	var requestUrl string
	flag.StringVar(&requestUrl, "url", "", "set request url")
	flag.Parse()

	if valid := govalidator.IsRequestURL(requestUrl); !valid {
		fmt.Printf("%v is a invalid url\n", requestUrl)
		os.Exit(1)
		return
	}

	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Printf("error occured in request to %s\n%s\n", requestUrl, err.Error())
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
