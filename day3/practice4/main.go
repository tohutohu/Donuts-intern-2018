package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/go/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.Start(":1323")
}
