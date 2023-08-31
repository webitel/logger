package rabbit

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/webitel/logger/app"
	"github.com/webitel/logger/model"
	"github.com/webitel/wlog"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	errors "github.com/webitel/engine/model"
)

type HandleFunc func(context.Context, *amqp.Delivery) errors.AppError

type RabbitListener struct {
	config     *model.RabbitConfig
	handleFunc HandleFunc
	connection *amqp.Connection
	channel    *amqp.Channel
	exit       chan errors.AppError
}

//func BuildAndServe(app *app.App, config *model.RabbitConfig, errChan chan errors.AppError) {
//	handler, err := NewHandler(app)
//	if err != nil {
//		errChan <- err
//		return
//	}
//	listener, err := NewListener(config, errChan)
//	if err != nil {
//		errChan <- err
//		return
//	}
//
//	listener.Start()
//}

func Build(app *app.App, config *model.RabbitConfig, errChan chan errors.AppError) (*RabbitListener, errors.AppError) {
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
	l.channel.Cancel("", true)
}

func (l *RabbitListener) Start() {
	conn, err := amqp.Dial(l.config.Url)
	if err != nil {
		l.exit <- errors.NewInternalError("rabbit.listener.listen.server_connect.fail", err.Error())
		return
	}
	l.connection = conn
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		l.exit <- errors.NewInternalError("rabbit.listener.listen.channel_connect.fail", err.Error())
		return
	}
	defer channel.Close()
	wlog.Info("connecting to the exchange")
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
		l.exit <- errors.NewInternalError("rabbit.listener.listen.exchange_declare.fail", err.Error())
		return
	}
	// queue, err := channel.QueueDeclarePassive(
	wlog.Info("connecting or creating a queue 'logger.service'")
	queue, err := channel.QueueDeclare(
		"logger.service",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		l.exit <- errors.NewInternalError("rabbit.listener.listen.queue_declare.fail", err.Error())
		return
	}
	wlog.Info("binding queue..")
	err = channel.QueueBind(
		queue.Name, // queue name
		"logger.#", // routing key
		"logger",   // exchange
		false,
		nil,
	)
	if err != nil {
		l.exit <- errors.NewInternalError("rabbit.listener.listen.queue_bind.fail", err.Error())
		return
	}
	err = channel.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("basic.qos: %v", err)
	}
	delivery, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		l.exit <- errors.NewInternalError("rabbit.listener.listen.start_consuming.fail", err.Error())
		return
	}
	//var forever chan struct{}
	go func() {
		var (
			message amqp.Delivery
			appErr  errors.AppError
		)
		wlog.Info("waiting for the messages..")
		for message = range delivery {
			wlog.Info("message received")
			appErr = l.handleFunc(context.Background(), &message)
			if appErr != nil {
				if checkNoRows(appErr) {
					wlog.Debug(fmt.Sprintf("error processing message, foreign key error, skipping message.. error: %s", appErr.Error()))
					appErr = nil
					err = message.Ack(false)
					if err != nil {
						appErr = errors.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
					}
				} else {
					wlog.Debug(fmt.Sprintf("error while processing the messasge! nacking.. error: %s", appErr.Error()))
					message.Nack(false, true)
				}
				wlog.Info("message processed with errors!")
			} else {
				wlog.Info("message processed! acknowledging.. ")
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

func checkNoRows(err errors.AppError) bool {
	return strings.Contains(err.Error(), sql.ErrNoRows.Error())
}
