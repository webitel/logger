package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webitel/logger/internal/model"
	"time"
)

type rabbitQueueConsumer struct {
	handleFunc    HandleFunc
	handleTimeout time.Duration
	delivery      <-chan amqp.Delivery
	stopper       chan any
	name          string
}

func BuildRabbitQueueConsumer(delivery <-chan amqp.Delivery, handleFunc HandleFunc, consumerName string, handleTimeout time.Duration) (*rabbitQueueConsumer, model.AppError) {
	if handleFunc == nil {
		return nil, model.NewInternalError("rabbit.consumer.build.check_args.handle_function", "handle function not specified")
	}
	if delivery == nil {
		return nil, model.NewInternalError("rabbit.consumer.build.check_args.delivery_channel", "delivery channel is nil")
	}
	if handleTimeout == 0 {
		handleTimeout = 5 * time.Second
	}
	return &rabbitQueueConsumer{
		handleTimeout: handleTimeout,
		handleFunc:    handleFunc,
		delivery:      delivery,
		stopper:       make(chan any),
		name:          consumerName,
	}, nil
}

func (l *rabbitQueueConsumer) Stop() {
	l.stopper <- "gracefully"
	return
}

func (l *rabbitQueueConsumer) Start() model.AppError {
	if l.delivery == nil {
		return model.NewInternalError("rabbit.consumer.start.check_args.delivery_channel", "delivery channel is nil")
	}
	if l.handleFunc == nil {
		return model.NewInternalError("rabbit.consumer.start.check_args.handle_func", "handle function not specified")
	}
	if l.stopper == nil {
		return model.NewInternalError("rabbit.consumer.start.check_args.stopper_channel", "stopper channel is nil")
	}
	go l.handleFunc(l.handleTimeout, l.delivery, l.stopper)
	return nil
}

/*
AcknowledgeFunc allows to define the reaction to the amqp.Delivery.

Will run in goroutine and should handle logic for the acknowledging messages.

delivery - channel where amqp.messages will be delivered

stopper - channel for stopping the routine

handleFunc - function used to handle the exact amqp.message content
*/
type HandleFunc func(handleTimeout time.Duration, delivery <-chan amqp.Delivery, stopper chan any)
