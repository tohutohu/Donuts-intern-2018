package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	r := &Redis{}
	r.client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	return r
}

func (r *Redis) run(roomEventCh chan *roomEvent) {
	pubsub := r.client.Subscribe("room_events")
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				break
			}
			fmt.Println(msg)
			data := &roomEvent{}
			json.Unmarshal([]byte(msg.Payload), data)
			roomEventCh <- data
		}
	}
}
