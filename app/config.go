package app

import (
	"context"
	"strings"
	"time"

	"github.com/webitel/engine/auth_manager"
	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"
	proto "github.com/webitel/protos/logger"
)

func (a *App) UpdateConfig(ctx context.Context, in *proto.UpdateConfigRequest, domainId int, userId int) (*proto.Config, errors.AppError) {
	var (
		result *model.Config
	)
	if in == nil {
		return nil, errors.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	oldModel, err := a.storage.Config().GetById(ctx, nil, int(in.GetConfigId()))
	if err != nil {
		return nil, err
	}

	model, err := a.convertUpdateConfigMessageToModel(in, domainId)

	if err != nil {
		return nil, err
	}
	result, err = a.storage.Config().Update(ctx, model, []string{}, userId)
	if err != nil {
		return nil, err
	}

	a.UpdateDeleteWatcherForConfig(oldModel.Enabled == true && result.Enabled == false, result.DaysToStore, result.Id, domainId, result.Object.Id.Int())

	res, err := a.convertConfigModelToMessage(result)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *App) PatchUpdateConfig(ctx context.Context, in *proto.PatchConfigRequest, domainId int, userId int) (*proto.Config, errors.AppError) {
	var (
		result *model.Config
	)
	if in == nil {
		return nil, errors.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	oldModel, err := a.storage.Config().GetById(ctx, nil, int(in.GetConfigId()))
	if err != nil {
		return nil, err
	}

	model, err := a.convertPatchConfigMessageToModel(in, domainId)

	if err != nil {
		return nil, err
	}
	result, err = a.storage.Config().Update(ctx, model, in.GetFields(), userId)
	if err != nil {
		return nil, err
	}

	a.UpdateDeleteWatcherForConfig(oldModel.Enabled == true && result.Enabled == false, result.DaysToStore, result.Id, domainId, result.Object.Id.Int())

	res, err := a.convertConfigModelToMessage(result)
	if err != nil {
		return nil, err
	}
	return res, nil

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

func (a *App) UpdateDeleteWatcherForConfig(statusChangedToDisabled bool, daysToStore, configId, domainId, objectId int) {
	watcherName := FormatKey(DeleteWatcherPrefix, domainId, objectId)
	if statusChangedToDisabled {
		a.DeleteWatcherByKey(watcherName)
	} else {
		if a.GetWatcherByKey(watcherName) != nil {
			a.UpdateDeleteWatcherWithNewInterval(configId, daysToStore)
		} else {
			a.InsertNewDeleteWatcher(configId, daysToStore)
		}
	}
}

func (a *App) InsertConfig(ctx context.Context, in *proto.CreateConfigRequest, domainId int, userId int) (*proto.Config, errors.AppError) {
	var (
		newModel *model.Config
	)
	if in == nil {
		errors.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	model, err := a.convertCreateConfigMessageToModel(in, domainId)
	if err != nil {
		return nil, err
	}
	newModel, err = a.storage.Config().Insert(ctx, model, userId)
	if err != nil {
		return nil, err
	}
	if newModel.Enabled {
		a.InsertNewDeleteWatcher(newModel.Id, newModel.DaysToStore)
	}

	res, err := a.convertConfigModelToMessage(newModel)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *App) GetConfigByObjectId(ctx context.Context /*opt *model.SearchOptions,*/, domainId int, objectId int) (*proto.Config, errors.AppError) {

	var res *proto.Config

	//cacheKey := FormatKey("config.objectId", domainId, objectId)

	//value, err := a.cache.Get(ctx, cacheKey)
	//if err != nil {
	newModel, appErr := a.storage.Config().GetByObjectId(ctx, domainId, objectId)
	if appErr != nil {
		return nil, appErr
	}
	res, appErr = a.convertConfigModelToMessage(newModel)
	if appErr != nil {
		return nil, appErr
	}
	//err := a.cache.Set(ctx, cacheKey, res, MEMORY_CACHE_DEFAULT_EXPIRES)
	//if err != nil {
	//	wlog.Debug(fmt.Sprintf("can't set cache value. error: %s", err.Error()))
	//}
	//} else {
	//	res = value.Raw().(*proto.Config)
	//}

	return res, nil

}

func (a *App) CheckConfigStatus(ctx context.Context, in *proto.CheckConfigStatusRequest) (*proto.ConfigStatus, errors.AppError) {

	var (
		response proto.ConfigStatus
	)

	//cacheKey := FormatKey("config.objectId", in.GetDomainId(), in.GetObjectName())
	objectName := strings.ToLower(in.GetObjectName())
	domainId := in.GetDomainId()

	//value, err := a.cache.Get(ctx, cacheKey)
	//if err != nil {
	searchResult, appErr := a.storage.Config().Get(ctx, nil, nil, model.FilterBunch{
		Bunch: []*model.Filter{
			{
				Column:         "wbt_class.name",
				Value:          objectName,
				ComparisonType: model.ILike,
			},
			{
				Column:         "object_config.domain_id",
				Value:          domainId,
				ComparisonType: model.Equal,
			}},
		ConnectionType: model.AND,
	})
	if appErr != nil {
		if IsErrNoRows(appErr) {
			response.IsEnabled = false
			return &response, nil
		}

		return nil, appErr
	}
	enabled := searchResult[0].Enabled
	response.IsEnabled = enabled

	return &response, nil

}

func (a *App) GetConfigById(ctx context.Context, rbac *model.RbacOptions, id int) (*proto.Config, errors.AppError) {
	var res *proto.Config
	newModel, appErr := a.storage.Config().GetById(ctx, rbac, id)
	if appErr != nil {
		return nil, appErr
	}
	res, appErr = a.convertConfigModelToMessage(newModel)
	if appErr != nil {
		return nil, appErr
	}
	return res, nil
}

func (a *App) DeleteConfig(ctx context.Context, id int32) errors.AppError {
	appErr := a.storage.Config().Delete(ctx, id)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a *App) DeleteConfigs(ctx context.Context, rbac *model.RbacOptions, ids []int32) errors.AppError {
	appErr := a.storage.Config().DeleteMany(ctx, rbac, ids)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (a *App) ConfigCheckAccess(ctx context.Context, domainId, id int64, groups []int, access auth_manager.PermissionAccess) (bool, errors.AppError) {
	available, err := a.storage.Config().CheckAccess(ctx, domainId, id, groups, access.Value())
	if err != nil {
		if IsErrNoRows(err) {
			return false, nil
		}
		return false, err
	}
	return available, nil

}

func (a *App) GetAllConfigs(ctx context.Context, rbac *model.RbacOptions, domainId int, in *proto.SearchConfigRequest) (*proto.Configs, errors.AppError) {
	var (
		rows      []*proto.Config
		res       proto.Configs
		searchOpt *model.SearchOptions
	)

	searchOpt = ExtractSearchOptions(in)
	modelConfigs, err := a.storage.Config().Get(
		ctx,
		searchOpt,
		rbac,
		model.Filter{
			Column:         "object_config.domain_id",
			Value:          domainId,
			ComparisonType: model.Equal,
		},
	)
	res.Page = int32(searchOpt.Page)
	if err != nil {
		if IsErrNoRows(err) {
			return &res, nil
		} else {
			return nil, err
		}
	}
	for _, v := range modelConfigs {
		proto, err := a.convertConfigModelToMessage(v)
		if err != nil {
			return nil, err
		}
		rows = append(rows, proto)
	}
	if len(rows)-1 == searchOpt.Size {
		res.Next = true
		res.Items = rows[0 : len(rows)-1]
	} else {
		res.Items = rows
	}

	return &res, nil

}

func (a *App) convertUpdateConfigMessageToModel(in *proto.UpdateConfigRequest, domainId int) (*model.Config, errors.AppError) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		//Storage.Id:  int(in.GetStorageId()),
		DomainId:    domainId,
		Description: *model.NewNullString(in.GetDescription()),
	}
	a.calculateNextPeriod(config)

	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, errors.NewInternalError("app.config.convert_update_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func (a *App) convertPatchConfigMessageToModel(in *proto.PatchConfigRequest, domainId int) (*model.Config, errors.AppError) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		//Storage.Id:  int(in.GetStorageId()),
		DomainId:    domainId,
		Description: *model.NewNullString(in.GetDescription()),
	}
	a.calculateNextPeriod(config)

	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, errors.NewInternalError("app.config.convert_patch_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func (a *App) convertCreateConfigMessageToModel(in *proto.CreateConfigRequest, domainId int) (*model.Config, errors.AppError) {

	if in.GetDaysToStore() <= 0 {
		return nil, errors.NewBadRequestError("app.config.convert_create_config_message.bad_args", "days to store should be greater than one")
	}
	config := &model.Config{

		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		//StorageId:   int(in.GetStorageId()),
		DomainId:    domainId,
		Description: *model.NewNullString(in.GetDescription()),
	}
	a.calculateNextPeriod(config)
	objectId, err := model.NewNullInt(in.GetObject().GetId())
	if err != nil {
		return nil, errors.NewInternalError("app.config.convert_create_config_message.convert_object_id.fail", err.Error())
	}
	config.Object.Id = objectId

	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, errors.NewInternalError("app.config.convert_create_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func (a *App) convertConfigModelToMessage(in *model.Config) (*proto.Config, errors.AppError) {

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

func (a *App) calculateNextPeriod(in *model.Config) {
	in.NextUploadOn = *model.NewNullTime(time.Now().Add(time.Hour * 24 * time.Duration(in.Period)))
}
