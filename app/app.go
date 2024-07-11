package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webitel/logger/auth"
	authmodel "github.com/webitel/logger/auth/model"
	"github.com/webitel/logger/auth/webitel_manager"
	broker "github.com/webitel/logger/broker/rabbit"
	"github.com/webitel/wlog"
	"strconv"
	"strings"
	"time"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/webitel/logger/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	storage_grpc "buf.build/gen/go/webitel/storage/grpc/go/_gogrpc"
	"github.com/webitel/logger/storage"
	"github.com/webitel/logger/storage/postgres"
	"github.com/webitel/logger/watcher"
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
	config         *model.AppConfig
	storage        storage.Storage
	file           storage_grpc.FileServiceClient
	logUploaders   map[string]*watcher.UploadWatcher
	logCleaners    map[string]*watcher.Watcher
	exitChan       chan model.AppError
	rabbit         *broker.RabbitBroker
	server         *AppServer
	storageConn    *grpc.ClientConn
	sessionManager auth.AuthManager
	webitelAppConn *grpc.ClientConn
}

func New(config *model.AppConfig) (*App, model.AppError) {
	app := &App{config: config, exitChan: make(chan model.AppError)}
	var err error

	// init of database
	if config.Database == nil {
		model.NewInternalError("app.app.new.database_config.bad_arguments", "error creating storage, config is nil")
	}
	app.storage = BuildDatabase(config.Database)

	// init of rabbit
	r, appErr := broker.BuildRabbit(app.config.Rabbit, app.exitChan)
	if appErr != nil {
		return nil, appErr
	}
	app.rabbit = r

	// init of grpc server
	s, appErr := BuildServer(app, app.config.Consul, app.exitChan)
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

	app.file = storage_grpc.NewFileServiceClient(app.storageConn)

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

func IsErrNoRows(err model.AppError) bool {
	return strings.Contains(err.Error(), sql.ErrNoRows.Error())
}

func (a *App) Start() model.AppError {

	err := a.storage.Open()
	if err != nil {
		return err
	}

	appErr := a.initializeWatchers()
	if appErr != nil {
		return appErr
	}
	// * run rabbit listener
	a.rabbit.Start()
	appErr = a.BrokerListenNewRecordLogs()
	if appErr != nil {
		return appErr
	}
	appErr = a.BrokerListenLoginAttempts()
	if appErr != nil {
		return appErr
	}
	// * run grpc server
	go a.server.Start()
	//go ServeRequests(a, a.config.Consul, a.exitChan)
	return <-a.exitChan
}

func (a *App) Stop() model.AppError {
	// close massive modules
	a.StopAllWatchers()
	a.rabbit.Stop()
	a.server.Stop()

	// close db connection
	a.storage.Close()

	// close grpc connections
	a.storageConn.Close()
	a.webitelAppConn.Close()

	return nil
}

type Searcher interface {
	GetPage() int32
	GetSize() int32
	GetQ() string
	GetSort() string
	GetFields() []string
}

func ExtractSearchOptions(t Searcher) *model.SearchOptions {
	var res model.SearchOptions
	if t.GetSort() != "" {
		res.Sort = ConvertSort(t.GetSort())
	}
	if t.GetSize() <= 0 || t.GetSize() > MaxPageSize {
		res.Size = DefaultPageSize
	} else {
		res.Size = int(t.GetSize())
	}
	if t.GetPage() <= 0 {
		res.Page = DefaultPage
	} else {
		res.Page = int(t.GetPage())
	}
	if t.GetQ() != "" {
		//	if input := strings.Replace(t.GetQ(), "*", "%", -1); input == "" {
		res.Search = strings.Replace(t.GetQ(), "*", "%", -1)
		//	}

	}
	if s := t.GetFields(); len(s) != 0 {
		res.Fields = s
	}
	return &res
}

func BuildDatabase(config *model.DatabaseConfig) storage.Storage {
	return postgres.New(config)
}

func (a *App) AuthorizeFromContext(ctx context.Context) (*authmodel.Session, model.AppError) {
	session, err := a.sessionManager.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, model.NewUnauthorizedError("app.app.authorize_from_context.validate_session.expired", "session expired")
	}
	return session, nil
}

func (a *App) MakePermissionError(session *authmodel.Session) model.AppError {
	if session == nil {
		return model.NewForbiddenError("app.permissions.check_access.denied", "access denied")
	}
	return model.NewForbiddenError("app.permissions.check_access.denied", fmt.Sprintf("userId=%d, access denied", session.GetUserId()))
}

func (a *App) MakeScopeError(session *authmodel.Session, scope *authmodel.Scope, access authmodel.AccessMode) model.AppError {
	if session == nil || session.GetUser() == nil || scope == nil {
		return model.NewForbiddenError("app.scope.check_access.denied", fmt.Sprintf("access denied"))
	}
	return model.NewForbiddenError("app.scope.check_access.denied", fmt.Sprintf("access denied scope=%s access=%d for user %d", scope.Name, access, session.GetUserId()))
}

func (a *App) StopAllWatchers() {
	for _, cleaner := range a.logCleaners {
		cleaner.Stop()
	}
	for _, uploader := range a.logUploaders {
		uploader.Stop()
	}
}

func ConvertSort(in string) string {
	if len(in) < 2 || (in[0] != '+' && in[0] != '-') {
		return ""
	}
	if in[0] == '+' {
		return fmt.Sprintf("%s:%s", "ASC", in[1:])
	} else {
		return fmt.Sprintf("%s:%s", "DESC", in[1:])
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

type Lister interface {
	GetSize() int32
}

// C type of items to filter
func GetListResult[C any](s Lister, items []C) (bool, []C) {
	if int32(len(items)-1) == s.GetSize() {
		return true, items[0 : len(items)-1]
	}
	return false, items
}

// C type of input, K type of output
func ConvertToOutputBulk[C any, K any](items []C, convertFunc func(C) (K, model.AppError)) ([]K, model.AppError) {
	var result []K
	for _, item := range items {
		out, err := convertFunc(item)
		if err != nil {
			return nil, err
		}
		result = append(result, out)
	}
	return result, nil
}

// C type of input, K type of output
func CalculateListResultMetadata[C any, K any](s Lister, items []C, convertFunc func(C) (K, model.AppError)) (bool, []K, model.AppError) {
	var (
		result []K
		err    model.AppError
	)
	next, filteredInput := GetListResult[C](s, items)
	result, err = ConvertToOutputBulk[C, K](filteredInput, convertFunc)
	if err != nil {
		return false, nil, err
	}
	return next, result, nil
}

func (a *App) BrokerListenNewRecordLogs() model.AppError {
	// create or connect the logger exchange
	appErr := a.rabbit.ExchangeDeclare("logger", "topic", broker.ExchangeEnableDurable)
	if appErr != nil {
		return appErr
	}
	// bind all logs from webitel exchange to the logger exchange
	appErr = a.rabbit.ExchangeBind("logger", "logger.#", "webitel", true, nil)
	if appErr != nil {
		return appErr
	}
	// declare new queue logger.service
	queueName, appErr := a.rabbit.QueueDeclare("logger.service", broker.QueueEnableDurable)
	if appErr != nil {
		return appErr
	}
	err := a.rabbit.QueueBind(
		"logger",
		queueName,
		"logger.#",
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = a.rabbit.QueueStartConsume(queueName, "", LogRecordAcknowledger, a.HandleRabbitRecordLogMessage, 5*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) BrokerListenLoginAttempts() model.AppError {
	sourceExchangeName := "webitel"
	sourceExchangeKind := "topic"
	// create or connect the logger exchange
	appErr := a.rabbit.ExchangeDeclare(sourceExchangeName, sourceExchangeKind, broker.ExchangeEnableDurable)
	if appErr != nil {
		return appErr
	}
	// declare new queue logger.service
	queueName, appErr := a.rabbit.QueueDeclare("logger.login", broker.QueueEnableDurable)
	if appErr != nil {
		return appErr
	}
	err := a.rabbit.QueueBind(
		sourceExchangeName,
		queueName,
		"login.#",
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = a.rabbit.QueueStartConsume(queueName, "", DefaultAcknowledger, a.HandleRabbitLoginMessage, 5*time.Second)
	if err != nil {
		return err
	}
	return nil
}

var (
	DefaultAcknowledger = func(timeout time.Duration, del <-chan amqp.Delivery, stopper chan any, handleFunc broker.HandleFunc) {
		var (
			//message amqp.Delivery
			appErr model.AppError
		)
		for {
			select {
			case message, closed := <-del:
				if !closed {
					wlog.Debug(fmt.Sprintf("[broker.handler]: channel closed"))
					return
				}
				// adding timeout on each handle
				ctx, cancelContext := context.WithTimeout(context.Background(), timeout)

				// try to handle the message
				appErr = handleFunc(ctx, &message)
				cancelContext()
				if appErr != nil {
					wlog.Debug(fmt.Sprintf("[broker.handler]: (%s) %s", message.RoutingKey, appErr.Error()))
				}
				wlog.Debug(fmt.Sprintf("[broker.handler]: (%s) processed", message.RoutingKey))
			case <-stopper:
				return
			}
		}
	}
	LogRecordAcknowledger = func(timeout time.Duration, del <-chan amqp.Delivery, stopper chan any, handleFunc broker.HandleFunc) {
		var (
			appErr model.AppError
			err    error
		)
		for {
			select {
			case message, closed := <-del:

				if !closed {
					wlog.Debug(fmt.Sprintf("[broker.handler]: channel closed"))
					return
				}
				// adding timeout on each handle
				ctx, cancelContext := context.WithTimeout(context.Background(), timeout)

				// try to handle the message
				appErr = handleFunc(ctx, &message)
				cancelContext()
				if appErr != nil {
					if errors.Is(appErr, sql.ErrNoRows) { // TODO: real foreign key check by postgres model
						message.Nack(false, false)
					} else {
						message.Nack(false, true)
					}
					continue
				}
				err = message.Ack(false)
				if err != nil {
					appErr = model.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
				}
				if appErr != nil {
					wlog.Debug(fmt.Sprintf("[broker.handler]: %s", appErr.Error()))
				}
			case <-stopper:
				return
			}
		}

	}
)

func (a *App) HandleRabbitLoginMessage(ctx context.Context, message *amqp.Delivery) model.AppError {
	var (
		m model.BrokerLoginMessage
	)
	err := json.Unmarshal(message.Body, &m)
	if err != nil {
		message.Nack(false, false)
		return model.NewInternalError("app.app.handle_rabbit_login_message.unmarshal.error", err.Error())
	}

	splittedKey := strings.Split(message.RoutingKey, ".")
	if len(splittedKey) < 4 {
		message.Nack(false, true)
		return model.NewInternalError("app.app.handle_rabbit_login_message.unmarshal.error", err.Error())
	}

	databaseModel, appErr := m.ConvertToDatabaseModel()
	if err != nil {
		message.Nack(false, false)
		return appErr
	}

	_, err = a.storage.LoginAttempt().Insert(ctx, databaseModel)
	if err != nil {
		message.Nack(false, false)
		return model.NewInternalError("app.app.handle_rabbit_login_message.insert.error", err.Error())
	}
	message.Ack(false)
	return nil
}

func (a *App) HandleRabbitRecordLogMessage(ctx context.Context, message *amqp.Delivery) model.AppError {
	var (
		m      model.BrokerRecordLogMessage
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
		appErr := a.InsertLogByRabbitMessageBulk(ctx, rabbitMessages, domain)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}
