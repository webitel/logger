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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
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
	tracer         trace.Tracer
	storage        storage.Storage
	file           storage_grpc.FileServiceClient
	logUploaders   map[string]*watcher.UploadWatcher
	logCleaners    map[string]*watcher.Watcher
	rabbitExitChan chan model.AppError
	serverExitChan chan model.AppError
	rabbit         *broker.RabbitBroker
	server         *AppServer
	storageConn    *grpc.ClientConn
	sessionManager auth.AuthManager
	webitelAppConn *grpc.ClientConn
}

func New(config *model.AppConfig) (*App, model.AppError) {
	app := &App{config: config, rabbitExitChan: make(chan model.AppError), serverExitChan: make(chan model.AppError), tracer: otel.GetTracerProvider().Tracer("app_internal")}
	var err error
	// init of database
	if &config.Database == nil {
		return nil, model.NewInternalError("app.app.new.database_config.bad_arguments", "error creating storage, config is nil")
	}
	app.storage = BuildDatabase(config.Database)

	// init of rabbit
	r, appErr := broker.BuildRabbit(app.config.Rabbit, app.rabbitExitChan)
	if appErr != nil {
		return nil, appErr
	}
	app.rabbit = r

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
	appErr = a.rabbit.Start()
	if appErr != nil {
		return appErr
	}
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
	select {
	case appErr = <-a.rabbitExitChan:
		a.server.Stop()
	case appErr = <-a.serverExitChan:
		a.rabbit.Stop()
	}
	a.StopAllWatchers()
	appErr = a.storage.Close()
	if appErr != nil {
		slog.Error(appErr.Error())
	}

	return appErr
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

	formatLog := func(s string) string {
		return fmt.Sprintf("[broker.record_logs.listener]: %s", s)
	}
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
	handler := func(timeout time.Duration, del <-chan amqp.Delivery, stopper chan any) {
		var (
			appErr model.AppError
			err    error
		)
		for {
			select {
			case message, closed := <-del:
				logAttr := slog.Group("message", slog.String("routing", message.RoutingKey), slog.String("body", string(message.Body)))
				//logAttr := []any{"routing", message.RoutingKey, "body", string(message.Body)}
				if !closed {
					slog.Warn(formatLog("channel closed, ending listening"), logAttr)
					return
				}
				// adding timeout on each handle
				ctx, cancelContext := context.WithTimeout(context.Background(), timeout)

				// try to handle the message
				appErr = a.HandleRabbitRecordLogMessage(ctx, &message)
				cancelContext()
				if appErr != nil {
					if errors.Is(appErr, sql.ErrNoRows) { // TODO: real foreign key check by postgres model
						message.Nack(false, false)
					} else {
						message.Nack(false, true)
					}
					slog.Warn(formatLog(appErr.Error()), logAttr)
					continue
				}
				err = message.Ack(false)
				if err != nil {
					appErr = model.NewInternalError("rabbit.listener.listen.acknowledge.fail", err.Error())
					err = nil
				}
				if appErr != nil {
					slog.Warn(appErr.Error(), logAttr)
				}
				slog.Info(formatLog("message processed"), logAttr)
			case <-stopper:
				slog.Info(formatLog(fmt.Sprintf("consuming stopped")))
				return
			}
		}

	}
	err = a.rabbit.QueueStartConsume(queueName, "", handler, 5*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) BrokerListenLoginAttempts() model.AppError {
	formatLog := func(s string) string {
		return fmt.Sprintf("[broker.login_attempts.listener]: %s", s)
	}
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

	handler := func(timeout time.Duration, del <-chan amqp.Delivery, stopper chan any) {
		var (
			//message amqp.Delivery
			appErr model.AppError
		)
		for {
			select {
			case message, closed := <-del:

				logAttr := slog.Group("message", slog.String("routing", message.RoutingKey), slog.String("body", string(message.Body)))
				if !closed {
					slog.Warn(formatLog("channel closed, ending listening"), logAttr)
					return
				}
				// adding timeout on each handle
				ctx, cancelContext := context.WithTimeout(context.Background(), timeout)

				// try to handle the message
				appErr = a.HandleRabbitLoginMessage(ctx, &message)
				cancelContext()
				if appErr != nil {
					slog.Warn(formatLog(appErr.Error()), logAttr)
				}
				slog.Info(formatLog("message processed"), logAttr)
			case <-stopper:
				slog.Info(formatLog(fmt.Sprintf("consuming stopped")))
				return
			}
		}
	}
	err = a.rabbit.QueueStartConsume(queueName, "", handler, 5*time.Second)
	if err != nil {
		return err
	}
	return nil
}

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
		return model.NewInternalError("app.app.handle_rabbit_login_message.key.error", "provided routing key is not matching with this handler")
	}

	databaseModel, appErr := m.ConvertToDatabaseModel()
	if appErr != nil {
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

func GroupAttributesAndBindToSpan(ctx context.Context, rootName string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	group := model.NewAttributeGroup(rootName, attrs...)
	span.SetAttributes(group.Unparse()...)
}

func GroupIncomingAttributesAndBindToSpan(ctx context.Context, attrs ...attribute.KeyValue) {
	GroupAttributesAndBindToSpan(ctx, "in", attrs...)
}

func GroupOutgoingAttributesAndBindToSpan(ctx context.Context, attrs ...attribute.KeyValue) {
	GroupAttributesAndBindToSpan(ctx, "out", attrs...)
}
