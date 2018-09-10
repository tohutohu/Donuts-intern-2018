// Package main provides ...
package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	client := NewClient()

	pubsub := client.Subscribe("mychannel1")

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := pubsub.Receive()
	if err != nil {
		panic(err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()

	// Publish a message.
	err = client.Publish("mychannel1", "hello").Err()
	if err != nil {
		panic(err)
	}

	time.AfterFunc(time.Second, func() {
		// When pubsub is closed channel is closed too.
		_ = pubsub.Close()
	})

	// Consume messages.
	for {
		msg, ok := <-ch
		if !ok {
			break
		}
		fmt.Println(msg.Channel, msg.Payload)
	}
}

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
