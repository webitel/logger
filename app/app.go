package app

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/pkg/cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/webitel/engine/auth_manager"
	"github.com/webitel/engine/discovery"
	strg "github.com/webitel/logger/api/storage"
	"github.com/webitel/logger/storage"
	"github.com/webitel/logger/storage/postgres"
	"github.com/webitel/logger/watcher"
	"google.golang.org/grpc/metadata"

	errors "github.com/webitel/engine/model"
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
	config           *model.AppConfig
	storage          storage.Storage
	file             strg.FileServiceClient
	uploadWatchers   map[string]*watcher.UploadWatcher
	deleteWatchers   map[string]*watcher.Watcher
	serviceDiscovery discovery.ServiceDiscovery
	sessionManager   auth_manager.AuthManager
	cache            cache.CacheStore
	exitChan         chan errors.AppError
	rabbit           *RabbitListener
	server           *AppServer
	grpcConn         *grpc.ClientConn
}

func New(config *model.AppConfig) (*App, errors.AppError) {
	app := &App{config: config, exitChan: make(chan errors.AppError)}

	// service registration
	disc, err := discovery.NewServiceDiscovery(config.Consul.Id, config.Consul.Address, func() (bool, error) {
		return true, nil
	})

	if err != nil {
		return nil, errors.NewInternalError("app.app.new.discovery_connection.fail", err.Error())
	}
	app.serviceDiscovery = disc
	// init of auth manager
	app.sessionManager = auth_manager.NewAuthManager(SessionCacheSize, SessionCacheTime, disc)
	if err := app.sessionManager.Start(); err != nil {
		return nil, errors.NewInternalError("app.app.new.auth_manager_start.fail", err.Error())
	}
	// init of cache storage
	app.cache = cache.NewMemoryCache(&cache.MemoryCacheConfig{
		Size:          200000,
		DefaultExpiry: 120,
	})

	// init of database
	if config.Database == nil {
		errors.NewInternalError("app.app.new.database_config.bad_arguments", "error creating storage, config is nil")
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

	app.grpcConn, err = grpc.Dial(fmt.Sprintf("consul://%s/storage?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.NewInternalError("app.app.new_app.grpc_conn.error", err.Error())
	}

	app.file = strg.NewFileServiceClient(app.grpcConn)
	return app, nil
}

func (a *App) GetConfig() *model.AppConfig {
	return a.config
}

func IsErrNoRows(err errors.AppError) bool {
	return strings.Contains(err.Error(), sql.ErrNoRows.Error())
}

func (a *App) Start() errors.AppError {

	err := a.storage.Open()
	if err != nil {
		return err
	}

	appErr := a.initializeWatchers()
	if appErr != nil {
		return appErr
	}
	// * Build and run rabbit1 listener
	go a.rabbit.Start()
	// * Build and run grpc server
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

func (a *App) GetSessionFromCtx(ctx context.Context) (*auth_manager.Session, errors.AppError) {
	var session *auth_manager.Session
	var err errors.AppError
	var token []string
	var info metadata.MD
	var ok bool

	v := ctx.Value(RequestContextName)
	info, ok = v.(metadata.MD)

	// todo
	if !ok {
		info, ok = metadata.FromIncomingContext(ctx)
	}

	if !ok {
		return nil, errors.NewForbiddenError("app.grpc.get_context", "Not found")
	} else {
		token = info.Get("X-Webitel-Access")
	}

	if len(token) < 1 {
		return nil, errors.NewInternalError("context.session_expired.app_error", "token not found")
	}

	session, err = a.GetSession(token[0])
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		return nil, errors.NewForbiddenError("context.session_expired.app_error", "token="+token[0])
	}

	return session, nil
}

func (a *App) MakePermissionError(session *auth_manager.Session, permission auth_manager.SessionPermission, access auth_manager.PermissionAccess) errors.AppError {

	return errors.NewForbiddenError("context.permissions.app_error", fmt.Sprintf("userId=%d, permission=%s access=%s", session.UserId, permission.Name, access.Name()))
}

func (a *App) Stop() errors.AppError {
	if a.serviceDiscovery != nil {
		a.serviceDiscovery.Shutdown()
	}
	a.storage.Close()
	a.rabbit.Stop()
	a.grpcConn.Close()

	return nil
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
func ConvertToOutputBulk[C any, K any](items []C, convertFunc func(C) (K, errors.AppError)) ([]K, errors.AppError) {
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
func CalculateListResultMetadata[C any, K any](s Lister, items []C, convertFunc func(C) (K, errors.AppError)) (bool, []K, errors.AppError) {
	var (
		result []K
		err    errors.AppError
	)
	next, filteredInput := GetListResult[C](s, items)
	result, err = ConvertToOutputBulk[C, K](filteredInput, convertFunc)
	if err != nil {
		return false, nil, err
	}
	return next, result, nil
}
