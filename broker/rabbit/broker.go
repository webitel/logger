package rabbit

import (
	"fmt"
	"log"
	"time"

	"github.com/webitel/logger/model"
	"github.com/webitel/wlog"

	amqp "github.com/rabbitmq/amqp091-go"
)

const MaxReconnectAttempts = 10

type RabbitBroker struct {
	config            *model.RabbitConfig
	connection        *amqp.Connection
	channel           *amqp.Channel
	amqpCloseNotifier chan *amqp.Error
	consumers         map[string]*rabbitQueueConsumer
	emergencyStopper  chan model.AppError
	gracefulStopper   chan any
	reconnectAttempts int
}

func BuildRabbit(config *model.RabbitConfig, errChan chan model.AppError) (*RabbitBroker, model.AppError) {
	return &RabbitBroker{
		config:           config,
		emergencyStopper: errChan,
		consumers:        make(map[string]*rabbitQueueConsumer),
		gracefulStopper:  make(chan any),
	}, nil

}

// Start starts the channel between rabbitMQ server and this server
func (l *RabbitBroker) Start() {
	var (
		appErr model.AppError
	)
	appErr = l.connect()
	if appErr != nil {
		l.emergencyStopper <- appErr
	}
	go stopHandler(l)
	wlog.Info(fmtBrokerLog("connection opened"))
}

// Stop stops all consumers and connections of rabbit
func (l *RabbitBroker) Stop() {
	// send to the gracefulStopper message that signalizes graceful stop
	l.gracefulStopper <- "graceful"
	for _, consumer := range l.consumers {
		consumer.Stop()
	}
	if l.channel != nil {
		l.channel.Close()
	}
	if l.connection != nil {
		l.connection.Close()
	}
	defer wlog.Info(fmtBrokerLog("connection closed"))
}

var (
	stopHandler = func(l *RabbitBroker) {
		for {
			select {
			case amqpErr, _ := <-l.amqpCloseNotifier:
				if amqpErr.Reason != "" {
					wlog.Info(fmtBrokerLog(fmt.Sprintf("reason: %s", amqpErr.Reason)))
				}

				for {
					if l.reconnectAttempts >= MaxReconnectAttempts { // if max reconnect attempts reached -- stop execution
						l.Stop()
						l.emergencyStopper <- model.NewInternalError("app.broker.stop_handler_routine.reconnect_attempts.reached_limit", "max reconnection attempts")
						return // end goroutine execution
					}
					reconnectErr := l.reconnect()
					if reconnectErr != nil {
						wlog.Info(fmtBrokerLog(reconnectErr.Error()))
					} else {
						l.reconnectAttempts = 0
					}
					time.Sleep(time.Second * 10)
				}

			case <-l.emergencyStopper:
				l.Stop()
				return
			case <-l.gracefulStopper:
				return
			}
		}
	}
)

func (l *RabbitBroker) connect() model.AppError {
	conn, err := amqp.Dial(l.config.Url)
	if err != nil {
		return model.NewInternalError("rabbit.listener.listen.server_connect.fail", err.Error())
	}
	l.connection = conn
	channel, err := conn.Channel()
	if err != nil {
		return model.NewInternalError("rabbit.listener.listen.channel_connect.fail", err.Error())
	}
	l.channel = channel
	l.amqpCloseNotifier = make(chan *amqp.Error)
	l.channel.NotifyClose(l.amqpCloseNotifier)

	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("basic.qos: %v", err)
	}
	return nil
}

func (l *RabbitBroker) reconnect() model.AppError {
	// try to create new connection channel
	err := l.connect()
	if err != nil {
		l.reconnectAttempts++
		return err
	}
	for s, consumer := range l.consumers {
		// make a new delivery channel with new connection
		ch, err := l.Consume(s, consumer.consumerName)
		if err != nil {
			return err
		}
		consumer.delivery = ch
		// start listen to the new delivery channel
		consumer.Start()
	}
	l.reconnectAttempts = 0
	return nil
}

func (l *RabbitBroker) ExchangeDeclare(exchangeName string, kind string, opts ...ExchangeDeclareOption) model.AppError {
	var decarationOptions ExchangeDeclareOptions
	for _, opt := range opts {
		opt(&decarationOptions)
	}

	err := l.channel.ExchangeDeclare(
		exchangeName,                 // name
		kind,                         // type
		decarationOptions.Durable,    // durable
		decarationOptions.AutoDelete, // auto-deleted
		decarationOptions.Internal,   // internal
		decarationOptions.NoWait,     // no-wait
		decarationOptions.Args,       // arguments
	)
	if err != nil {
		return model.NewInternalError("rabbit.listener.exchange_declare.request.fail", err.Error())
	}
	wlog.Info(fmtBrokerLog(fmt.Sprintf("[%s] exchange declared", exchangeName)))
	return nil
}

func (l *RabbitBroker) QueueDeclare(queueName string, opts ...QueueDeclareOption) (string, model.AppError) {
	var declarationOptions QueueDeclareOptions
	for _, opt := range opts {
		opt(&declarationOptions)
	}

	_, err := l.channel.QueueDeclare(
		queueName,
		declarationOptions.Durable,
		declarationOptions.AutoDelete,
		declarationOptions.Exclusive,
		declarationOptions.NoWait,
		declarationOptions.Args,
	)
	if err != nil {
		return "", model.NewInternalError("rabbit.listener.queue_declare.request.fail", err.Error())
	}
	wlog.Info(fmtBrokerLog(fmt.Sprintf("[%s] queue declared", queueName)))
	return queueName, nil
}

func (l *RabbitBroker) QueueBind(exchangeName string, queueName string, routingKey string, noWait bool, args map[string]any) model.AppError {
	err := l.channel.QueueBind(queueName, routingKey, exchangeName, noWait, args)
	if err != nil {
		return model.NewInternalError("rabbit.listener.queue_bind.request.fail", err.Error())
	}
	wlog.Info(fmtBrokerLog(fmt.Sprintf("[%s]->(%s)->[%s] queue bind", exchangeName, routingKey, queueName)))
	return nil
}

func (l *RabbitBroker) Consume(queueName string, consumerName string) (<-chan amqp.Delivery, model.AppError) {
	ch, err := l.channel.Consume(queueName, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, model.NewInternalError("rabbit.listener.consume.request.fail", err.Error())
	}
	wlog.Info(fmtBrokerLog(fmt.Sprintf("[%s] queue started to consume", queueName)))
	return ch, nil
}

func (l *RabbitBroker) ExchangeBind(destination string, key string, source string, noWait bool, args map[string]any) model.AppError {
	err := l.channel.ExchangeBind(destination, key, source, noWait, args)
	if err != nil {
		return model.NewInternalError("rabbit.listener.exchange_bind.request.fail", err.Error())
	}
	wlog.Info(fmtBrokerLog(fmt.Sprintf("[%s]->(%s)->[%s] exchange binded", source, key, destination)))
	return nil
}

func (l *RabbitBroker) QueueStartConsume(queueName string, consumerName string, acknowledgeFunc AcknowledgeFunc, handleFunc HandleFunc, handleTimeout time.Duration) model.AppError {
	// make a connection
	ch, err := l.Consume(queueName, consumerName)
	if err != nil {
		return err
	}
	// initialize handler
	queue, err := BuildRabbitQueueConsumer(ch, acknowledgeFunc, handleFunc, consumerName, handleTimeout)
	if err != nil {
		return err
	}
	// start new consumer
	queue.Start()

	// insert handler in the registry
	l.consumers[queueName] = queue

	return nil
}

func (l *RabbitBroker) QueueStopConsume(queueName string) model.AppError {
	if consumer, consumerFound := l.consumers[queueName]; consumerFound {
		consumer.Stop()
		delete(l.consumers, queueName)
	}
	return nil
}

func fmtBrokerLog(description string) string {
	return fmt.Sprintf("broker: %s", description)
}
