package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println(c.Path())
			return next(c)
		}
	})
	e.GET("/v2/:path", func(c echo.Context) error {
		path := c.Param("path")
		return c.String(http.StatusOK, path)
	})

	e.Start(":1323")
}
