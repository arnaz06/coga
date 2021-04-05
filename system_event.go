package coga

import (
	"context"
	"encoding/json"
	"reflect"
	"time"
)

type contextKey int

const (
	// ContextKeyPublisher is the context key for publisher.
	ContextKeyPublisher contextKey = iota
)

// Publisher is the interface that wraps the basic Publish method.
type Publisher interface {
	// Publish publishes system events.
	Publish(ctx context.Context, e SystemEvent)
}

// SystemEvent holds information of system event.
type SystemEvent struct {
	Name        string
	Body        Message
	PublishTime time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (e SystemEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":        e.Name,
		"body":        e.Body,
		"publishTime": e.PublishTime.Format(time.RFC3339Nano),
	})
}

// NewSystemEvent creates a system event using name inferred from the eventBody type name.
func NewSystemEvent(eventBody Message) SystemEvent {
	name := reflect.TypeOf(eventBody).Name()
	return SystemEvent{
		Name:        name,
		Body:        eventBody,
		PublishTime: time.Now(),
	}
}

// PublisherFromContext get Publisher from the ctx.
func PublisherFromContext(ctx context.Context) Publisher {
	pub, ok := ctx.Value(ContextKeyPublisher).(Publisher)
	if !ok {
		return nil
	}

	return pub
}

// PublishSystemEvent publishes a system event.
func PublishSystemEvent(ctx context.Context, eventBody Message) {
	e := NewSystemEvent(eventBody)

	publisher := PublisherFromContext(ctx)
	if publisher == nil {
		return
	}

	publisher.Publish(ctx, e)
}
