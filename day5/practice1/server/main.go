package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type Server struct {
	hubs map[string]*Hub
}

func NewServer() *Server {
	s := &Server{}
	s.hubs = make(map[string]*Hub)
	return s
}

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	s := NewServer()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		go func(con *websocket.Conn) {
			con.WriteMessage(websocket.TextMessage, []byte("which rooms you join ?"))
			con.WriteMessage(websocket.TextMessage, []byte(strings.Join(s.getRoomList(), "\n")))

			_, message, err := con.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(message))
			roomName := string(message)
			var hub *Hub
			if h, ok := s.hubs[roomName]; !ok {
				con.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s is not found. make room %s and join", roomName, roomName)))
				hub = newHub()
				go hub.run()
				s.hubs[roomName] = hub
			} else {
				hub = h
			}

			client := &Client{
				conn: con,
				hub:  hub,
				send: make(chan []byte, 256),
			}
			hub.register <- client
			go client.readPump()
			go client.writePump()
		}(conn)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *Server) getRoomList() []string {
	res := []string{}
	for k, _ := range s.hubs {
		res = append(res, k)
	}
	return res
}
