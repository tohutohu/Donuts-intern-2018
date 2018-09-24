package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"
)

type M struct {
	Message string `json:"message"`
}

func main() {
	username := os.Getenv("MYSQL_USERNAME")
	if username == "" {
		username = "root"
	}

	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		password = "mysql"
	}

	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}

	dbname := os.Getenv("MYSQL_DB_NAME")
	if dbname == "" {
		dbname = "db"
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname))

	if err != nil {
		panic(err)
	}

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB(), "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	e.GET("/ping", func(c echo.Context) error {
		sess, _ := session.Get("sessions", c)
		sess.Values["user_id"] = uint(12)
		sess.Save(c.Request(), c.Response())

		return c.String(http.StatusOK, "pong")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func WithLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, M{"something wrong in getting session"})
		}

		if sess.Values["user_id"] == nil {
			return c.JSON(http.StatusForbidden, M{"please login"})
		}
		c.Set("userID", int(sess.Values["user_id"].(uint)))

		return next(c)
	}
}

