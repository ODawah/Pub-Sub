package main

import (
	"Pub-Sub/broker"
	"fmt"
)

func main() {
	broker.AblyInit()
	err := broker.AblyInit()
	if err != nil {
		panic(err)
	}
	ch := make(chan string)
	broker.Subscribe(broker.Client.Channels.Get("chat"), ch)
	for {
		select {
		case <-ch:
			fmt.Println(<-ch)
		}
	}
}
