package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	notifier "github.com/webitel/webitel-go-kit/pkg/watcher"
)

type Publisher interface {
	Publish(ctx context.Context, routingKey string, body []byte, headers amqp091.Table) error
}

type NoopPublisher struct{}

func (n *NoopPublisher) Publish(_ context.Context, _ string, _ []byte, _ amqp091.Table) error {
	return nil
}

type CrudObserver struct {
	publisher Publisher
}

func NewCrudObserver(publisher Publisher) *CrudObserver {
	pub := publisher
	if pub == nil {
		pub = &NoopPublisher{}
	}
	return &CrudObserver{publisher: pub}
}

func (l *CrudObserver) Update(eventType notifier.EventType, m map[string]any) error {
	obj, ok := m["object"]
	if !ok {
		return errors.New("object not found")
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	succeded, ok := m["succeeded"]
	if !ok {
		return errors.New("succeded not found")
	}
	var operationStatus string
	if succeded.(bool) {
		operationStatus = "success"
	} else {
		operationStatus = "failure"
	}

	err = l.publisher.Publish(context.Background(), fmt.Sprintf("%s_%s.%s", m["objclass"], eventType, operationStatus), bytes, nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *CrudObserver) GetId() string {
	return "CrudObserver"
}
