package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webitel/logger/internal/model"
	broker "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq"
	slogadapter "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq/pkg/adapter/slog"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

func (a *App) initLogsConsumption(exchangeConfig *broker.ExchangeConfig) error {
	// bind all logs from webitel exchange to the logger exchange
	err := a.rabbitConn.BindExchange(context.Background(), "webitel", "logger", "logger.#", true, nil)
	if err != nil {
		return err
	}
	// declare new queue logger.service
	queueConfig, err := broker.NewQueueConfig("logger.service", broker.WithQueueDurable(true))
	if err != nil {
		return err
	}
	err = a.rabbitConn.DeclareQueue(context.Background(), queueConfig, exchangeConfig, "logger.#")
	if err != nil {
		return err
	}
	consumerConf, err := broker.NewConsumerConfig(fmt.Sprintf("logger_logs_%s", a.config.Consul.Id))
	if err != nil {
		return err
	}
	consumer := broker.NewConsumer(a.rabbitConn, queueConfig, consumerConf, a.HandleLog, slogadapter.NewSlogLogger(slog.Default()))
	return consumer.Start(context.Background())
}

func (a *App) HandleLog(ctx context.Context, message amqp.Delivery) error {
	var (
		m      model.BrokerRecordLogMessage
		domain int64
		object string
	)
	err := json.Unmarshal(message.Body, &m)
	if err != nil {
		slog.Debug(fmt.Sprintf("error unmarshalling message. details: %s", err.Error()))
		return nil
	}

	splittedKey := strings.Split(message.RoutingKey, ".")
	if len(splittedKey) >= 3 {
		domain, _ = strconv.ParseInt(splittedKey[1], 10, 64)
		object = splittedKey[2]
	}
	for _, record := range m.Records {
		log := &model.Log{
			Action:   m.Action,
			UserIp:   &m.UserIp,
			NewState: record.NewState.Body,
			Object:   &model.Object{Name: &object},
			Author:   &model.Author{Id: &m.UserId},
			Record:   &model.Record{Id: &record.Id},
		}
		date := time.Unix(m.Date, 0)
		if m.Date == 0 {
			date = time.Now()
		}
		log.Date = &date

		err = a.CreateLog(ctx, log, int(domain))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initLoginConsumption() error {
	sourceExchangeName := "webitel"
	queueName := "logger.login"
	// create or connect the logger exchange
	exchangeConf, err := broker.NewExchangeConfig(sourceExchangeName, broker.ExchangeTypeTopic)
	if err != nil {
		return err
	}
	queueConf, err := broker.NewQueueConfig(queueName, broker.WithQueueDurable(true))
	if err != nil {
		return err
	}
	consConfig, err := broker.NewConsumerConfig(fmt.Sprintf("logger_login_%s", a.config.Consul.Id))
	if err != nil {
		return err
	}
	err = a.rabbitConn.DeclareQueue(context.Background(), queueConf, exchangeConf, "login.#")
	if err != nil {
		return err
	}
	cons := broker.NewConsumer(a.rabbitConn, queueConf, consConfig, a.HandleLogin, slogadapter.NewSlogLogger(slog.Default()))
	return cons.Start(context.Background())
}

func (a *App) HandleLogin(ctx context.Context, message amqp.Delivery) error {
	var (
		m model.BrokerLoginMessage
	)
	err := json.Unmarshal(message.Body, &m)
	if err != nil {
		return err
	}

	splittedKey := strings.Split(message.RoutingKey, ".")
	if len(splittedKey) < 4 {
		return errors.New("provided routing key is not matching with this handler")
	}

	databaseModel, appErr := m.ConvertToDatabaseModel()
	if appErr != nil {
		return appErr
	}

	_, err = a.storage.LoginAttempt().Insert(ctx, databaseModel)
	if err != nil {
		return err
	}
	return nil
}

type PopulateConfigEventRequest struct {
	DomainId int `json:"domainId,omitempty"`
}

type Config struct {
	Name  string `json:"name,omitempty"`
	State bool   `json:"state,omitempty"`
}

type PopulateConfigEventResponse struct {
	DomainId int       `json:"domainId,omitempty"`
	Configs  []*Config `json:"configs,omitempty"`
}

func (a *App) initPopulateEventConsumption(exchangeConfig *broker.ExchangeConfig) error {
	// create or connect the logger exchange
	queueConf, err := broker.NewQueueConfig("logger.populate", broker.WithQueueDurable(true))
	if err != nil {
		return err
	}
	err = a.rabbitConn.DeclareQueue(context.Background(), queueConf, exchangeConfig, "populate_configs")
	if err != nil {
		return err
	}
	consConfig, err := broker.NewConsumerConfig(fmt.Sprintf("logger_populate_%s", a.config.Consul.Id))
	if err != nil {
		return err
	}

	cons := broker.NewConsumer(a.rabbitConn, queueConf, consConfig, a.HandlePopulateConfigs, slogadapter.NewSlogLogger(slog.Default()))
	return cons.Start(context.Background())
}

// populate_configs
func (a *App) HandlePopulateConfigs(ctx context.Context, message amqp.Delivery) error {
	var event PopulateConfigEventRequest
	err := json.Unmarshal(message.Body, &event)
	if err != nil {
		return err
	}
	res, err := a.storage.Config().Select(ctx, &model.SearchOptions{Size: -1, Fields: []string{model.ConfigFields.Object, model.ConfigFields.Enabled}}, nil, &model.Filter{
		Column:         model.ConfigFields.DomainId,
		Value:          event.DomainId,
		ComparisonType: model.Equal,
	})
	if err != nil {
		return err
	}
	populatedConfigs := PopulateConfigEventResponse{
		DomainId: event.DomainId,
	}
	for _, re := range res {
		state := re.Enabled
		objName := re.Object.GetName()
		if objName == nil {
			// we don't know name of object, set enabled to true to not miss potential log
			slog.DebugContext(ctx, "object name not found in config", slog.Int("eventId", event.DomainId))
			continue
		}
		populatedConfigs.Configs = append(populatedConfigs.Configs, &Config{
			Name:  *objName,
			State: state,
		})
	}
	bytes, err := json.Marshal(&populatedConfigs)
	if err != nil {
		return err
	}
	err = a.brokerPublisher.Publish(ctx, "config_population", bytes, nil)
	if err != nil {
		return err
	}
	return nil
}
