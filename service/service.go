package service

import (
	"context"
	"fmt"

	"github.com/arnaz06/coga"
	log "github.com/sirupsen/logrus"
)

type MessageHandler struct {
	Ctx             context.Context
	transactionList []coga.TransactionList
}

func NewMessaheHandler(Ctx context.Context, transactionList []coga.TransactionList) MessageHandler {
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
		log.Infof("\n publish to  ===== %s ===== \n", nextTopicName)
		coga.PublishSystemEvent(mh.Ctx, nextTopicName, m)
	case coga.EventRollback:
		tl := mh.ResolveRollback(mh.transactionList, m.Service)
		for _, t := range tl {
			coga.PublishSystemEvent(mh.Ctx, t.Topic, coga.Message{
				ID:      m.ID,
				Event:   coga.EventRollback,
				Data:    m.Data,
				Service: t.ServiceName,
			})

			log.Infof("rollback to %s", t.ServiceName)
		}
	default:
		return fmt.Errorf("unsupported event type: %s", m.Event)
	}
	return nil
}

func (mh MessageHandler) ResolveRollback(tl []coga.TransactionList, serviceName string) (res []coga.TransactionList) {
	res = make([]coga.TransactionList, 0)
	for _, t := range tl {
		if t.ServiceName == "coga-service" {
			continue
		}
		if t.ServiceName == serviceName {
			break
		}
		res = append(res, coga.TransactionList{
			ServiceName: t.ServiceName,
			Topic:       t.Topic,
		})
	}
	return res
}

func (mh MessageHandler) ResolveNextService(tl []coga.TransactionList, serviceName string) (nextServiceName string, nextTopicName string) {
	temp := 0
	for i, t := range tl {
		if i+1 == len(tl) {
			temp = -1
			break
		}
		if t.ServiceName == serviceName {
			temp += 1
			break
		}
		temp += 1
	}

	if temp >= 0 {
		nextServiceName = tl[temp].ServiceName
		nextTopicName = tl[temp].Topic
	}
	return nextServiceName, nextTopicName
}
