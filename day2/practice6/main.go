package main

import (
	"html/template"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/:fileName", func(c echo.Context) error {
		fileName := c.Param("fileName")
		temp, err := template.ParseFiles("./template/header.html", "./template/"+fileName, "./template/footer.html")
		if err != nil {
			return err
		}
		return temp.ExecuteTemplate(c.Response(), "content", nil)
	})

	e.Static("/static", "./static")

	e.Start(":1323")
}
