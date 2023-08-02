package app

import (
	"context"
	"fmt"
	"github.com/webitel/wlog"
	"time"
	"webitel_logger/model"
	"webitel_logger/proto"

	errors "github.com/webitel/engine/model"
)

const (
	MEMORY_CACHE_DEFAULT_EXPIRES = 2 * 60
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
	watcherName := FormatKey(DeleteWatcherPrefix, result.DomainId, result.Object.Id.Int())
	if oldModel.Enabled == true && result.Enabled == false {
		a.DeleteWatcherByKey(watcherName)
	} else {
		if a.GetWatcherByKey(watcherName) != nil {
			a.UpdateDeleteWatcherWithNewInterval(result.Id, result.DaysToStore)
		} else {
			a.InsertNewDeleteWatcher(result.Id, result.DaysToStore)
		}
	}

	//}

	res, err := a.convertConfigModelToMessage(result)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *App) PatchUpdateConfig(ctx context.Context, in *proto.PatchUpdateConfigRequest, domainId int, userId int) (*proto.Config, errors.AppError) {
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

	model, err := a.convertPatchUpdateConfigMessageToModel(in, domainId)

	if err != nil {
		return nil, err
	}
	result, err = a.storage.Config().Update(ctx, model, in.GetFields(), userId)
	if err != nil {
		return nil, err
	}
	watcherName := FormatKey(DeleteWatcherPrefix, result.DomainId, result.Object.Id.Int())
	if oldModel.Enabled == true && result.Enabled == false {
		a.DeleteWatcherByKey(watcherName)
	} else {
		if a.GetWatcherByKey(watcherName) != nil {
			a.UpdateDeleteWatcherWithNewInterval(result.Id, result.DaysToStore)
		} else {
			a.InsertNewDeleteWatcher(result.Id, result.DaysToStore)
		}
	}

	//}

	res, err := a.convertConfigModelToMessage(result)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *App) InsertConfig(ctx context.Context, in *proto.InsertConfigRequest, domainId int, userId int) (*proto.Config, errors.AppError) {
	var (
		newModel *model.Config
	)
	if in == nil {
		errors.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	model, err := a.convertInsertConfigMessageToModel(in, domainId)
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

	cacheKey := FormatKey("config.objectId", domainId, objectId)

	value, err := a.cache.Get(ctx, cacheKey)
	if err != nil {
		newModel, appErr := a.storage.Config().GetByObjectId(ctx, domainId, objectId)
		if appErr != nil {
			return nil, appErr
		}
		res, appErr = a.convertConfigModelToMessage(newModel)
		if appErr != nil {
			return nil, appErr
		}
		err := a.cache.Set(ctx, cacheKey, res, MEMORY_CACHE_DEFAULT_EXPIRES)
		if err != nil {
			wlog.Debug(fmt.Sprintf("can't set cache value. error: %s", err.Error()))
		}
	} else {
		res = value.Raw().(*proto.Config)
	}

	return res, nil

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

//func (a *App) ConfigCheckAccess(ctx context.Context, domainId, id int64, groups []int, access auth_manager.PermissionAccess) (bool, errors.AppError) {
//	available, err := a.storage.Config().CheckAccess(ctx, domainId, id, groups, access)
//	if err != nil {
//		return false, err
//	}
//	return available, nil
//
//}
//
//func (a *App) ConfigCheckAccessByObjectId(ctx context.Context, domainId, objectId int64, groups []int, access auth_manager.PermissionAccess) (bool, errors.AppError) {
//	available, err := a.storage.Config().CheckAccess(ctx, domainId, objectId, groups, access)
//	if err != nil {
//		return false, err
//	}
//	return available, nil
//
//}

func (a *App) GetAllConfigs(ctx context.Context, opt *model.SearchOptions, rbac *model.RbacOptions, domainId int) ([]*proto.Config, errors.AppError) {
	var res []*proto.Config
	modelConfigs, err := a.storage.Config().Get(
		ctx,
		opt,
		rbac,
		model.Filter{
			Column:         "object_config.domain_id",
			Value:          domainId,
			ComparisonType: model.Equal,
		},
	)
	if err != nil {
		if IsErrNoRows(err) {
			return res, nil
		} else {
			return nil, err
		}
	}
	for _, v := range *modelConfigs {
		proto, err := a.convertConfigModelToMessage(&v)
		if err != nil {
			return nil, err
		}
		res = append(res, proto)
	}

	return res, nil

}

func (a *App) convertUpdateConfigMessageToModel(in *proto.UpdateConfigRequest, domainId int) (*model.Config, errors.AppError) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      in.GetPeriod(),
		//Storage.Id:  int(in.GetStorageId()),
		DomainId: domainId,
	}
	a.calculateNextPeriod(config)
	config.Storage.Id = model.NewNullInt(int(in.GetStorageId()))
	return config, nil
}

func (a *App) convertPatchUpdateConfigMessageToModel(in *proto.PatchUpdateConfigRequest, domainId int) (*model.Config, errors.AppError) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      in.GetPeriod(),
		//Storage.Id:  int(in.GetStorageId()),
		DomainId: domainId,
	}
	a.calculateNextPeriod(config)
	config.Storage.Id = model.NewNullInt(int(in.GetStorageId()))
	return config, nil
}

func (a *App) convertInsertConfigMessageToModel(in *proto.InsertConfigRequest, domainId int) (*model.Config, errors.AppError) {
	config := &model.Config{

		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      in.GetPeriod(),
		//StorageId:   int(in.GetStorageId()),
		DomainId: domainId,
	}
	a.calculateNextPeriod(config)
	config.Object.Id = model.NewNullInt(int(in.GetObjectId()))
	config.Storage.Id = model.NewNullInt(int(in.GetStorageId()))
	return config, nil
}

func (a *App) convertConfigModelToMessage(in *model.Config) (*proto.Config, errors.AppError) {
	conf := &proto.Config{
		Id:          int32(in.Id),
		Enabled:     in.Enabled,
		DaysToStore: int32(in.DaysToStore),
		Period:      in.Period,
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
	switch in.Period {
	case "everyday":
		in.NextUploadOn = (model.NullTime)(time.Now().AddDate(0, 0, 1))
	case "everyweek":
		in.NextUploadOn = (model.NullTime)(time.Now().AddDate(0, 0, 7))
	case "everytwoweeks":
		in.NextUploadOn = (model.NullTime)(time.Now().AddDate(0, 0, 14))
	default:
		in.Period = "everymonth"
		in.NextUploadOn = (model.NullTime)(time.Now().AddDate(0, 1, 0))
	}
}
