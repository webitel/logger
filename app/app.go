package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/webitel/logger/auth"
	authmodel "github.com/webitel/logger/auth/model"
	"github.com/webitel/logger/auth/webitel_manager"
	"strings"

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
	rabbit         *RabbitListener
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

	// init of rabbit1
	r, appErr := BuildRabbit(app, app.config.Rabbit, app.exitChan)
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

	app.storageConn, err = grpc.Dial(fmt.Sprintf("consul://%s/storage?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, model.NewInternalError("app.app.new_app.grpc_conn.error", err.Error())
	}

	app.file = storage_grpc.NewFileServiceClient(app.storageConn)

	app.webitelAppConn, err = grpc.Dial(fmt.Sprintf("consul://%s/go.webitel.app?wait=14s", config.Consul.Address),
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
	// * run grpc server
	go a.server.Start()
	//go ServeRequests(a, a.config.Consul, a.exitChan)
	return <-a.exitChan
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
