package main

import (
	"Pub-Sub/broker"
	"Pub-Sub/persistence"
	"fmt"
	"github.com/ably/ably-go/ably"
)

func main() {
	err := persistence.InitMongo("<MONGO-URI>")
	if err != nil {
		panic(err)
	}
	err = broker.AblyInit()
	if err != nil {
		panic(err)
	}
	ch := make(chan *ably.Message)
	broker.Subscribe(broker.Client.Channels.Get("chat"), ch)
	for {
		select {
		case <-ch:
			ev := <-ch
			fmt.Println(ev)
		}
	}
}
