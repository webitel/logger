package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/webitel/logger/model"
	"github.com/webitel/wlog"

	amqp "github.com/rabbitmq/amqp091-go"
)

const MaxReconnectAttempts = 10

type RabbitListener struct {
	config            *model.RabbitConfig
	handleFunc        HandleFunc
	connection        *amqp.Connection
	channel           *amqp.Channel
	delivery          <-chan amqp.Delivery
	amqpCloseNotifier chan *amqp.Error
	emergencyStopper  chan model.AppError
	gracefullStopper  chan any
	reconnectAttempts int
}

func BuildRabbit(app *App, config *model.RabbitConfig, errChan chan model.AppError) (*RabbitListener, model.AppError) {
	handler, err := NewHandler(app)
	if err != nil {
		return nil, err
	}
	return NewListener(config, handler.Handle, errChan)

}

func NewListener(config *model.RabbitConfig, f HandleFunc, errChan chan model.AppError) (*RabbitListener, model.AppError) {
	if config == nil {
		return nil, model.NewInternalError("rabbit.listener.new_rabit_listener.arguments_check.config_nil", "rabbit config is nil")
	}

	return &RabbitListener{config: config, handleFunc: f, emergencyStopper: errChan, gracefullStopper: make(chan any)}, nil
}

func (l *RabbitListener) Stop() {
	l.gracefullStopper <- "gracefully"
	if l.channel != nil {
		l.channel.Close()
	}
	if l.connection != nil {
		l.connection.Close()
	}
	defer wlog.Info(fmtBrokerLog("connection closed"))
}

var (
	stopHandler = func(l *RabbitListener) {
		for {
			select {
			case amqpErr, _ := <-l.amqpCloseNotifier:

				// if close has a reason -- log
				if amqpErr != nil {
					wlog.Info(fmtBrokerLog(fmt.Sprintf("reason: %s", amqpErr.Reason)))
				}
				if l.reconnectAttempts >= MaxReconnectAttempts { // if max reconnect attempts reached -- stop execution
					l.emergencyStopper <- model.NewInternalError("app.rabbit_listener.start.reconnect_routine.max_reconnect_attempts.error", fmtBrokerLog("connection lost"))
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
			case <-l.gracefullStopper:
				return
			}

		}
	}

	messageHandler = func(l *RabbitListener) {
		var (
			message amqp.Delivery
			appErr  model.AppError
			err     error
		)
		wlog.Info(fmtBrokerLog("waiting for the messages.."))
		for message = range l.delivery {
			wlog.Info(fmtBrokerLog("message received"))
			ctx, cancelContext := context.WithTimeout(context.Background(), 5*time.Second)
			appErr = l.handleFunc(ctx, &message)
			if appErr != nil {
				if checkNoRows(appErr) { // TODO: real foreign key check by postgres model
					wlog.Debug(fmtBrokerLog(fmt.Sprintf("error processing message, foreign key error, skipping message.. error: %s", appErr.Error())))
					appErr = nil
					err = message.Ack(false)
					if err != nil {
						appErr = model.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
					}
				} else {
					wlog.Debug(fmtBrokerLog(fmt.Sprintf("error while processing the messasge! nacking.. error: %s", appErr.Error())))
					message.Nack(false, true)
				}
				wlog.Info(fmtBrokerLog("message processed with model!"))
			} else {
				wlog.Info(fmtBrokerLog("message processed! acknowledging.. "))
				err = message.Ack(false)
				if err != nil {
					appErr = model.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
				}
			}
			cancelContext()
			if appErr != nil {
				l.emergencyStopper <- appErr
				break
			}

		}
		wlog.Info(fmtBrokerLog("stopped waiting for the messages"))
	}
)

func (l *RabbitListener) Start() {
	var (
		appErr model.AppError
	)
	appErr = l.connect()
	if appErr != nil {
		l.emergencyStopper <- appErr
	}
	go stopHandler(l)
	go messageHandler(l)
	// stop handler
	//go func() {
	//	for {
	//		select {
	//		case amqpErr, _ := <-l.amqpCloseNotifier:
	//
	//			// if close has a reason -- log
	//			if amqpErr != nil {
	//				wlog.Info(fmtBrokerLog(fmt.Sprintf("reason: %s", amqpErr.Reason)))
	//			}
	//			if l.reconnectAttempts >= MaxReconnectAttempts { // if max reconnect attempts reached -- stop execution
	//				l.emergencyStopper <- model.NewInternalError("app.rabbit_listener.start.reconnect_routine.max_reconnect_attempts.error", fmtBrokerLog("connection lost"))
	//				return
	//			}
	//			// try to reconnect
	//			err := l.reconnect()
	//			if err != nil {
	//				wlog.Info(fmtBrokerLog(err.Error()))
	//				l.reconnectAttempts++
	//			} else {
	//				l.reconnectAttempts = 0
	//			}
	//			time.Sleep(time.Second * 10)
	//		case <-l.gracefullStopper:
	//			return
	//		}
	//
	//	}
	//}()
	//
	//// listener
	//go func() {
	//	var (
	//		message amqp.Delivery
	//		appErr  model.AppError
	//	)
	//	wlog.Info(fmtBrokerLog("waiting for the messages.."))
	//	for message = range l.delivery {
	//		wlog.Info(fmtBrokerLog("message received"))
	//		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	//		defer cancelFunc()
	//		appErr = l.handleFunc(ctx, &message)
	//		if appErr != nil {
	//			if checkNoRows(appErr) { // TODO: real foreign key check by postgres model
	//				wlog.Debug(fmtBrokerLog(fmt.Sprintf("error processing message, foreign key error, skipping message.. error: %s", appErr.Error())))
	//				appErr = nil
	//				err = message.Ack(false)
	//				if err != nil {
	//					appErr = model.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
	//				}
	//			} else {
	//				wlog.Debug(fmtBrokerLog(fmt.Sprintf("error while processing the messasge! nacking.. error: %s", appErr.Error())))
	//				message.Nack(false, true)
	//			}
	//			wlog.Info(fmtBrokerLog("message processed with model!"))
	//		} else {
	//			wlog.Info(fmtBrokerLog("message processed! acknowledging.. "))
	//			err = message.Ack(false)
	//			if err != nil {
	//				appErr = model.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
	//			}
	//		}
	//		if appErr != nil {
	//			l.emergencyStopper <- appErr
	//			break
	//		}
	//	}
	//}()
	wlog.Info(fmtBrokerLog("connection opened"))
}

func (l *RabbitListener) connect() model.AppError {
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
		return model.NewInternalError("rabbit.listener.listen.exchange_declare.fail", err.Error())
	}

	wlog.Info(fmtBrokerLog("binding webitel exchange to logger"))
	err = channel.ExchangeBind("logger", "logger.#", "webitel", true, nil)
	if err != nil {
		return model.NewInternalError("rabbit.listener.listen.webitel_exchange_bind.fail", err.Error())
	}

	wlog.Info(fmtBrokerLog("binding chat exchange to logger"))
	err = channel.ExchangeBind("logger", "logger.#", "chat", true, nil)
	if err != nil {
		return model.NewInternalError("rabbit.listener.listen.webitel_exchange_bind.fail", err.Error())
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
		return model.NewInternalError("rabbit.listener.listen.queue_declare.fail", err.Error())
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
		return model.NewInternalError("rabbit.listener.listen.queue_bind.fail", err.Error())
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
		return model.NewInternalError("rabbit.listener.listen.start_consuming.fail", err.Error())
	}
	l.delivery = del
	return nil
}

func (l *RabbitListener) reconnect() model.AppError {
	err := l.connect()
	if err != nil {
		return err
	}
	go messageHandler(l)
	return nil
}

type BrokerHandler struct {
	app *App
}

func NewHandler(app *App) (*BrokerHandler, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("rabbit.handler.new_handler.arguments_check.app_nil", "can't configure handler, app is nil")
	}
	return &BrokerHandler{app: app}, nil
}

func (h *BrokerHandler) Handle(ctx context.Context, message *amqp.Delivery) model.AppError {
	var (
		m      model.BrokerLogMessage
		domain int64
		object string
	)
	err := json.Unmarshal(message.Body, &m)
	if err != nil {
		wlog.Debug(fmt.Sprintf("error unmarshalling message. details: %s", err.Error()))
		return nil
	}

	splittedKey := strings.Split(message.RoutingKey, ".")
	if len(splittedKey) >= 3 {
		domain, _ = strconv.ParseInt(splittedKey[1], 10, 64)
		object = splittedKey[2]
	}
	if m.Records != nil {
		var rabbitMessages []*model.RabbitMessage
		for _, v := range m.Records {
			rabbitMessage := &model.RabbitMessage{
				//ObjectId: object,
				NewState: v.NewState.GetBody(),
				UserId:   m.UserId,
				UserIp:   m.UserIp,
				Action:   m.Action,
				Date:     m.Date,
				RecordId: v.Id,
				Schema:   object,
			}
			rabbitMessages = append(rabbitMessages, rabbitMessage)
		}
		appErr := h.app.InsertLogByRabbitMessageBulk(ctx, rabbitMessages, domain)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

type HandleFunc func(context.Context, *amqp.Delivery) model.AppError

func checkNoRows(err model.AppError) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func fmtBrokerLog(description string) string {
	return fmt.Sprintf("broker: %s", description)
}
