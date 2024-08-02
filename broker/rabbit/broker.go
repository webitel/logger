package rabbit

import (
	"fmt"
	"github.com/webitel/webitel-go-kit/logs"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webitel/logger/model"
)

const MaxReconnectAttempts = 100

type RabbitBroker struct {
	config            *model.RabbitConfig
	connection        *amqp.Connection
	channel           *amqp.Channel
	amqpCloseNotifier chan *amqp.Error
	consumers         map[string]*rabbitQueueConsumer
	emergencyStopper  chan<- model.AppError
	gracefulStopper   chan any
}

func BuildRabbit(config *model.RabbitConfig, errChan chan<- model.AppError) (*RabbitBroker, model.AppError) {
	err := config.Normalize()
	if err != nil {
		return nil, err
	}
	return &RabbitBroker{
		config:           config,
		emergencyStopper: errChan,
		consumers:        make(map[string]*rabbitQueueConsumer),
		gracefulStopper:  make(chan any),
	}, nil

}

// Start starts the channel between rabbitMQ server and this server
func (l *RabbitBroker) Start() model.AppError {
	var (
		appErr model.AppError
	)
	appErr = l.connect()
	if appErr != nil {
		return appErr
	}

	appErr = l.StartAllConsumers()
	if appErr != nil {
		return appErr
	}
	go stopHandler(l)
	return nil
}

// Stop stops all consumers and connections of rabbit
func (l *RabbitBroker) Stop() {
	// send to the gracefulStopper message that signalizes graceful stop (accepted in the stopHandler)
	l.gracefulStopper <- "graceful"
	l.StopAllConsumers()
	if l.channel != nil {
		l.channel.Close()
	}
	if l.connection != nil {
		l.connection.Close()
	}
	defer logs.Info(fmtBrokerLog("connection gracefully closed"))
}

var (
	stopHandler = func(l *RabbitBroker) {
		for {
			select {
			case amqpErr, _ := <-l.amqpCloseNotifier:
				if amqpErr.Reason != "" {
					logs.Info(fmtBrokerLog(fmt.Sprintf("connection lost, %s", amqpErr.Reason)))
				}

				var (
					continueReconnection = true
					reconnectAttempts    int
				)

				for continueReconnection {
					if reconnectAttempts >= MaxReconnectAttempts { // if max reconnect attempts reached -- stop execution
						l.emergencyStopper <- model.NewInternalError("app.broker.stop_handler_routine.reconnect_attempts.reached_limit", "max reconnection attempts")
						return
					}
					reconnectErr := l.reconnect()
					if reconnectErr != nil {
						reconnectAttempts++
						logs.Info(fmtBrokerLog(reconnectErr.Error()))
						//time.Sleep(time.Second * 10)
					} else {
						continueReconnection = false
					}

				}
			case <-l.gracefulStopper:
				return
			}
		}
	}
)

func (l *RabbitBroker) connect() model.AppError {
	conn, err := amqp.Dial(*l.config.Url)
	if err != nil {
		return model.NewInternalError("rabbit.listener.listen.server_connect.fail", err.Error())
	}
	l.connection = conn
	channel, err := conn.Channel()
	if err != nil {
		return model.NewInternalError("rabbit.listener.listen.channel_connect.fail", err.Error())
	}
	l.channel = channel
	l.amqpCloseNotifier = l.channel.NotifyClose(make(chan *amqp.Error))

	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("basic.qos: %v", err)
	} // log.Fatalf?
	logs.Info(fmtBrokerLog("connection and amqp channel are opened"))
	return nil
}

func (l *RabbitBroker) reconnect() model.AppError {
	// try to create new connection channel
	logs.Debug(fmtBrokerLog("trying to reconnect"))
	err := l.connect()
	if err != nil {
		return err
	}
	for s, consumer := range l.consumers {
		// make a new delivery channel with new connection
		ch, err := l.Consume(s, consumer.name)
		if err != nil {
			return err
		}
		consumer.delivery = ch
		// start listen to the new delivery channel
		consumer.Start()
	}
	return nil
}

// StopAllConsumers stops all consumers if exist
func (l *RabbitBroker) StopAllConsumers() model.AppError {
	for _, consumer := range l.consumers {
		consumer.Stop()
	}
	return nil
}

// StartAllConsumers starts all consumers if exist
func (l *RabbitBroker) StartAllConsumers() model.AppError {
	for _, consumer := range l.consumers {
		err := consumer.Start()
		if err != nil {
			return err
		}
	}
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
	logs.Info(fmtBrokerLog(fmt.Sprintf("[%s] exchange declared", exchangeName)))
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
	logs.Info(fmtBrokerLog(fmt.Sprintf("[%s] queue declared", queueName)))
	return queueName, nil
}

func (l *RabbitBroker) QueueBind(exchangeName string, queueName string, routingKey string, noWait bool, args map[string]any) model.AppError {
	err := l.channel.QueueBind(queueName, routingKey, exchangeName, noWait, args)
	if err != nil {
		return model.NewInternalError("rabbit.listener.queue_bind.request.fail", err.Error())
	}
	logs.Info(fmtBrokerLog(fmt.Sprintf("[%s]->(%s)->[%s] queue bind", exchangeName, routingKey, queueName)))
	return nil
}

func (l *RabbitBroker) Consume(queueName string, consumerName string) (<-chan amqp.Delivery, model.AppError) {
	ch, err := l.channel.Consume(queueName, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, model.NewInternalError("rabbit.listener.consume.request.fail", err.Error())
	}
	logs.Info(fmtBrokerLog(fmt.Sprintf("[%s] queue started to consume", queueName)))
	return ch, nil
}

func (l *RabbitBroker) ExchangeBind(destination string, key string, source string, noWait bool, args map[string]any) model.AppError {
	err := l.channel.ExchangeBind(destination, key, source, noWait, args)
	if err != nil {
		return model.NewInternalError("rabbit.listener.exchange_bind.request.fail", err.Error())
	}
	logs.Info(fmtBrokerLog(fmt.Sprintf("[%s]->(%s)->[%s] exchange binded", source, key, destination)))
	return nil
}

func (l *RabbitBroker) QueueStartConsume(queueName string, consumerName string, handleFunc HandleFunc, handleTimeout time.Duration) model.AppError {
	// make a connection
	ch, err := l.Consume(queueName, consumerName)
	if err != nil {
		return err
	}
	// initialize handler
	queue, err := BuildRabbitQueueConsumer(ch, handleFunc, consumerName, handleTimeout)
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
