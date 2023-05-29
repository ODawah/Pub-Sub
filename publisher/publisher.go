package main

import (
	"Pub-Sub/broker"
	"encoding/json"
	"github.com/ably/ably-go/ably"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	err := broker.AblyInit()
	if err != nil {
		panic(err)
	}
	channelName := "chat"
	channel := broker.Client.Channels.Get(channelName)

	r := Routes(channel)
	http.ListenAndServe(":8080", r)
}

func Routes(channel *ably.RealtimeChannel) http.Handler {
	r := chi.NewRouter()
	r.Post("/api/publish", func(w http.ResponseWriter, r *http.Request) {
		var message *ably.Message
		if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		got := broker.Publishing(channel, message)
		if got == false {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	return r
}
