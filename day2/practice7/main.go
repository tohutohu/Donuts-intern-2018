package main

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/hello/:name", func(c echo.Context) error {
		name := c.Param("name")
		temp, err := pongo2.FromFile("./template/hello.html")
		if err != nil {
			return err
		}

		return temp.ExecuteWriter(pongo2.Context{"name": name}, c.Response())
	})

	e.Start(":1323")
}
