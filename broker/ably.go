package broker

import (
	"Pub-Sub/persistence"
	"context"
	"fmt"
	"github.com/ably/ably-go/ably"
	"go.mongodb.org/mongo-driver/bson"
)

var Client *ably.Realtime

func AblyInit() error {
	client, err := ably.NewRealtime(ably.WithKey("<ABLY-API-KEY>"))
	if err != nil {
		panic(err)
		return err
	}
	Client = client
	return nil
}

func Publishing(channel *ably.RealtimeChannel, message *ably.Message) bool {
	for {
		// Publish the message typed in to the Ably Channel
		err := channel.Publish(context.Background(), "message", message)
		// await confirmation that message was received by Ably
		if err != nil {
			err := fmt.Errorf("publishing to channel: %w", err)
			fmt.Println(err)
			return false
		}
		return true
	}
}

func Subscribe(channel *ably.RealtimeChannel, local chan *ably.Message) string {
	collection := persistence.MongoDB.Collection("farmers")
	var message string
	_, err := channel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		message = fmt.Sprintf("Received message from %v: '%v'\n", msg.ClientID, msg)
		local <- msg
		event := bson.M{"clientID": msg.ClientID, "message": msg}
		_, err := collection.InsertOne(context.Background(), event)
		if err != nil {
			panic(err)
		}

		fmt.Println(message)
	})
	if err != nil {
		err := fmt.Errorf("subscribing to a channel: %w", err)
		fmt.Println(err)

	}
	return message
}
