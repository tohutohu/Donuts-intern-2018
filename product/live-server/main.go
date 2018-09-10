package main

import (
	"fmt"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/tohutohu/Donuts/day6/practice1/live-server/db"
	"github.com/tohutohu/Donuts/day6/practice1/live-server/router"
)

func main() {
	db, err := db.New()
	if err != nil {
		panic(err)
	}

	r := router.New(db)

	e := echo.New()
	//e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	//	AllowOrigins:     []string{"http://localhost:1234"},
	//	AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	//	AllowCredentials: true,
	//}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println(c.Request().RequestURI)
			for k, v := range c.QueryParams() {
				fmt.Printf("%s : %v\n", k, v)
			}
			fmt.Println()
			return next(c)
		}
	})

	e.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/rtmp/on_publish", r.StartLive)
	e.GET("/rtmp/on_publish_done", r.EndLive)

	e.GET("/api/lives", r.GetLives)
	e.POST("/api/lives", r.GetLiveEndpoint)
	// e.GET("/api/lives", getLives)

	e.Logger.Fatal(e.Start(":1323"))
}
