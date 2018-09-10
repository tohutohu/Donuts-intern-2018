// Package main provides ...
package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Server struct {
	redis       *Redis
	hubs        map[string]*Hub
	hubsMutex   *sync.Mutex
	id          int64
	roomEventCh chan *roomEvent
}

type roomEvent struct {
	EventID  int    `json:"event_id"`
	Name     string `json:"name"`
	ServerID int64  `json:"server_id"`
}

const (
	CREATE = iota
	DELETE
)

func NewServer() *Server {
	s := &Server{}
	s.hubs = make(map[string]*Hub)
	s.hubsMutex = &sync.Mutex{}
	s.redis = NewRedis()
	s.roomEventCh = make(chan *roomEvent)
	return s
}

func (s *Server) run() {
	s.id = time.Now().UnixNano()
	go s.redis.run(s.roomEventCh)

	for {
		select {
		case event := <-s.roomEventCh:
			if event.ServerID == s.id {
				continue
			}
			switch event.EventID {
			case CREATE:
				s.getHub(event.Name)
			case DELETE:
				s.hubsMutex.Lock()
				delete(s.hubs, event.Name)
				s.hubsMutex.Unlock()
			}

		}
	}
}

func (s *Server) getRoomList() []string {
	res := []string{}
	for k, _ := range s.hubs {
		res = append(res, k)
	}
	return res
}

func (s *Server) CleanRooms() {
	for roomName, room := range s.hubs {
		if room.CountConnectedClients() == 0 {
			fmt.Println("clean", roomName)
			s.SendRoomEvent(DELETE, roomName)
			delete(s.hubs, roomName)
		}
	}
}

func (s *Server) getHub(roomName string) *Hub {
	s.hubsMutex.Lock()
	defer s.hubsMutex.Unlock()
	if h, ok := s.hubs[roomName]; !ok {
		// ws.WriteJSON(Message{Type: "create_room", Data: H{"message": roomName + " is created"}})
		s.SendRoomEvent(CREATE, roomName)
		hub := newHub(roomName)
		go hub.run(s.redis.client)
		s.hubs[roomName] = hub
		return hub
	} else {
		return h
	}
}

func (s *Server) SendRoomEvent(event int, roomName string) error {
	data := roomEvent{
		EventID:  event,
		Name:     roomName,
		ServerID: s.id,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.redis.client.Publish("room_events", string(body)).Err()
}
