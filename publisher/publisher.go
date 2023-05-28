package main

import (
	"Pub-Sub/broker"
	"context"
	"fmt"
	"time"
)

func main() {
	err := broker.AblyInit()
	if err != nil {
		panic(err)
	}
	channel := "chat"
	publisher := broker.Client.Channels.Get(channel)

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		err := publisher.Publish(context.Background(), "message", "Hello World!")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Sent message %d\n", i)
	}
}
