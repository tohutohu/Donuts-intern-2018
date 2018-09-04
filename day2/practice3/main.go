package main

import (
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/:param", func(c echo.Context) error {
		param := c.Param("param")
		return c.HTML(http.StatusOK, heredoc.Docf(`
			<h1>%s</h1>
		`, param))
	})

	e.Start(":1323")
}
