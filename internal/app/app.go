package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webitel/logger/internal/auth"
	autherror "github.com/webitel/logger/internal/auth/errors"
	"github.com/webitel/logger/internal/auth/manager/webitel_app"
	handlergrpc "github.com/webitel/logger/internal/handler/grpc"
	"github.com/webitel/logger/internal/model"
	"github.com/webitel/logger/internal/registry"
	"github.com/webitel/logger/internal/registry/consul"
	"github.com/webitel/logger/internal/storage"
	"github.com/webitel/logger/internal/storage/postgres"
	"github.com/webitel/logger/internal/watcher"
	broker "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq"
	slogadapter "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq/pkg/adapter/slog"
	notifier "github.com/webitel/webitel-go-kit/pkg/watcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	storagegrpc "github.com/webitel/logger/api/storage"
)

const (
	DeleteWatcherPrefix = "config.delete.watcher"
	UploadWatcherPrefix = "config.upload.watcher"
)

type App struct {
	config         *model.AppConfig
	storage        storage.Storage
	file           storagegrpc.FileServiceClient
	grpcServer     *grpc.Server
	registry       registry.ServiceRegistrar
	sessionManager auth.Manager
	watcherManager notifier.Manager

	logUploaders    map[string]*watcher.UploadWatcher
	logCleaners     map[string]*watcher.Watcher
	brokerPublisher broker.Publisher

	// active connections
	rabbitConn     *broker.Connection
	storageConn    *grpc.ClientConn
	webitelAppConn *grpc.ClientConn

	// emergency channel, error in this channel signals fatal application error
	emergencyStop chan error
}

func (app *App) Database() storage.Storage {
	return app.storage
}

func New(config *model.AppConfig) (*App, error) {
	app := &App{config: config, emergencyStop: make(chan error), watcherManager: notifier.NewDefaultWatcherManager(true)}
	var err error
	// init of database
	if config.Database == nil {
		return nil, errors.New("error creating storage, config is nil")
	}
	app.storage = BuildDatabase(config.Database)
	// registry
	app.registry, err = consul.New(config.GrpcAddr, config.Consul)
	if err != nil {
		return nil, err
	}

	// GRPC handlers initialization
	app.grpcServer, err = handlergrpc.Build(app)
	if err != nil {
		return nil, err
	}

	// init service connections
	app.storageConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/storage?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	app.file = storagegrpc.NewFileServiceClient(app.storageConn)

	app.webitelAppConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/go.webitel.app?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	app.sessionManager, err = webitel_app.New(app.webitelAppConn)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) GetConfig() *model.AppConfig {
	return a.config
}

func (a *App) Start() error {

	err := a.storage.Open()
	if err != nil {
		return err
	}

	err = a.initializeWatchers()
	if err != nil {
		return err
	}

	err = a.initBroker()
	if err != nil {
		return err
	}

	// config watcher
	if a.config.Features.EnableSagaEvents {
		configSagaWatcher := notifier.NewDefaultWatcher()
		sagaObserver := NewSagaObserver(a.brokerPublisher)
		configSagaWatcher.Attach(notifier.EventTypeCreate, sagaObserver)
		configSagaWatcher.Attach(notifier.EventTypeUpdate, sagaObserver)
		configSagaWatcher.Attach(notifier.EventTypeDelete, sagaObserver)
		a.watcherManager.AddWatcher(ConfigNotifierObject, configSagaWatcher)

		logSagaWatcher := notifier.NewDefaultWatcher()
		logSagaWatcher.Attach(notifier.EventTypeCreate, sagaObserver)
		logSagaWatcher.Attach(notifier.EventTypeUpdate, sagaObserver)
		logSagaWatcher.Attach(notifier.EventTypeDelete, sagaObserver)
		a.watcherManager.AddWatcher(LogsNotifierObject, logSagaWatcher)
	}

	// * run grpc server
	listener, err := net.Listen("tcp", a.config.GrpcAddr)
	if err != nil {
		return err
	}
	err = a.registry.Register()
	if err != nil {
		return err
	}
	defer func(listener net.Listener) {
		err = listener.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(listener)
	go func() {
		err = a.grpcServer.Serve(listener)
		if err != nil {
			a.emergencyStop <- err
		}
	}()
	err = <-a.emergencyStop
	_ = a.Stop()

	return err
}

func (a *App) Stop() error {
	// close massive modules
	a.StopAllWatchers()
	_ = a.rabbitConn.Close()
	a.grpcServer.GracefulStop()

	// close db connection
	_ = a.storage.Close()

	// close grpc connections
	_ = a.storageConn.Close()
	_ = a.webitelAppConn.Close()

	return nil
}

func BuildDatabase(config *model.DatabaseConfig) storage.Storage {
	return postgres.New(config)
}

func (a *App) AuthorizeFromContext(ctx context.Context, objclass string, accessMode auth.AccessMode) (auth.Auther, error) {
	return a.sessionManager.AuthorizeFromContext(ctx, objclass, accessMode)
}

func (a *App) MakeScopeError(objclass string) error {
	return autherror.NewForbiddenError("application.objclass.access.denied", fmt.Sprintf("%s: access denied", objclass))
}

func (a *App) StopAllWatchers() {
	for _, cleaner := range a.logCleaners {
		cleaner.Stop()
	}
	for _, uploader := range a.logUploaders {
		uploader.Stop()
	}
}

func (a *App) initBroker() error {
	var err error
	a.rabbitConn, err = broker.NewConnection(&broker.Config{
		URL:            a.config.Rabbit.Url,
		ConnectTimeout: 1 * time.Minute,
	}, slogadapter.NewSlogLogger(slog.Default()))
	if err != nil {
		return err
	}
	exchangeConf, err := broker.NewExchangeConfig("logger", broker.ExchangeTypeTopic)
	if err != nil {
		return err
	}
	err = a.initExchange(exchangeConf)
	if err != nil {
		return err
	}

	pubConfig, err := broker.NewPublisherConfig()
	if err != nil {
		return err
	}
	a.brokerPublisher, err = broker.NewPublisher(a.rabbitConn, exchangeConf, pubConfig, slogadapter.NewSlogLogger(slog.Default()))
	if err != nil {
		return err
	}

	err = a.initLogsConsumption(exchangeConf)
	if err != nil {
		return err
	}
	if a.config.Features.EnableLoginConsumption {
		err = a.initLoginConsumption()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initExchange(config *broker.ExchangeConfig) error {
	return a.rabbitConn.DeclareExchange(context.Background(), config)
}

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
	consumerConf, err := broker.NewConsumerConfig(a.config.Consul.Id)
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
		date := time.UnixMilli(m.Date)
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
	// create or connect the logger exchange
	exchangeConf, err := broker.NewExchangeConfig(sourceExchangeName, broker.ExchangeTypeTopic)
	if err != nil {
		return err
	}
	queueConf, err := broker.NewQueueConfig("logger.login", broker.WithQueueDurable(true))
	if err != nil {
		return err
	}
	consConfig, err := broker.NewConsumerConfig("logger.login")
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
