package main

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		file, err := ioutil.ReadFile("./index.html")
		if err != nil {
			return err
		}
		return c.HTMLBlob(http.StatusOK, file)
	})

	e.Start(":1323")
}
