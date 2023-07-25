package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/quevivasbien/bird-game/game"
	"github.com/valyala/fasthttp"
)

type CloseCode int

const (
	ContinueCode CloseCode = iota
	EmptyCode
)

type Subscription[T game.HasID] struct {
	data  chan T
	close chan CloseCode
}

type Manager[T game.HasID] struct {
	items map[string]T
	subs  map[string](map[string]Subscription[T])
}

func (m Manager[T]) Get(id string) (T, bool) {
	item, exists := m.items[id]
	return item, exists
}

func (m Manager[T]) Put(item T) {
	id := item.GetID()
	m.items[id] = item
	subs, exists := m.subs[id]
	if !exists {
		m.subs[id] = make(map[string]Subscription[T])
		return
	}
	for _, s := range subs {
		s.data <- item
	}
}

func (m Manager[T]) Delete(id string, code CloseCode) {
	delete(m.items, id)
	subs, exists := m.subs[id]
	if !exists {
		return
	}
	for _, s := range subs {
		s.close <- code
	}
	delete(m.subs, id)
}

func (m Manager[T]) Subscribe(id string, subscriber string, c *fiber.Ctx) (Subscription[T], error) {
	_, exists := m.subs[id]
	if !exists {
		return Subscription[T]{}, fmt.Errorf("Attempted to subscribe to an item, %s, that doesn't exist", id)
	}
	sub := Subscription[T]{make(chan T), make(chan CloseCode)}
	m.subs[id][subscriber] = sub

	if c == nil {
		return sub, nil
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			select {
			case item := <-sub.data:
				data, err := json.Marshal(item)
				if err != nil {
					log.Println("Got error when processing stream notification:", err)
					break
				}
				msg := fmt.Sprintf("event: update\ndata: %s\n\n", data)
				log.Printf("Sending message:\n%v", msg)
				fmt.Fprintf(w, msg)
			case code := <-sub.close:
				if code == ContinueCode {
					log.Printf("Notifying of continue signal")
					fmt.Fprintf(w, "event: continue\ndata: %d\n\n", code)
				} else {
					log.Printf("Notifying of deletion signal; code = %v", code)
					fmt.Fprintf(w, "event: delete\ndata: %d\n\n", code)
				}
				return
			}
			err := w.Flush()
			if err != nil {
				log.Printf("Error while flushing: %v. Closing stream.", err)
				return
			}
		}
	}))

	return sub, nil
}

func (m Manager[T]) Unsubscribe(id string, subscriber string) {
	_, exists := m.subs[id]
	if !exists {
		return
	}
	delete(m.subs[id], subscriber)
}

func MakeManager[T game.HasID]() Manager[T] {
	return Manager[T]{
		items: make(map[string]T),
		subs:  make(map[string]map[string]Subscription[T]),
	}
}
