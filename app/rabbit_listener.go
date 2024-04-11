package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/webitel/logger/model"
	"github.com/webitel/wlog"

	amqp "github.com/rabbitmq/amqp091-go"

	errors "github.com/webitel/engine/model"
)

const MaxReconnectAttempts = 50

type HandleFunc func(context.Context, *amqp.Delivery) errors.AppError

type RabbitListener struct {
	config            *model.RabbitConfig
	handleFunc        HandleFunc
	connection        *amqp.Connection
	channel           *amqp.Channel
	delivery          <-chan amqp.Delivery
	amqpCloseNotifier chan *amqp.Error
	exit              chan errors.AppError
	reconnectAttempts int
}

func BuildRabbit(app *App, config *model.RabbitConfig, errChan chan errors.AppError) (*RabbitListener, errors.AppError) {
	handler, err := NewHandler(app)
	if err != nil {
		return nil, err
	}
	return NewListener(config, handler.Handle, errChan)

}

func NewListener(config *model.RabbitConfig, f HandleFunc, errChan chan errors.AppError) (*RabbitListener, errors.AppError) {
	if config == nil {
		return nil, errors.NewInternalError("rabbit.listener.new_rabit_listener.arguments_check.config_nil", "rabbit config is nil")
	}

	return &RabbitListener{config: config, handleFunc: f, exit: errChan}, nil
}

func (l *RabbitListener) Stop() {
	if l.channel != nil {
		l.channel.Close()
	}
	if l.connection != nil {
		l.connection.Close()
	}
}

func (l *RabbitListener) Start() {
	var (
		err    error
		appErr errors.AppError
	)
	appErr = l.connect()
	if appErr != nil {
		l.exit <- appErr
	}

	go func() {
		for {
			// check and wait for closed channel
			amqpErr, _ := <-l.amqpCloseNotifier

			wlog.Info(fmtBrokerLog("connection closed... "))
			// if close has a reason -- log
			if amqpErr != nil {
				wlog.Info(fmtBrokerLog(fmt.Sprintf("reason: %s", amqpErr.Reason)))
			}
			if l.reconnectAttempts >= MaxReconnectAttempts { // if max reconnect attempts reached -- stop execution
				l.exit <- errors.NewInternalError("app.rabbit_listener.start.reconnect_routine.max_reconnect_attempts.error", fmtBrokerLog("connection lost"))
				return
			}
			// try to reconnect
			err := l.reconnect()
			if err != nil {
				wlog.Info(fmtBrokerLog(err.Error()))
				l.reconnectAttempts++
			} else {
				l.reconnectAttempts = 0
			}
			time.Sleep(time.Second * 10)
		}
	}()
	//var forever chan struct{}
	go func() {
		var (
			message amqp.Delivery
			appErr  errors.AppError
		)
		wlog.Info(fmtBrokerLog("waiting for the messages.."))
		for message = range l.delivery {
			wlog.Info(fmtBrokerLog("message received"))
			appErr = l.handleFunc(context.Background(), &message)
			if appErr != nil {
				if checkNoRows(appErr) { // TODO: real foreign key check by postgres errors
					wlog.Debug(fmtBrokerLog(fmt.Sprintf("error processing message, foreign key error, skipping message.. error: %s", appErr.Error())))
					appErr = nil
					err = message.Ack(false)
					if err != nil {
						appErr = errors.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
					}
				} else {
					wlog.Debug(fmtBrokerLog(fmt.Sprintf("error while processing the messasge! nacking.. error: %s", appErr.Error())))
					message.Nack(false, true)
				}
				wlog.Info(fmtBrokerLog("message processed with errors!"))
			} else {
				wlog.Info(fmtBrokerLog("message processed! acknowledging.. "))
				err = message.Ack(false)
				if err != nil {
					appErr = errors.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
				}
			}
			if appErr != nil {
				l.exit <- appErr
				break
			}
		}
	}()
	<-l.exit

}

func (l *RabbitListener) connect() errors.AppError {
	conn, err := amqp.Dial(l.config.Url)
	if err != nil {
		return errors.NewInternalError("rabbit.listener.listen.server_connect.fail", err.Error())
	}
	l.connection = conn
	channel, err := conn.Channel()
	if err != nil {
		return errors.NewInternalError("rabbit.listener.listen.channel_connect.fail", err.Error())
	}
	l.channel = channel
	l.amqpCloseNotifier = make(chan *amqp.Error)
	l.channel.NotifyClose(l.amqpCloseNotifier)
	wlog.Info(fmtBrokerLog("connecting to the exchange"))
	err = channel.ExchangeDeclare(
		"logger", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return errors.NewInternalError("rabbit.listener.listen.exchange_declare.fail", err.Error())
	}
	// queue, err := channel.QueueDeclarePassive(
	wlog.Info(fmtBrokerLog("binding webitel exchange to logger"))
	err = channel.ExchangeBind("logger", "logger.#", "webitel", true, nil)
	if err != nil {
		return errors.NewInternalError("rabbit.listener.listen.webitel_exchange_bind.fail", err.Error())
	}
	wlog.Info(fmtBrokerLog("connecting or creating a queue 'logger.service'"))
	queue, err := channel.QueueDeclare(
		"logger.service",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.NewInternalError("rabbit.listener.listen.queue_declare.fail", err.Error())
	}
	wlog.Info(fmtBrokerLog("binding queue.."))
	err = channel.QueueBind(
		queue.Name, // queue name
		"logger.#", // routing key
		"logger",   // exchange
		false,
		nil,
	)
	if err != nil {
		return errors.NewInternalError("rabbit.listener.listen.queue_bind.fail", err.Error())
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("basic.qos: %v", err)
	}
	del, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return errors.NewInternalError("rabbit.listener.listen.start_consuming.fail", err.Error())
	}
	l.delivery = del
	return nil
}

func (l *RabbitListener) reconnect() errors.AppError {
	err := l.connect()
	if err != nil {
		return err
	}
	return nil
}

func checkNoRows(err errors.AppError) bool {
	return strings.Contains(err.Error(), sql.ErrNoRows.Error())
}

func fmtBrokerLog(description string) string {
	return fmt.Sprintf("broker: %s", description)
}
