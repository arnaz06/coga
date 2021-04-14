package coga_test

import (
	"context"
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
	coga.PublishSystemEvent(context.Background(), "topic-key", coga.Message{
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

		fmt.Printf("name: '%s'\n", e.Name)
		fmt.Printf("body: %+v\n", e.Body)
	}))

	ctx := context.WithValue(context.Background(), "topic-key", ebus)

	coga.PublishSystemEvent(ctx, "topic-key", coga.Message{
		ID:      "123",
		Service: "order-service",
		Event:   coga.EventStart,
		Data:    []byte(`{"message":"lorem ipsum"}`),
	})

	//Output:
	// name: 'saga-orch'
	// body: {ID:123 Service:order-service Event:start Data:[123 34 109 101 115 115 97 103 101 34 58 34 108 111 114 101 109 32 105 112 115 117 109 34 125]}
}
