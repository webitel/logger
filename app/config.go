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

func (a *App) UpdateConfig(ctx context.Context, in *proto.UpdateConfigRequest) (*proto.Config, errors.AppError) {
	var (
		newModel    *model.Config
		noRowsError bool
	)
	if in == nil {
		errors.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	oldModel, err := a.storage.Config().GetByObjectId(ctx, int(in.GetDomainId()), int(in.GetObjectId()))
	if err != nil {
		if noRowsError = IsErrNoRows(err); !noRowsError {
			return nil, err
		}
	}

	model, err := a.convertConfigMessageToModel(in)

	if err != nil {
		return nil, err
	}
	if noRowsError {
		newModel, err = a.storage.Config().Insert(ctx, model)
		if err != nil {
			return nil, err
		}
		if newModel.Enabled {
			a.InsertNewDeleteWatcher(newModel.DomainId, newModel.ObjectId, newModel.DaysToStore)
		}
	} else {
		newModel, err = a.storage.Config().Update(ctx, model)
		if err != nil {
			return nil, err
		}
		watcherName := FormatConfigKey(DeleteWatcherPrefix, newModel.DomainId, newModel.ObjectId)
		if oldModel.Enabled == true && newModel.Enabled == false {
			a.DeleteWatcherByKey(watcherName)
		} else {
			if a.GetWatcherByKey(watcherName) != nil {
				a.UpdateDeleteWatcherWithNewInterval(newModel.DomainId, newModel.ObjectId, newModel.DaysToStore)
			} else {
				a.InsertNewDeleteWatcher(newModel.DomainId, newModel.ObjectId, newModel.DaysToStore)
			}
		}

	}

	res, err := a.convertConfigModelToMessage(newModel)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *App) GetConfigByObjectId(ctx context.Context /*opt *model.SearchOptions,*/, domainId int, objectId int) (*proto.Config, errors.AppError) {

	var res *proto.Config

	cacheKey := FormatConfigKey("config", domainId, objectId)

	value, err := a.cache.Get(ctx, cacheKey)
	if err != nil {
		newModel, appErr := a.storage.Config().GetByObjectId(ctx /*opt,*/, domainId, objectId)
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

func (a *App) GetAllConfigs(ctx context.Context, opt *model.SearchOptions, domainId int) ([]*proto.Config, errors.AppError) {
	var res []*proto.Config
	modelConfigs, err := a.storage.Config().GetAll(ctx, opt, domainId)
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

func (a *App) convertConfigMessageToModel(in *proto.UpdateConfigRequest) (*model.Config, errors.AppError) {
	config := &model.Config{
		ObjectId:    int(in.GetObjectId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      in.GetPeriod(),
		StorageId:   int(in.GetStorageId()),
		DomainId:    int(in.GetDomainId()),
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
		in.NextUploadOn = time.Now().AddDate(0, 0, 1)
	case "everyweek":
		in.NextUploadOn = time.Now().AddDate(0, 0, 7)
	case "everytwoweeks":
		in.NextUploadOn = time.Now().AddDate(0, 0, 14)
	default:
		in.Period = "everymonth"
		in.NextUploadOn = time.Now().AddDate(0, 1, 0)
	}
}
