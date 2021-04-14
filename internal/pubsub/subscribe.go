package pubsub

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	"cloud.google.com/go/pubsub"
	"github.com/arnaz06/coga"
	"github.com/arnaz06/coga/service"
	log "github.com/sirupsen/logrus"
)

type subMessageHandler struct {
	subscription   *pubsub.Subscription
	messageHandler service.MessageHandler
}

func NewSubMessageHandler(subscription *pubsub.Subscription, messageHandler service.MessageHandler) subMessageHandler {
	return subMessageHandler{
		subscription:   subscription,
		messageHandler: messageHandler,
	}
}

func (smh subMessageHandler) Pull() {
	go func() {
		err := smh.subscription.Receive(smh.messageHandler.Ctx, func(ctx context.Context, message *pubsub.Message) {
			var cogaMessage map[string]interface{}
			err := json.Unmarshal([]byte(message.Data), &cogaMessage)
			if err != nil {
				message.Ack()
				log.Errorf("error unmarshal message.Data, err:%+v, data:%+v", err, message.Data)
			}
			err = smh.messageHandler.Handle(coga.ToSystemEvent(cogaMessage).Body)
			if err != nil {
				message.Ack()
				log.Errorf("error handle message, err:%+v", err)
			}
			message.Ack()
		})

		if err != nil {
			log.Errorf("error pull message from pubsub, err:%+v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
}
