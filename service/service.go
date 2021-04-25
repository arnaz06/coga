package service

import (
	"context"
	"fmt"

	"github.com/arnaz06/coga"
	log "github.com/sirupsen/logrus"
)

type MessageHandler struct {
	Ctx             context.Context
	transactionList map[string]string
}

func NewMessaheHandler(Ctx context.Context, transactionList map[string]string) MessageHandler {
	return MessageHandler{
		Ctx:             Ctx,
		transactionList: transactionList,
	}
}

func (mh MessageHandler) Handle(m coga.Message) error {
	switch m.Event {
	case coga.EventStart:
		nextServiceName, nextTopicName := mh.ResolveNextService(mh.transactionList, m.Service)
		if nextTopicName == "" {
			return fmt.Errorf("wrong service Name, %s", m.Service)
		}
		m.Service = nextServiceName
		log.Infof("\n publish to %s = %+v \n", nextTopicName, m)
		coga.PublishSystemEvent(mh.Ctx, nextTopicName, m)
	case coga.EventRollback:
		for serviceName, transaction := range mh.transactionList {
			if m.Service == serviceName {
				break
			}
			coga.PublishSystemEvent(mh.Ctx, transaction, coga.Message{
				ID:      m.ID,
				Event:   coga.EventRollback,
				Data:    m.Data,
				Service: serviceName,
			})

			log.Infof("rollback to %s", serviceName)
		}
	default:
		return fmt.Errorf("unsupported event type: %s", m.Event)
	}
	return nil
}

func (mh MessageHandler) ResolveNextService(transactionList map[string]string, serviceName string) (nextServiceName string, nextTopicName string) {
	sliceTl := make([]string, 0)

	for i, _ := range transactionList {
		sliceTl = append(sliceTl, i)
	}

	var nextIndex int
	for i, s := range sliceTl {
		if i+1 == len(sliceTl) {
			fmt.Println(s)
			nextIndex = -1
			break
		}

		if s == serviceName {
			nextIndex = i + 1
			break
		}
	}
	if nextIndex >= 0 {
		nextTopicName = transactionList[sliceTl[nextIndex]]
		nextServiceName = sliceTl[nextIndex]
	}

	return nextServiceName, nextTopicName
}
