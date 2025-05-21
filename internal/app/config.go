package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/webitel/logger/internal/auth"
	"github.com/webitel/logger/internal/watcher"
	"log/slog"

	proto "github.com/webitel/logger/api/logger"
	"github.com/webitel/logger/internal/model"
)

func (a *App) UpdateConfig(ctx context.Context, in *model.Config) (*model.Config, error) {
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Edit)
	if err != nil {
		return nil, err
	}

	// OBAC check
	if !session.CheckObacAccess() {
		return nil, a.MakeScopeError(session.GetMainObjClassName())
	}

	// RBAC check
	if session.IsRbacCheckRequired() {
		access, err := a.ConfigCheckAccess(ctx, int64(in.Id), session.GetUserId(), session.GetRoles(), session.GetMainAccessMode())
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, a.MakeScopeError(session.GetMainObjClassName())
		}
	}

	var (
		newConfig *model.Config
	)
	if in == nil {
		return nil, errors.New("config proto is nil")
	}

	oldConfig, err := a.storage.Config().GetById(ctx, nil, in.Id, session.GetDomainId())
	if err != nil {
		return nil, err
	}

	in.NextUploadOn = model.NullTime(calculateNextPeriodFromNow(int32(in.Period)))
	newConfig, err = a.storage.Config().Update(ctx, in, []string{}, session.GetUserId())
	if err != nil {
		return nil, err
	}
	a.UpdateConfigWatchers(ctx, oldConfig, newConfig)
	return newConfig, nil

}

func (a *App) GetSystemObjects(ctx context.Context, in *proto.ReadSystemObjectsRequest) (*proto.SystemObjects, error) {
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
	objects, err := a.storage.Config().GetAvailableSystemObjects(ctx, int(session.GetDomainId()), in.GetIncludeExisting(),
		filters...)
	if err != nil {
		return nil, err
	}
	var r []*proto.Lookup
	for _, v := range objects {
		r = append(r, &proto.Lookup{
			Id:   v.Id.Int32(),
			Name: v.Name.String(),
		})
	}

	return &proto.SystemObjects{
		Items: r,
	}, nil

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
			params.UserId = newConfig.UpdatedBy.Int64()
			params.StorageId = newConfig.Storage.Id.Int()
			params.Period = newConfig.Period
			params.NextUploadOn = newConfig.NextUploadOn.Time()
		} else {
			a.InsertLogUploader(configId, nil, &watcher.UploadWatcherParams{
				StorageId:    newConfig.Storage.Id.Int(),
				Period:       newConfig.Period,
				NextUploadOn: newConfig.NextUploadOn.Time(),
				LastLogId:    newConfig.LastUploadedLog.Int(),
				DomainId:     domainId,
			})
		}
	}
	// else status still disabled
}

func (a *App) InsertConfig(ctx context.Context, in *model.Config) (*model.Config, error) {
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Add)
	if err != nil {
		return nil, err
	}
	// OBAC check
	if !session.CheckObacAccess() {
		return nil, a.MakeScopeError(session.GetMainObjClassName())
	}
	var (
		newModel *model.Config
	)
	if in == nil {
		return nil, errors.New("config is nil")
	}

	in.NextUploadOn = *model.NewNullTime(calculateNextPeriodFromNow(int32(in.Period)))
	newModel, err = a.storage.Config().Insert(ctx, in, session.GetUserId())
	if err != nil {
		return nil, err
	}
	if newModel.Enabled {
		a.InsertLogCleaner(newModel.Id, nil, newModel.DaysToStore)
		a.InsertLogUploader(newModel.Id, nil, &watcher.UploadWatcherParams{
			StorageId:    newModel.Storage.Id.Int(),
			Period:       newModel.Period,
			NextUploadOn: newModel.NextUploadOn.Time(),
			LastLogId:    0,
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

func (a *App) DeleteConfig(ctx context.Context, ids []int32) error {
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Delete)
	if err != nil {
		return err
	}
	if !session.CheckObacAccess() {
		return a.MakeScopeError(session.GetMainObjClassName())
	}
	var rbac *model.RbacOptions
	if session.IsRbacCheckRequired() {
		rbac = &model.RbacOptions{
			Groups: session.GetRoles(),
			Access: session.GetMainAccessMode().Value(),
		}
	}
	_, appErr := a.storage.Config().DeleteMany(ctx, rbac, ids, session.GetDomainId())
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a *App) ConfigCheckAccess(ctx context.Context, id int64, domainId int64, groups []int64, access auth.AccessMode) (bool, error) {
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
