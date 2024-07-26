package main

import (
	"sync"
)

// MessageBroker struct
type MessageBroker struct {
	subscribers map[string][]chan string
	sync.Mutex
}

// NewMessageBroker initializes a new MessageBroker
func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		subscribers: make(map[string][]chan string),
	}
}

// Subscribe adds a subscriber to a specific topic
func (mb *MessageBroker) Subscribe(topic string) <-chan string {
	mb.Lock()
	defer mb.Unlock()

	ch := make(chan string)
	mb.subscribers[topic] = append(mb.subscribers[topic], ch)
	return ch
}

// Publish sends a message to all subscribers of a topic
func (mb *MessageBroker) Publish(topic, msg string) {
	mb.Lock()
	defer mb.Unlock()

	if subscribers, found := mb.subscribers[topic]; found {
		for _, ch := range subscribers {
			go func(c chan string) {
				c <- msg
			}(ch)
		}
	}
}

// Unsubscribe removes a subscriber from a specific topic
func (mb *MessageBroker) Unsubscribe(topic string, sub <-chan string) {
	mb.Lock()
	defer mb.Unlock()

	if subscribers, found := mb.subscribers[topic]; found {
		for i, ch := range subscribers {
			if ch == sub {
				mb.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
				close(ch)
				break
			}
		}
	}
}
