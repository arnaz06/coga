package coga_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/arnaz06/coga"
)

type eventHandler interface {
	Handle(coga.SystemEvent)
}

type eventHandlerFunc func(coga.SystemEvent)

func (f eventHandlerFunc) Handle(e coga.SystemEvent) {
	f(e)
}

type bus struct {
	mu       sync.RWMutex
	handlers []eventHandler
}

func (b *bus) Publish(ctx context.Context, e coga.SystemEvent) {
	b.mu.RLock()
	for _, h := range b.handlers {
		h.Handle(e)
	}
	b.mu.RUnlock()
}

func (b *bus) Subscribe(h eventHandler) {
	b.mu.Lock()
	b.handlers = append(b.handlers, h)
	b.mu.Unlock()
}

func ExamplePublishSystemEvent_withoutPublisher() {
	coga.PublishSystemEvent(context.Background(), coga.Message{
		ID:      "123",
		Service: "order-service",
		Event:   coga.EventStart,
		Data:    []byte(`{"message":"lorem ipsum"}`),
	})

	// Output:
}

func ExamplePublishSystemEvent_withEbus() {
	ebus := new(bus)
	ebus.Subscribe(eventHandlerFunc(func(e coga.SystemEvent) {
		itemJSON, err := json.Marshal(e.Body)
		if err != nil {
			panic(err)
		}

		fmt.Printf("name: '%s'\n", e.Name)
		fmt.Printf("item: %s\n", itemJSON)
	}))

	ctx := context.WithValue(context.Background(), coga.ContextKeyPublisher, ebus)

	coga.PublishSystemEvent(ctx, coga.Message{
		ID:      "123",
		Service: "order-service",
		Event:   coga.EventStart,
		Data:    []byte(`{"message":"lorem ipsum"}`),
	})

	//Output:
	// name: 'Message'
	// item: {"id":"123","service":"order-service","event":"start","data":"eyJtZXNzYWdlIjoibG9yZW0gaXBzdW0ifQ=="}
}
