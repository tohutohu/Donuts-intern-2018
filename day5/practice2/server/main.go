// Package main provides ...
package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader websocket.Upgrader
)

type Server struct {
	hubs      map[string]*Hub
	hubsMutex *sync.Mutex
}

func NewServer() *Server {
	s := &Server{}
	s.hubs = make(map[string]*Hub)
	s.hubsMutex = &sync.Mutex{}
	return s
}

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
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	e.GET("/ws/:room", joinChatroom)
	e.Start(":1323")
}

func joinChatroom(c echo.Context) error {
	roomName := c.Param("room")
	fmt.Println(roomName)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	var hub *Hub
	s.hubsMutex.Lock()
	if h, ok := s.hubs[roomName]; !ok {
		ws.WriteJSON(Message{Type: "create_room", Data: H{"message": roomName + " is created"}})
		hub = newHub()
		go hub.run()
		s.hubs[roomName] = hub
	} else {
		hub = h
	}
	s.hubsMutex.Unlock()

	client := &Client{
		conn: ws,
		hub:  hub,
		send: make(chan interface{}, 256),
	}
	hub.register <- client
	go client.readPump()
	client.writePump()
	return nil
}

func (s *Server) getRoomList() []string {
	res := []string{}
	for k, _ := range s.hubs {
		res = append(res, k)
	}
	return res
}
