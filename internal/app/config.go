package app

import (
	"context"
	"fmt"
	authmodel "github.com/webitel/logger/internal/auth/model"
	"github.com/webitel/logger/internal/watcher"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"

	proto "github.com/webitel/logger/api/logger"
	"github.com/webitel/logger/internal/model"
)

func (a *App) UpdateConfig(ctx context.Context, in *model.Config) (*model.Config, model.AppError) {
	var (
		newConfig *model.Config
	)
	if in == nil {
		return nil, model.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	session, err := a.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, a.MakeScopeError(session, scope, accessMode)
	}

	// RBAC check
	if session.UseRbacAccess(scope, accessMode) {
		access, err := a.ConfigCheckAccess(ctx, int64(in.Id), session.GetUserId(), session.GetAclRoles(), accessMode)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, a.MakeScopeError(session, scope, accessMode)
		}
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
	session, err := a.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	scope := model.ScopeLog
	accessMode := authmodel.Read
	// OBAC check
	if !session.HasObacAccess(scope, accessMode) {
		return nil, a.MakeScopeError(session, scope, authmodel.Read)
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
	ctx, span := a.tracer.Start(ctx, "worker.UpdateConfigWatchers")
	defer span.End()
	configId := newConfig.Id
	domainId := newConfig.DomainId
	span.SetAttributes(attribute.Int("config.id", configId))
	if !newConfig.Enabled && oldConfig.Enabled { // changed to disabled
		a.DeleteWatchers(configId)
		span.AddEvent("config workers deleted")
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
		span.AddEvent("config workers updated")
	}
	// else status still disabled
}

func (a *App) InsertConfig(ctx context.Context, in *model.Config) (*model.Config, model.AppError) {
	var (
		newModel *model.Config
	)
	if in == nil {
		return nil, model.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	session, err := a.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Add
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, a.MakeScopeError(session, scope, accessMode)
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

func (a *App) GetConfigByObjectId(ctx context.Context, objectId int, domainId int) (*model.Config, model.AppError) {
	res, appErr := a.storage.Config().GetByObjectId(ctx, domainId, objectId)
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (a *App) CheckConfigStatus(ctx context.Context, objectName string, domainId int) (bool, model.AppError) {
	searchResult, appErr := a.storage.Config().Select(ctx, nil, nil, &model.FilterNode{
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
	if appErr != nil {
		if IsErrNoRows(appErr) {
			return false, nil
		}

		return false, appErr
	}
	enabled := searchResult[0].Enabled

	return enabled, nil

}

func (a *App) GetConfigById(ctx context.Context, rbac *model.RbacOptions, id int) (*model.Config, model.AppError) {
	session, err := a.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	scope := model.ScopeLog
	accessMode := authmodel.Read
	// OBAC check
	if !session.HasObacAccess(scope, accessMode) {
		return nil, a.MakeScopeError(session, scope, accessMode)
	}
	// RBAC check
	if session.UseRbacAccess(scope, accessMode) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: accessMode.Value(),
		}
	}
	res, appErr := a.storage.Config().GetById(ctx, rbac, id, session.GetDomainId())
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (a *App) DeleteConfig(ctx context.Context, ids []int32) model.AppError {
	session, err := a.AuthorizeFromContext(ctx)
	if err != nil {
		return err
	}
	accessMode := authmodel.Edit
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return a.MakeScopeError(session, scope, accessMode)
	}
	var rbac *model.RbacOptions
	if session.UseRbacAccess(scope, accessMode) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: accessMode.Value(),
		}
	}
	appErr := a.storage.Config().DeleteMany(ctx, rbac, ids, session.GetDomainId())
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a *App) ConfigCheckAccess(ctx context.Context, id int64, domainId int64, groups []int64, access authmodel.AccessMode) (bool, model.AppError) {
	return a.storage.Config().CheckAccess(ctx, domainId, id, groups, access.Value())

}

func (a *App) SearchConfig(ctx context.Context, rbac *model.RbacOptions, searchOpt *model.SearchOptions) ([]*model.Config, model.AppError) {
	session, err := a.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	scope := model.ScopeLog
	accessMode := authmodel.Read
	// OBAC check
	if !session.HasObacAccess(scope, accessMode) {
		return nil, a.MakeScopeError(session, scope, accessMode)
	}
	// RBAC check
	if session.UseRbacAccess(scope, accessMode) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: accessMode.Value(),
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
