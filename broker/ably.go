package broker

import (
	"context"
	"fmt"
	"github.com/ably/ably-go/ably"
)

var Client *ably.Realtime

func AblyInit() error {
	client, err := ably.NewRealtime(ably.WithKey("w5_O4Q.U7ulRg:h_hrf0BxtKFnBXB1-NyJKyLlBe_Q9T9nAVQC8QWk7us"))
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
		err := channel.Publish(context.Background(), "message", message.Data)
		// await confirmation that message was received by Ably
		if err != nil {
			err := fmt.Errorf("publishing to channel: %w", err)
			fmt.Println(err)
			return false
		}
		return true
	}
}

func Subscribe(channel *ably.RealtimeChannel, local chan string) string {
	var message string
	_, err := channel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		message = fmt.Sprintf("Received message from %v: '%v'\n", msg.ClientID, msg.Data)
		fmt.Println(message)
		local <- message
	})
	if err != nil {
		err := fmt.Errorf("subscribing to a channel: %w", err)
		fmt.Println(err)

	}
	return message
}
