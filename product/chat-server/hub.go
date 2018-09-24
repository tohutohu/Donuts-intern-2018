package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Hub struct {
	clients map[*Client]struct{}

	// Inbound messages from the clients.
	broadcast chan WSMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	roomName string
}

func (h *Hub) CountConnectedClients() int {
	return len(h.clients)
}

func newHub(roomName string) *Hub {
	return &Hub{
		broadcast:  make(chan WSMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
		roomName:   roomName,
	}
}

func (h *Hub) run(rc *redis.Client) {
	pubsub := rc.Subscribe(h.roomName)
	defer pubsub.Unsubscribe(h.roomName)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	poi := time.NewTicker(10 * time.Second)

	// Go channel which receives messages.
	ch := pubsub.Channel()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}

			go func(roomName string) {
				rc.Incr(roomName + "_cum_count")
				rc.Incr(roomName + "_now_count")
			}(h.roomName)

			go func(c *Client, roomName string) {
				cmd := rc.LRange(roomName, 0, 9)
				if cmd.Err() != nil {
					fmt.Println(cmd.Err())
					return
				}
				data := WSMessage{Type: "initial-data"}
				messages := []map[string]string{}
				for _, v := range cmd.Val() {
					m := map[string]string{}
					if err := json.Unmarshal([]byte(v), &m); err != nil {
						fmt.Println(err)
						continue
					}
					messages = append(messages, m)
				}
				data.Data = messages
				c.send <- data
			}(client, h.roomName)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)

				close(client.send)
				go func(roomName string) {
					rc.Decr(roomName + "_now_count")
				}(h.roomName)
			}

		case message := <-h.broadcast:
			body, _ := json.Marshal(message)
			fmt.Println("send ", string(body))
			rc.Publish(h.roomName, string(body))
			m, _ := json.Marshal(message.Data)
			rc.ZIncrBy("live-ranking", 1, h.roomName)
			rc.LPush(h.roomName, string(m))

		case message := <-ch:
			data := &WSMessage{}
			fmt.Println("recieve ", message)
			json.Unmarshal([]byte(message.Payload), data)
			for client := range h.clients {
				select {
				case client.send <- *data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		case <-poi.C:
			now, err := rc.Get(h.roomName + "_now_count").Result()
			if err != nil {
				fmt.Println(err)
				break
			}
			cum, err := rc.Get(h.roomName + "_cum_count").Result()
			if err != nil {
				fmt.Println(err)
				break
			}
			wsMes := WSMessage{
				Type: "room-info",
				Data: map[string]string{
					"now": now,
					"cum": cum,
				},
			}
			for client := range h.clients {
				select {
				case client.send <- wsMes:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}
}
