// Package main provides ...
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader websocket.Upgrader
)

type H map[string]string

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	s *Server
)

func main() {
	e := echo.New()
	s = NewServer()
	go s.run()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	e.GET("/ws/:room", joinChatroom)
	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	e.Start(":" + port)
}

func joinChatroom(c echo.Context) error {
	roomName := c.Param("room")
	fmt.Println(roomName)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	hub := s.getHub(roomName)

	client := &Client{
		conn: ws,
		hub:  hub,
		send: make(chan interface{}, 256),
	}
	hub.register <- client
	go client.readPump()
	client.writePump()
	s.CleanRooms()
	return nil
}
