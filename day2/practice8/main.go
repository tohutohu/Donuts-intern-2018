package main

import (
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

type N struct {
	Num int
}

func (n *N) Square() int {
	return n.Num * n.Num
}

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

	e.GET("/square/:num", func(c echo.Context) error {
		num, err := strconv.Atoi(c.Param("num"))
		if err != nil {
			return c.String(http.StatusBadRequest, "num is not a Number")
		}
		temp, err := pongo2.FromFile("./template/square.html")
		if err != nil {
			return err
		}
		n := &N{num}

		ctx := pongo2.Context{
			"N": n,
		}

		return temp.ExecuteWriter(ctx, c.Response())
	})

	e.Start(":1323")
}
