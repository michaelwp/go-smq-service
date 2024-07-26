package main

import (
	"fmt"
	"net/http"
)

var broker = NewMessageBroker()

func publishHandler(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	message := r.URL.Query().Get("message")

	if topic == "" || message == "" {
		http.Error(w, "Missing topic or message", http.StatusBadRequest)
		return
	}

	broker.Publish(topic, message)
	fmt.Fprintf(w, "Published message to topic: %s\n", topic)
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")

	if topic == "" {
		http.Error(w, "Missing topic", http.StatusBadRequest)
		return
	}

	ch := broker.Subscribe(topic)
	defer broker.Unsubscribe(topic, ch)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
