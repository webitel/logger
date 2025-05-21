package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webitel/logger/internal/auth"
	authmodel "github.com/webitel/logger/internal/auth/model"
	"github.com/webitel/logger/internal/auth/webitel_manager"
	"github.com/webitel/logger/internal/storage"
	"github.com/webitel/logger/internal/storage/postgres"
	"github.com/webitel/logger/internal/watcher"
	broker "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq"
	slogadapter "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq/pkg/adapter/slog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
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
	SessionCacheSize    = 35000
	SessionCacheTime    = 60 * 5

	DefaultPageSize = 40
	DefaultPage     = 1
	MaxPageSize     = 40000
)

type App struct {
	config          *model.AppConfig
	tracer          trace.Tracer
	storage         storage.Storage
	file            storagegrpc.FileServiceClient
	logUploaders    map[string]*watcher.UploadWatcher
	logCleaners     map[string]*watcher.Watcher
	serverExitChan  chan model.AppError
	rabbitConn      *broker.Connection
	rabbitConsumers []broker.Consumer

	server         *AppServer
	storageConn    *grpc.ClientConn
	sessionManager auth.AuthManager
	webitelAppConn *grpc.ClientConn
}

func (app *App) Database() storage.Storage {
	return app.storage
}

func New(config *model.AppConfig) (*App, error) {
	app := &App{config: config, serverExitChan: make(chan model.AppError), tracer: otel.GetTracerProvider().Tracer("app_internal")}
	var err error
	// init of database
	if &config.Database == nil {
		return nil, model.NewInternalError("app.app.new.database_config.bad_arguments", "error creating storage, config is nil")
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

	// init of grpc server
	s, appErr := BuildServer(app, app.config.Consul, app.serverExitChan)
	if appErr != nil {
		return nil, appErr
	}
	app.server = s

	// init service connections
	app.storageConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/storage?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, model.NewInternalError("app.app.new_app.grpc_conn.error", err.Error())
	}

	app.file = storagegrpc.NewFileServiceClient(app.storageConn)

	app.webitelAppConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/go.webitel.app?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, model.NewInternalError("app.app.new_app.grpc_conn.error", err.Error())
	}

	app.sessionManager, appErr = webitel_manager.NewWebitelAppAuthManager(app.webitelAppConn)
	if appErr != nil {
		return nil, appErr
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
	err = a.BrokerListenNewRecordLogs()
	if err != nil {
		return err
	}
	err = a.BrokerListenLoginAttempts()
	if err != nil {
		return err
	}
	// * run grpc server
	go a.server.Start()
	err = <-a.serverExitChan
	a.StopAllWatchers()
	a.rabbitConn.Close()
	err = a.storage.Close()
	if err != nil {
		slog.Error(err.Error())
	}

	return appErr
}

func (a *App) Stop() model.AppError {
	// close massive modules
	a.StopAllWatchers()
	//a.rabbit.Stop()
	a.server.Stop()

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

func (a *App) AuthorizeFromContext(ctx context.Context) (*authmodel.Session, error) {
	span := trace.SpanFromContext(ctx)
	session, err := a.sessionManager.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, model.NewUnauthorizedError("app.app.authorize_from_context.validate_session.expired", "session expired")
	}
	span.SetAttributes(attribute.Int64("caller_user.id", session.GetUserId()), attribute.Int64("caller_user.domain", session.GetDomainId()))
	return session, nil
}

func (a *App) MakePermissionError(session *authmodel.Session) model.AppError {
	if session == nil {
		return model.NewForbiddenError("app.permissions.check_access.denied", "access denied")
	}
	return model.NewForbiddenError("app.permissions.check_access.denied", fmt.Sprintf("userId=%d, access denied", session.GetUserId()))
}

func (a *App) MakeScopeError(session *authmodel.Session, scope string, access authmodel.AccessMode) model.AppError {
	if session == nil || session.GetUser() == nil || scope == "" {
		return model.NewForbiddenError("app.scope.check_access.denied", fmt.Sprintf("access denied"))
	}
	return model.NewForbiddenError("app.scope.check_access.denied", fmt.Sprintf("access denied scope=%s access=%d for user %d", scope, access, session.GetUserId()))
}

func (a *App) StopAllWatchers() {
	for _, cleaner := range a.logCleaners {
		cleaner.Stop()
	}
	for _, uploader := range a.logUploaders {
		uploader.Stop()
	}
}

func (a *App) BrokerListenNewRecordLogs() error {
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

func (a *App) BrokerListenLoginAttempts() error {
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
