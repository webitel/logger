package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/webitel/engine/auth_manager"
	"github.com/webitel/engine/discovery"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/pkg/cache"
	"github.com/webitel/logger/storage"
	"github.com/webitel/logger/watcher"
	"google.golang.org/grpc/metadata"
	"strings"

	errors "github.com/webitel/engine/model"
)

const (
	DeleteWatcherPrefix = "config.watcher"
	SESSION_CACHE_SIZE  = 35000
	SESSION_CACHE_TIME  = 60 * 5
	RequestContextName  = "grpc_ctx"

	DEFAULT_PAGE_SIZE = 40
	DEFAULT_PAGE      = 1
	MAX_PAGE_SIZE     = 40000
)

type App struct {
	config           *model.AppConfig
	storage          storage.Storage
	watchers         map[string]*watcher.Watcher
	serviceDiscovery discovery.ServiceDiscovery
	sessionManager   auth_manager.AuthManager
	cache            cache.CacheStore
}

func New(store storage.Storage, config *model.AppConfig) (*App, errors.AppError) {
	if store == nil {
		return nil, errors.NewInternalError("app.app.new.check_arguments.fail", "store is nil")
	}
	app := &App{storage: store, config: config}

	appErr := app.initializeWatchers()
	if appErr != nil {
		return nil, appErr
	}
	disc, err := discovery.NewServiceDiscovery(config.Consul.Id, config.Consul.Address, func() (bool, error) {
		return true, nil
	})
	if err != nil {
		return nil, errors.NewInternalError("app.app.new.discovery_connection.fail", err.Error())
	}
	app.sessionManager = auth_manager.NewAuthManager(SESSION_CACHE_SIZE, SESSION_CACHE_TIME, disc)
	if err := app.sessionManager.Start(); err != nil {
		return nil, errors.NewInternalError("app.app.new.auth_manager_start.fail", err.Error())
	}
	app.cache = cache.NewMemoryCache(&cache.MemoryCacheConfig{
		Size:          200000,
		DefaultExpiry: 120,
	})
	return app, nil
}

func (a *App) GetConfig() *model.AppConfig {
	return a.config
}

func IsErrNoRows(err errors.AppError) bool {
	return strings.Contains(err.Error(), sql.ErrNoRows.Error())
}

type Search interface {
	GetPage() int32
	GetSize() int32
	GetQ() string
	GetSort() string
	GetFields() []string
}

func ExtractSearchOptions(t Search) *model.SearchOptions {
	var res model.SearchOptions
	if t.GetSort() != "" {
		res.Sort = ConvertSort(res.Sort)
	}
	if t.GetSize() <= 0 || t.GetSize() > MAX_PAGE_SIZE {
		res.Size = DEFAULT_PAGE_SIZE
	} else {
		res.Size = int(t.GetSize())
	}
	if t.GetPage() <= 0 {
		res.Page = DEFAULT_PAGE
	} else {
		res.Page = int(t.GetPage())
	}
	if t.GetQ() != "" {
		res.Search = strings.Replace(res.Search, "*", "%", -1)
	}
	if s := t.GetFields(); len(s) != 0 {
		res.Fields = s
	}
	return &res
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
		return nil, errors.NewInternalError("api.context.session_expired.app_error", "token not found")
	}

	session, err = a.GetSession(token[0])
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		return nil, errors.NewForbiddenError("api.context.session_expired.app_error", "token="+token[0])
	}

	return session, nil
}

func (a *App) MakePermissionError(session *auth_manager.Session, permission auth_manager.SessionPermission, access auth_manager.PermissionAccess) errors.AppError {

	return errors.NewForbiddenError("api.context.permissions.app_error", fmt.Sprintf("userId=%d, permission=%s access=%s", session.UserId, permission.Name, access.Name()))
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
