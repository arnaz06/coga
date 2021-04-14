package pubsub

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/arnaz06/coga"
	log "github.com/sirupsen/logrus"
)

type eventPublishingService struct {
	ctx   context.Context
	topic *pubsub.Topic
}

func NewEventPublishingService(ctx context.Context, topic *pubsub.Topic) coga.Publisher {
	return eventPublishingService{
		ctx:   ctx,
		topic: topic,
	}
}

func (ep eventPublishingService) Publish(ctx context.Context, systemEvent coga.SystemEvent) {
	messageJSON, err := json.Marshal(systemEvent)
	if err != nil {
		log.Errorf("error publish event, err:%+v event:%+v", err, systemEvent)
		return
	}

	log.Infof("publishing event: %+v", systemEvent)
	res := ep.topic.Publish(ep.ctx, &pubsub.Message{Data: messageJSON, PublishTime: systemEvent.PublishTime})
	id, err := res.Get(ep.ctx)
	if err != nil {
		log.Errorf("error get pubsub messageID, err:%+v", err)
	} else {
		log.Debugf("published event %+v with id: %s", systemEvent, id)
	}
}
