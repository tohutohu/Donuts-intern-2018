package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

var (
	db *gorm.DB
)

type Live struct {
	gorm.Model
	Name string
	E    string
	St   string
	Done bool
}

func main() {
	username := os.Getenv("MYSQL_USERNAME")
	if username == "" {
		username = "root"
	}

	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		password = "password"
	}

	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "db"
	}

	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}

	dbname := os.Getenv("MYSQL_DB_NAME")
	if dbname == "" {
		dbname = "live"
	}
	for {
		_db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname))

		if err == nil {
			db = _db
			break
		}
		fmt.Println(err)
		fmt.Println("DB Start up waiting...")
		time.Sleep(2 * time.Second)
	}
	db.AutoMigrate(&Live{})
	db.DropTableIfExists(&Live{})

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

	e.GET("/rtmp/on_publish", startLive)
	e.GET("/rtmp/on_publish_done", endLive)

	e.GET("/api/lives", getLives)

	e.GET("/api/lives", getLiveEndpoint)
	// e.GET("/api/lives", getLives)

	e.Logger.Fatal(e.Start(":1323"))
}

func getLiveEndpoint(c echo.Context) error {
	rp := strings.NewReplacer("+", "-", "/", "_", "=", "")
	unix := time.Now().Unix()
	h := md5.New()
	req := &struct {
		Name string `json:"name"`
	}{}
	if req.Name == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	c.Bind(req)
	h.Write([]byte("CocoroIsGodlive/" + req.Name + strconv.FormatInt(unix, 10)))
	fmt.Println(unix)
	url := rp.Replace(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	return c.JSON(http.StatusOK, map[string]string{"url": url})
}

func startLive(c echo.Context) error {
	name := c.QueryParam("name")
	e := c.QueryParam("e")
	st := c.QueryParam("st")

	live := &Live{
		Name: name,
		E:    e,
		St:   st,
		Done: false,
	}

	db.Create(live)

	return c.String(http.StatusOK, "OK")
}

func endLive(c echo.Context) error {
	name := c.QueryParam("name")
	e := c.QueryParam("e")
	st := c.QueryParam("st")

	live := &Live{
		Name: name,
		E:    e,
		St:   st,
	}

	db.First(live)
	if live.ID == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	live.Done = true
	db.Save(live)

	return c.String(http.StatusOK, "OK")
}

func getLives(c echo.Context) error {
	lives := []Live{}
	db.Find(&lives, "done = false")
	return c.JSON(http.StatusOK, lives)
}
