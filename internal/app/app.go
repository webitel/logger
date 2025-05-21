package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webitel/logger/internal/auth"
	autherror "github.com/webitel/logger/internal/auth/errors"
	"github.com/webitel/logger/internal/auth/manager/webitel_app"
	handlergrpc "github.com/webitel/logger/internal/handler/grpc"
	"github.com/webitel/logger/internal/registry"
	"github.com/webitel/logger/internal/registry/consul"
	"github.com/webitel/logger/internal/storage"
	"github.com/webitel/logger/internal/storage/postgres"
	"github.com/webitel/logger/internal/watcher"
	broker "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq"
	slogadapter "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq/pkg/adapter/slog"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/webitel/logger/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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

	logUploaders    map[string]*watcher.UploadWatcher
	logCleaners     map[string]*watcher.Watcher
	brokerConsumers map[string]broker.Consumer

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
	app := &App{config: config, emergencyStop: make(chan error)}
	var err error
	// init of database
	if &config.Database == nil {
		return nil, errors.New("error creating storage, config is nil")
	}
	app.storage = BuildDatabase(config.Database)

	// init of rabbit
	app.rabbitConn, err = broker.NewConnection(&broker.Config{
		URL:            config.Rabbit.Url,
		ConnectTimeout: 1 * time.Minute,
	}, slogadapter.NewSlogLogger(slog.Default()))
	if err != nil {
		return nil, err
	}

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

	appErr := a.initializeWatchers()
	if appErr != nil {
		return appErr
	}

	err = a.initLogsConsumption()
	if err != nil {
		return err
	}

	if a.config.Features.EnableLoginConsumption {
		err = a.initLoginConsumption()
		if err != nil {
			return err
		}
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
	defer listener.Close()
	go func() {
		err = a.grpcServer.Serve(listener)
		if err != nil {
			a.emergencyStop <- err
		}
	}()
	err = <-a.emergencyStop
	a.Stop()

	return appErr
}

func (a *App) Stop() error {
	// close massive modules
	a.StopAllWatchers()
	a.rabbitConn.Close()
	a.grpcServer.GracefulStop()

	// close db connection
	a.storage.Close()

	// close grpc connections
	a.storageConn.Close()
	a.webitelAppConn.Close()

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

func (a *App) initLogsConsumption() error {
	exchangeConf, err := broker.NewExchangeConfig("logger", broker.ExchangeTypeTopic)
	if err != nil {
		return err
	}
	err = a.rabbitConn.DeclareExchange(context.Background(), exchangeConf)
	if err != nil {
		return err
	}
	// bind all logs from webitel exchange to the logger exchange
	err = a.rabbitConn.BindExchange(context.Background(), "webitel", "logger", "logger.#", true, nil)
	if err != nil {
		return err
	}
	// declare new queue logger.service
	queueConfig, err := broker.NewQueueConfig("logger.service", broker.WithQueueDurable(true))
	if err != nil {
		return err
	}
	err = a.rabbitConn.DeclareQueue(context.Background(), queueConfig, exchangeConf, "logger.#")
	if err != nil {
		return err
	}
	consumerConf, err := broker.NewConsumerConfig(a.config.Consul.Id)
	if err != nil {
		return err
	}
	consumer := broker.NewConsumer(a.rabbitConn, queueConfig, consumerConf, a.HandleRabbitRecordLogMessage, slogadapter.NewSlogLogger(slog.Default()))
	return consumer.Start(context.Background())
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

	cons := broker.NewConsumer(a.rabbitConn, queueConf, consConfig, a.HandleRabbitLoginMessage, slogadapter.NewSlogLogger(slog.Default()))
	return cons.Start(context.Background())
}

func (a *App) HandleRabbitLoginMessage(ctx context.Context, message amqp.Delivery) error {
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

func (a *App) HandleRabbitRecordLogMessage(ctx context.Context, message amqp.Delivery) error {
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
	if m.Records != nil {
		var rabbitMessages []*model.RabbitMessage
		for _, v := range m.Records {
			rabbitMessage := &model.RabbitMessage{
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
		err = a.InsertLogByRabbitMessageBulk(ctx, rabbitMessages, domain)
		if err != nil {
			return err
		}
	}

	return nil
}
