package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ably/ably-go/ably"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"net/http"
)

func main() {
	client, err := ably.NewRealtime(ably.WithKey("w5_O4Q.EJkpXw:wKz-OjK5KMpImH5IZBA3Vq0kxH1niQMN7LsdaKKCKDU"))
	if err != nil {
		panic(err)
	}
	channel := client.Channels.Get("chat")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/message", func(w http.ResponseWriter, r *http.Request) {
		var message *ably.Message
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		got := Publishing(channel, message)
		if got == false {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/message", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	Subscribe(channel)
	http.ListenAndServe(":8080", r)

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

func Subscribe(channel *ably.RealtimeChannel) string {
	var message string
	_, err := channel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		message = fmt.Sprintf("Received message from %v: '%v'\n", msg.ClientID, msg.Data)
		fmt.Println(message)
	})
	if err != nil {
		err := fmt.Errorf("subscribing to a channel: %w", err)
		fmt.Println(err)

	}
	fmt.Println(message)
	return message
}

func getHistory(channel *ably.RealtimeChannel) {
	// Before subscribing for messages, check the channel's
	// History for any missed messages. By default a channel
	// will keep 2 minutes of history available, but this can
	// be extended to 48 hours
	pages, err := channel.History().Pages(context.Background())
	if err != nil || pages == nil {
		return
	}

	hasHistory := true

	for ; hasHistory; hasHistory = pages.Next(context.Background()) {
		for _, msg := range pages.Items() {
			fmt.Printf("Previous message from %v: '%v'\n", msg.ClientID, msg.Data)
		}
	}
}
