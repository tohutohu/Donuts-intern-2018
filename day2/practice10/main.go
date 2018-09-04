package main

import (
	"net/http"
	"strconv"

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

	e.GET("/square/:num", func(c echo.Context) error {
		num, err := strconv.Atoi(c.Param("num"))
		if err != nil {
			return c.String(http.StatusBadRequest, "num is not a Number")
		}
		temp, err := pongo2.FromFile("./template/square.html")
		if err != nil {
			return err
		}

		ctx := pongo2.Context{
			"num": num,
			"Square": func(n int) int {
				return n * n
			},
		}

		return temp.ExecuteWriter(ctx, c.Response())
	})

	e.GET("/sign/:num", func(c echo.Context) error {
		num, err := strconv.Atoi(c.Param("num"))
		if err != nil {
			return c.String(http.StatusBadRequest, "num is not a Number")
		}
		temp, err := pongo2.FromFile("./template/sign.html")
		if err != nil {
			return err
		}

		return temp.ExecuteWriter(pongo2.Context{"num": num}, c.Response())
	})

	e.GET("/fizzbuzz/:num", func(c echo.Context) error {
		num, err := strconv.Atoi(c.Param("num"))
		if err != nil {
			return c.String(http.StatusBadRequest, "num is not a Number")
		}
		temp, err := pongo2.FromFile("./template/fizzbuzz.html")
		if err != nil {
			return err
		}

		fizzbuzz := make([]string, num)
		for i, _ := range fizzbuzz {
			if (i+1)%15 == 0 {
				fizzbuzz[i] = "fizzbuzz"
			} else if (i+1)%5 == 0 {
				fizzbuzz[i] = "buzz"
			} else if (i+1)%3 == 0 {
				fizzbuzz[i] = "fizz"
			} else {
				fizzbuzz[i] = strconv.Itoa(i + 1)
			}
		}

		return temp.ExecuteWriter(pongo2.Context{"fizzbuzz": fizzbuzz}, c.Response())
	})

	e.Start(":1323")
}
