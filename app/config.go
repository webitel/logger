package app

import (
	"context"
	"fmt"
	authmodel "github.com/webitel/logger/auth/model"
	"go.opentelemetry.io/otel/attribute"
	"time"

	"github.com/webitel/logger/watcher"

	"log/slog"

	proto "buf.build/gen/go/webitel/logger/protocolbuffers/go"
	"github.com/webitel/logger/model"
)

func (a *App) UpdateConfig(ctx context.Context, in *model.Config, userId int64, domainId int64) (*model.Config, model.AppError) {
	var (
		newConfig *model.Config
	)
	if in == nil {
		return nil, model.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	oldConfig, err := a.storage.Config().GetById(ctx, nil, in.Id, domainId)
	if err != nil {
		return nil, err
	}
	newConfig, err = a.storage.Config().Update(ctx, in, []string{}, userId)
	if err != nil {
		return nil, err
	}
	a.UpdateConfigWatchers(ctx, oldConfig, newConfig)
	return newConfig, nil

}

func (a *App) PatchUpdateConfig(ctx context.Context, in *model.Config, fields []string, userId int64, domainId int64) (*model.Config, model.AppError) {
	var (
		newConfig *model.Config
	)
	if in == nil {
		return nil, model.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	oldConfig, err := a.storage.Config().GetById(ctx, nil, in.Id, domainId)
	if err != nil {
		return nil, err
	}
	newConfig, err = a.storage.Config().Update(ctx, in, fields, userId)
	if err != nil {
		return nil, err
	}
	a.UpdateConfigWatchers(ctx, oldConfig, newConfig)

	return newConfig, nil

}

func (a *App) GetSystemObjects(ctx context.Context, in *proto.ReadSystemObjectsRequest, domainId int) (*proto.SystemObjects, error) {
	var filters []string
	for _, name := range proto.AvailableSystemObjects_name {
		filters = append(filters, name)
	}
	objects, err := a.storage.Config().GetAvailableSystemObjects(ctx, domainId, in.GetIncludeExisting(),
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

func (a *App) InsertConfig(ctx context.Context, in *model.Config, userId int64) (*model.Config, model.AppError) {
	var (
		newModel *model.Config
	)
	if in == nil {
		return nil, model.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}

	newModel, err := a.storage.Config().Insert(ctx, in, userId)
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

func (a *App) GetConfigByObjectId(ctx context.Context /*opt *model.SearchOptions,*/, domainId int, objectId int) (*model.Config, model.AppError) {
	res, appErr := a.storage.Config().GetByObjectId(ctx, domainId, objectId)
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (a *App) CheckConfigStatus(ctx context.Context, objectName string, domainId int64) (bool, model.AppError) {
	searchResult, appErr := a.storage.Config().Get(ctx, nil, nil, &model.FilterNode{
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

func (a *App) GetConfigById(ctx context.Context, rbac *model.RbacOptions, id int, domainId int64) (*model.Config, model.AppError) {
	res, appErr := a.storage.Config().GetById(ctx, rbac, id, domainId)
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (a *App) DeleteConfig(ctx context.Context, id int32, domainId int64) model.AppError {
	appErr := a.storage.Config().Delete(ctx, id, domainId)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a *App) DeleteConfigs(ctx context.Context, rbac *model.RbacOptions, ids []int32, domainId int64) model.AppError {
	appErr := a.storage.Config().DeleteMany(ctx, rbac, ids, domainId)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a *App) ConfigCheckAccess(ctx context.Context, domainId, id int64, groups []int, access authmodel.AccessMode) (bool, model.AppError) {
	available, err := a.storage.Config().CheckAccess(ctx, domainId, id, groups, access.Value())
	if err != nil {
		if IsErrNoRows(err) {
			return false, nil
		}
		return false, err
	}
	return available, nil

}

func (a *App) GetAllConfigs(ctx context.Context, rbac *model.RbacOptions, searchOpt *model.SearchOptions, domainId int64) ([]*model.Config, model.AppError) {
	modelConfigs, err := a.storage.Config().Get(
		ctx,
		searchOpt,
		rbac,
		&model.Filter{
			Column:         "object_config.domain_id",
			Value:          domainId,
			ComparisonType: model.Equal,
		},
	)
	if err != nil {
		return nil, err
	}

	return modelConfigs, nil

}

func ConvertUpdateConfigMessageToModel(in *proto.UpdateConfigRequest, domainId int64) (*model.Config, model.AppError) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		//Storage.Id:  int(in.GetStorageId()),
		DomainId:     domainId,
		Description:  *model.NewNullString(in.GetDescription()),
		NextUploadOn: *model.NewNullTime(calculateNextPeriodFromNow(in.Period)),
	}

	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, model.NewInternalError("app.config.convert_update_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func ConvertPatchConfigMessageToModel(in *proto.PatchConfigRequest, domainId int64) (*model.Config, model.AppError) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		//Storage.Id:  int(in.GetStorageId()),
		DomainId:     domainId,
		Description:  *model.NewNullString(in.GetDescription()),
		NextUploadOn: *model.NewNullTime(calculateNextPeriodFromNow(in.Period)),
	}
	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, model.NewInternalError("app.config.convert_patch_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func ConvertCreateConfigMessageToModel(in *proto.CreateConfigRequest, domainId int64) (*model.Config, model.AppError) {

	if in.GetDaysToStore() <= 0 {
		return nil, model.NewBadRequestError("app.config.convert_create_config_message.bad_args", "days to store should be greater than one")
	}
	config := &model.Config{

		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		//StorageId:   int(in.GetStorageId()),
		DomainId:     domainId,
		Description:  *model.NewNullString(in.GetDescription()),
		NextUploadOn: *model.NewNullTime(calculateNextPeriodFromNow(in.Period)),
	}
	objectId, err := model.NewNullInt(in.GetObject().GetId())
	if err != nil {
		return nil, model.NewInternalError("app.config.convert_create_config_message.convert_object_id.fail", err.Error())
	}
	config.Object.Id = objectId

	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, model.NewInternalError("app.config.convert_create_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func ConvertConfigModelToMessage(in *model.Config) (*proto.Config, model.AppError) {

	conf := &proto.Config{
		Id:          int32(in.Id),
		Enabled:     in.Enabled,
		DaysToStore: int32(in.DaysToStore),
		Period:      int32(in.Period),
		Description: in.Description.String(),
		//DomainId:    int32(in.DomainId),
	}
	if !in.Object.IsZero() {
		conf.Object = &proto.Lookup{
			Id:   int32(in.Object.Id.Int()),
			Name: in.Object.Name.String(),
		}
	}
	if !in.Storage.IsZero() {
		conf.Storage = &proto.Lookup{
			Id:   int32(in.Storage.Id.Int()),
			Name: in.Storage.Name.String(),
		}
	}
	return conf, nil
}

func calculateNextPeriodFromDate(period int, from time.Time) *model.NullTime {
	return model.NewNullTime(from.Add(time.Hour * 24 * time.Duration(period)))

}

func calculateNextPeriodFromNow(period int32) time.Time {
	now := time.Now().Add(time.Hour * 24 * time.Duration(period))
	return now
}
