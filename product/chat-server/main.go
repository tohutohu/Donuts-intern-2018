// Package main provides ...
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/srinathgs/mysqlstore"
)

var (
	upgrader websocket.Upgrader
)

type H map[string]string

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	s *Server
)

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

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname))

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB(), "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	s = NewServer()
	go s.run()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ws/api/ranking", func(c echo.Context) error {
		res := s.redis.client.ZRangeWithScores("live-ranking", 0, 1000)
		return c.JSON(http.StatusOK, res.Val())
	})
	e.GET("/ws/:room", joinChatroom)
	e.Start(":1323")
}

func joinChatroom(c echo.Context) error {
	roomName := c.Param("room")

	var username string
	sess, _ := session.Get("sessions", c)
	if sess.Values["user_id"] != nil {
		id := sess.Values["user_id"].(uint)
		username = sess.Values["username"].(string)
		fmt.Println(id, username)
	}
	if username == "" {
		username = "名無しさん"
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	hub := s.getHub(roomName)

	client := &Client{
		conn: ws,
		hub:  hub,
		send: make(chan WSMessage, 256),
		name: username,
	}
	hub.register <- client
	go client.readPump()
	client.writePump()
	s.CleanRooms()
	return nil
}
