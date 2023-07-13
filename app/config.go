package app

import (
	"context"
	"time"
	"webitel_logger/model"
	"webitel_logger/proto"

	errors "github.com/webitel/engine/model"
)

func (a *App) UpdateConfig(ctx context.Context, in *proto.Config) (*proto.Config, errors.AppError) {
	var (
		newModel *model.Config
	)
	if in == nil {
		errors.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	config, err := a.storage.Config().GetByObjectId(ctx, int(in.GetObjectId()), int(in.GetDomainId()))
	if err != nil {
		return nil, err
	}

	model, err := a.convertConfigToModel(in)
	if err != nil {
		return nil, err
	}
	model.NextUploadOn = a.CalculateNextPeriod(model.Period)
	if config == nil {
		newModel, err = a.storage.Config().Insert(ctx, model)
		if err != nil {
			return nil, err
		}
	} else {
		newModel, err = a.storage.Config().Update(ctx, model)
		if err != nil {
			return nil, err
		}
	}

	res, err := a.convertModelToConfig(newModel)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *App) GetConfigByObjectId(ctx context.Context, objectId int, domainId int) (*proto.Config, errors.AppError) {
	newModel, err := a.storage.Config().GetByObjectId(ctx, objectId, domainId)
	if err != nil {
		return nil, err
	}
	res, err := a.convertModelToConfig(newModel)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *App) convertConfigToModel(in *proto.Config) (*model.Config, errors.AppError) {
	return &model.Config{
		ObjectId:    int(in.GetObjectId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      in.GetPeriod(),
		StorageId:   int(in.GetStorageId()),
		DomainId:    int(in.GetDomainId()),
	}, nil
}

func (a *App) convertModelToConfig(in *model.Config) (*proto.Config, errors.AppError) {
	return &proto.Config{
		ObjectId:    int32(in.ObjectId),
		Enabled:     in.Enabled,
		DaysToStore: int32(in.DaysToStore),
		Period:      in.Period,
		StorageId:   int32(in.StorageId),
		DomainId:    int32(in.DomainId),
	}, nil
}

func (a *App) CalculateNextPeriod(in string) time.Time {
	var res time.Time
	switch in {
	case "everyday":
		res = time.Now().AddDate(0, 0, 1)
	case "everyweek":
		res = time.Now().AddDate(0, 0, 7)
	case "everytwoweeks":
		res = time.Now().AddDate(0, 0, 14)
	default:
		res = time.Now().AddDate(0, 1, 0)
	}
	return res
}
