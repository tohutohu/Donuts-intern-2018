package main

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, heredoc.Doc(`
			<h1>hellow world</h1>
		`))
	})

	e.Start(":1323")
}
