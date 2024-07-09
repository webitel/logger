package rabbit

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webitel/logger/model"
	"time"
)

type rabbitQueueConsumer struct {
	handleFunc      HandleFunc
	handleTimeout   time.Duration
	acknowledgeFunc AcknowledgeFunc
	delivery        <-chan amqp.Delivery
	stopper         chan any
	consumerName    string
}

func BuildRabbitQueueConsumer(delivery <-chan amqp.Delivery, acknowledgeFunc AcknowledgeFunc, handleFunc HandleFunc, consumerName string, handleTimeout time.Duration) (*rabbitQueueConsumer, model.AppError) {
	if acknowledgeFunc == nil {
		return nil, model.NewInternalError("rabbit.consumer.build.check_args.handle_function", "acknowledge function not specified")
	}
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
		handleTimeout:   handleTimeout,
		handleFunc:      handleFunc,
		acknowledgeFunc: acknowledgeFunc,
		delivery:        delivery,
		stopper:         make(chan any),
		consumerName:    consumerName,
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
	go l.acknowledgeFunc(l.handleTimeout, l.delivery, l.stopper, l.handleFunc)
	return nil
}

// HandleFunc allows to define the reaction to the amqp.Delivery
type HandleFunc func(context.Context, *amqp.Delivery) model.AppError

/*
AcknowledgeFunc allows to define the reaction to the amqp.Delivery.

Will run in goroutine and should handle logic for the acknowledging messages.

delivery - channel where amqp.messages will be delivered

stopper - channel for stopping the routine

handleFunc - function used to handle the exact amqp.message content
*/
type AcknowledgeFunc func(handleTimeout time.Duration, delivery <-chan amqp.Delivery, stopper chan any, handleFunc HandleFunc)
