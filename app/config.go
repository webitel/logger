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
		newModel *model.Config
	)
	if in == nil {
		errors.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	oldModel, err := a.storage.Config().GetByObjectId(ctx, domainId, int(in.GetObjectId()))
	if err != nil {
		return nil, err
	}

	model, err := a.convertUpdateConfigMessageToModel(in, domainId)

	if err != nil {
		return nil, err
	}
	newModel, err = a.storage.Config().Update(ctx, model, userId)
	if err != nil {
		return nil, err
	}
	watcherName := FormatKey(DeleteWatcherPrefix, newModel.DomainId, newModel.ObjectId)
	if oldModel.Enabled == true && newModel.Enabled == false {
		a.DeleteWatcherByKey(watcherName)
	} else {
		if a.GetWatcherByKey(watcherName) != nil {
			a.UpdateDeleteWatcherWithNewInterval(newModel.Id, newModel.DaysToStore)
		} else {
			a.InsertNewDeleteWatcher(newModel.Id, newModel.DaysToStore)
		}
	}

	//}

	res, err := a.convertConfigModelToMessage(newModel)
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
	modelConfigs, err := a.storage.Config().GetAll(ctx, opt, rbac, domainId)
	if err != nil {
		return nil, err
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
		ObjectId:    int(in.GetObjectId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      in.GetPeriod(),
		StorageId:   int(in.GetStorageId()),
		DomainId:    domainId,
	}
	a.calculateNextPeriod(config)
	return config, nil
}

func (a *App) convertInsertConfigMessageToModel(in *proto.InsertConfigRequest, domainId int) (*model.Config, errors.AppError) {
	config := &model.Config{
		ObjectId:    int(in.GetObjectId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      in.GetPeriod(),
		StorageId:   int(in.GetStorageId()),
		DomainId:    domainId,
	}
	a.calculateNextPeriod(config)
	return config, nil
}

func (a *App) convertConfigModelToMessage(in *model.Config) (*proto.Config, errors.AppError) {
	return &proto.Config{
		Id:          int32(in.Id),
		ObjectId:    int32(in.ObjectId),
		Enabled:     in.Enabled,
		DaysToStore: int32(in.DaysToStore),
		Period:      in.Period,
		StorageId:   int32(in.StorageId),
		DomainId:    int32(in.DomainId),
	}, nil
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
