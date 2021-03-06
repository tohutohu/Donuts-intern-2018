package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
)

var (
	methodsList = []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}
)

func main() {
	var requestUrl string
	var method string
	flag.StringVar(&requestUrl, "url", "", "set request url")
	flag.StringVar(&method, "method", "GET", "change request method (default GET)")
	flag.Parse()
	method = strings.ToUpper(method)

	if !govalidator.IsRequestURL(requestUrl) {
		fmt.Printf("%v is a invalid url\n", requestUrl)
		os.Exit(1)
		return
	}

	if !isValidMethod(method) {
		fmt.Printf("%s is a invalid method\n", method)
		os.Exit(1)
		return
	}

	req, err := http.NewRequest(method, requestUrl, nil)
	if err != nil {
		fmt.Printf("error occured in create request\n%s\n", err.Error())
		os.Exit(1)
		return
	}

	resp, err := http.DefaultClient.Do(req)
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

func isValidMethod(method string) bool {
	for _, m := range methodsList {
		if m == method {
			return true
		}
	}
	return false
}
