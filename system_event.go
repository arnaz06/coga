package coga

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

// Publisher is the interface that wraps the basic Publish method.
type Publisher interface {
	// Publish publishes system events.
	Publish(ctx context.Context, e SystemEvent)
}

// SystemEvent represent of system event data.
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

// NewSystemEvent ...
func NewSystemEvent(eventBody Message) SystemEvent {
	return SystemEvent{
		Name:        "saga-orch",
		Body:        eventBody,
		PublishTime: time.Now(),
	}
}

func ToSystemEvent(message map[string]interface{}) SystemEvent {
	var res SystemEvent
	res.Name, _ = message["name"].(string)
	body, _ := message["body"].(map[string]interface{})
	res.Body.ID, _ = body["id"].(string)
	res.Body.Service, _ = body["service"].(string)
	res.Body.Event, _ = body["event"].(string)
	res.Body.Data, _ = body["data"].(map[string]interface{})
	return res
}

// PublisherFromContext is function to get Publisher from context.
func PublisherFromContext(ctx context.Context, topicKey string) Publisher {
	pub, ok := ctx.Value(topicKey).(Publisher)
	if !ok {
		log.Error("publisher not found")
		return nil
	}

	return pub
}

// PublishSystemEvent publishes a system event.
func PublishSystemEvent(ctx context.Context, topicKey string, eventBody Message) {
	e := NewSystemEvent(eventBody)

	publisher := PublisherFromContext(ctx, topicKey)
	if publisher == nil {
		log.Error("publisher not found")
		return
	}

	publisher.Publish(ctx, e)
}
