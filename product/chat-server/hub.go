package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type Hub struct {
	clients map[*Client]struct{}

	// Inbound messages from the clients.
	broadcast chan interface{}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	roomName string
}

type messageData struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

func (h *Hub) CountConnectedClients() int {
	return len(h.clients)
}

func newHub(roomName string) *Hub {
	return &Hub{
		broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
		roomName:   roomName,
	}
}

func (h *Hub) run(client *redis.Client) {
	pubsub := client.Subscribe(h.roomName)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			body, _ := json.Marshal(message)
			fmt.Println("send ", string(body))
			client.Publish(h.roomName, string(body))
		case message := <-ch:
			data := &messageData{}
			fmt.Println("recieve ", message)
			json.Unmarshal([]byte(message.Payload), data)
			for client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
