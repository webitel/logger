package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/webitel/logger/internal/auth"
	"github.com/webitel/logger/internal/watcher"
	notifier "github.com/webitel/webitel-go-kit/pkg/watcher"
	"log/slog"

	proto "github.com/webitel/logger/api/logger"
	"github.com/webitel/logger/internal/model"
)

const (
	ConfigNotifierObject = "config"
)

func (a *App) UpdateConfig(ctx context.Context, in *model.Config) (*model.Config, error) {
	var (
		err       error
		newConfig *model.Config
	)
	defer func() {
		notifyErr := a.watcherManager.Notify(ConfigNotifierObject, notifier.EventTypeUpdate, NewNotifierConfigArgs(err == nil, newConfig))
		if notifyErr != nil {
			slog.ErrorContext(ctx, notifyErr.Error())
		}
	}()
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Edit)
	if err != nil {
		return nil, err
	}
	// OBAC check
	if !session.CheckObacAccess() {
		err = a.MakeScopeError(session.GetMainObjClassName())
		return nil, err
	}
	if in == nil {
		err = errors.New("config is nil")
		return nil, err
	}
	// RBAC check
	if session.IsRbacCheckRequired() {
		access, err := a.ConfigCheckAccess(ctx, in.Id, session.GetUserId(), session.GetRoles(), session.GetMainAccessMode())
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, a.MakeScopeError(session.GetMainObjClassName())
		}
	}

	oldConfig, err := a.storage.Config().GetById(ctx, nil, in.Id, session.GetDomainId())
	if err != nil {
		return nil, err
	}
	in.NextUploadOn = calculateNextPeriodFromNow(int32(in.Period))
	newConfig, err = a.storage.Config().Update(ctx, in, []string{}, session.GetUserId())
	if err != nil {
		return nil, err
	}
	a.UpdateConfigWatchers(ctx, oldConfig, newConfig)
	return newConfig, nil

}

func (a *App) GetSystemObjects(ctx context.Context, in *proto.ReadSystemObjectsRequest) ([]*model.SystemObject, error) {
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Read)
	if err != nil {
		return nil, err
	}
	// OBAC check
	if !session.CheckObacAccess() {
		return nil, a.MakeScopeError(session.GetMainObjClassName())
	}
	var filters []string
	for _, name := range proto.AvailableSystemObjects_name {
		filters = append(filters, name)
	}
	return a.storage.Config().GetAvailableSystemObjects(ctx, session.GetDomainId(), in.GetIncludeExisting(),
		filters...)

}

func (a *App) UpdateConfigWatchers(ctx context.Context, oldConfig, newConfig *model.Config) {
	configId := newConfig.Id
	domainId := newConfig.DomainId
	if !newConfig.Enabled && oldConfig.Enabled { // changed to disabled
		a.DeleteWatchers(configId)
		slog.InfoContext(ctx, fmt.Sprintf("config with id %d disabled... watchers have been deleted !", configId))
	} else if newConfig.Enabled && oldConfig.Enabled || (newConfig.Enabled && !oldConfig.Enabled) { // status wasn't changed and it's still enabled OR changed to enabled
		// if days to store changed, then we need to update cleaner worker
		if newConfig.DaysToStore != oldConfig.DaysToStore {
			// if cleaned exists update, otherwise insert new
			if a.GetLogCleaner(configId) != nil {
				a.UpdateLogCleanerWithNewInterval(ctx, configId, newConfig.DaysToStore)
			} else {
				a.InsertLogCleaner(configId, nil, newConfig.DaysToStore)
			}
		}
		// find uploader params, if exists then update params otherwise insert new
		if params := a.GetLogUploaderParams(configId); params != nil {
			params.UserId = newConfig.Editor.GetId()
			params.StorageId = *newConfig.Storage.Id
			params.Period = newConfig.Period
		} else {
			a.InsertLogUploader(configId, nil, &watcher.UploadWatcherParams{
				StorageId:    *newConfig.Storage.Id,
				Period:       newConfig.Period,
				NextUploadOn: newConfig.NextUploadOn,
				LastLogId:    newConfig.LastUploadedLogId,
				DomainId:     domainId,
			})
		}
	}
	// else status still disabled
}

func (a *App) CreateConfig(ctx context.Context, in *model.Config) (*model.Config, error) {
	var (
		err      error
		newModel *model.Config
	)
	defer func() {
		notifyErr := a.watcherManager.Notify(ConfigNotifierObject, notifier.EventTypeCreate, NewNotifierConfigArgs(err == nil, newModel))
		if notifyErr != nil {
			slog.ErrorContext(ctx, notifyErr.Error())
		}
	}()
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Add)
	if err != nil {
		return nil, err
	}
	// OBAC check
	if !session.CheckObacAccess() {
		err = a.MakeScopeError(session.GetMainObjClassName())
		return nil, err
	}
	if in == nil {
		err = errors.New("config is nil")
		return nil, err
	}
	userId := session.GetUserId()
	domainId := session.GetDomainId()
	in.Author = &model.Author{Id: &userId}
	in.DomainId = domainId
	in.NextUploadOn = calculateNextPeriodFromNow(int32(in.Period))
	newModel, err = a.storage.Config().Insert(ctx, in)
	if err != nil {
		return nil, err
	}
	if newModel.Enabled {
		a.InsertLogCleaner(newModel.Id, nil, newModel.DaysToStore)
		a.InsertLogUploader(newModel.Id, nil, &watcher.UploadWatcherParams{
			StorageId:    *newModel.Storage.Id,
			Period:       newModel.Period,
			NextUploadOn: newModel.NextUploadOn,
			LastLogId:    nil,
			DomainId:     in.DomainId,
		})
	}

	return newModel, nil
}

func (a *App) GetConfigByObjectId(ctx context.Context, objectId int, domainId int) (*model.Config, error) {
	res, appErr := a.storage.Config().GetByObjectId(ctx, domainId, objectId)
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (a *App) CheckConfigStatus(ctx context.Context, objectName string, domainId int) (bool, error) {
	searchResult, err := a.storage.Config().Select(ctx, nil, nil, &model.FilterNode{
		Nodes: []any{
			&model.Filter{
				Column:         "wbt_class.name",
				Value:          objectName,
				ComparisonType: model.ILike,
			},
			&model.Filter{
				Column:         model.ConfigFields.DomainId,
				Value:          domainId,
				ComparisonType: model.Equal,
			}},
		Connection: model.AND,
	})
	if err != nil {
		return false, err
	}
	if len(searchResult) == 0 {
		return false, nil
	}
	enabled := searchResult[0].Enabled
	return enabled, nil

}

func (a *App) GetConfigById(ctx context.Context, rbac *model.RbacOptions, id int) (*model.Config, error) {
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Read)
	if err != nil {
		return nil, err
	}
	// OBAC check
	if !session.CheckObacAccess() {
		return nil, a.MakeScopeError(session.GetMainObjClassName())
	}
	// RBAC check
	if session.IsRbacCheckRequired() {
		rbac = &model.RbacOptions{
			Groups: session.GetRoles(),
			Access: session.GetMainAccessMode().Value(),
		}
	}
	res, appErr := a.storage.Config().GetById(ctx, rbac, id, session.GetDomainId())
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (a *App) DeleteConfig(ctx context.Context, ids []int) error {
	var err error
	defer func() {
		notifyErr := a.watcherManager.Notify(ConfigNotifierObject, notifier.EventTypeDelete, NewNotifierConfigArgs(err == nil, nil))
		if notifyErr != nil {
			slog.ErrorContext(ctx, notifyErr.Error())
		}
	}()
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Delete)
	if err != nil {
		return err
	}
	if !session.CheckObacAccess() {
		err = a.MakeScopeError(session.GetMainObjClassName())
		return err
	}
	var rbac *model.RbacOptions
	if session.IsRbacCheckRequired() {
		rbac = &model.RbacOptions{
			Groups: session.GetRoles(),
			Access: session.GetMainAccessMode().Value(),
		}
	}
	_, err = a.storage.Config().DeleteMany(ctx, rbac, ids, session.GetDomainId())
	if err != nil {
		return err
	}
	return nil
}

func (a *App) ConfigCheckAccess(ctx context.Context, id int, domainId int, groups []int, access auth.AccessMode) (bool, error) {
	return a.storage.Config().CheckAccess(ctx, domainId, id, groups, access.Value())

}

func (a *App) SearchConfig(ctx context.Context, rbac *model.RbacOptions, searchOpt *model.SearchOptions) ([]*model.Config, error) {
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Read)
	if err != nil {
		return nil, err
	}
	// OBAC check
	if !session.CheckObacAccess() {
		return nil, a.MakeScopeError(session.GetMainObjClassName())
	}
	// RBAC check
	if session.IsRbacCheckRequired() {
		rbac = &model.RbacOptions{
			Groups: session.GetRoles(),
			Access: session.GetMainAccessMode().Value(),
		}
	}
	modelConfigs, err := a.storage.Config().Select(
		ctx,
		searchOpt,
		rbac,
		&model.Filter{
			Column:         "object_config.domain_id",
			Value:          session.GetDomainId(),
			ComparisonType: model.Equal,
		},
	)
	if err != nil {
		return nil, err
	}

	return modelConfigs, nil

}

type NotifierConfigArgs struct {
	Config             *model.Config
	OperationSucceeded bool
}

func (n *NotifierConfigArgs) GetArgs() map[string]any {
	return map[string]any{
		"object":    n.Config,
		"objclass":  ConfigNotifierObject,
		"succeeded": n.OperationSucceeded,
	}
}

func NewNotifierConfigArgs(success bool, config *model.Config) *NotifierConfigArgs {
	return &NotifierConfigArgs{Config: config, OperationSucceeded: success}
}
