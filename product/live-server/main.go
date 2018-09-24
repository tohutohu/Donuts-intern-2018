package main

import (
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"
	"github.com/tohutohu/Donuts/product/live-server/db"
	"github.com/tohutohu/Donuts/product/live-server/router"
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

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB(), "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.GET("/rtmp/on_publish", r.StartLive)
	e.GET("/rtmp/on_publish_done", r.EndLive)

	e.GET("/api/lives", r.GetLives)
	e.POST("/api/lives", r.GetLiveEndpoint)
	e.POST("/api/users", r.PostUsers)
	e.POST("/api/login", r.PostLogin)
	e.POST("/api/whoami", r.GetWhoAmI)
	// e.GET("/api/lives", getLives)

	e.Logger.Fatal(e.Start(":1323"))
}
