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
		newModel    *model.Config
		noRowsError bool
	)
	if in == nil {
		errors.NewInternalError("app.app.update_config.check_arguments.fail", "config proto is nil")
	}
	_, err := a.storage.Config().GetByObjectId(ctx, int(in.GetObjectId()), int(in.GetDomainId()))
	if err != nil {
		if noRowsError = IsErrNoRows(err); !noRowsError {
			return nil, err
		}
	}

	model, err := a.convertConfigMessageToModel(in)
	if err != nil {
		return nil, err
	}
	a.CalculateNextPeriod(model)
	if noRowsError {
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

	res, err := a.convertConfigModelToMessage(newModel)
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
	res, err := a.convertConfigModelToMessage(newModel)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (a *App) GetAllConfigs(ctx context.Context, domainId int) ([]*proto.Config, errors.AppError) {
	var res []*proto.Config
	modelConfigs, err := a.storage.Config().GetAll(ctx, domainId)
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

func (a *App) convertConfigMessageToModel(in *proto.Config) (*model.Config, errors.AppError) {
	return &model.Config{
		Id:          int(in.GetId()),
		ObjectId:    int(in.GetObjectId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      in.GetPeriod(),
		StorageId:   int(in.GetStorageId()),
		DomainId:    int(in.GetDomainId()),
	}, nil
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

func (a *App) CalculateNextPeriod(in *model.Config) {
	switch in.Period {
	case "everyday":
		in.NextUploadOn = time.Now().AddDate(0, 0, 1).Unix()
	case "everyweek":
		in.NextUploadOn = time.Now().AddDate(0, 0, 7).Unix()
	case "everytwoweeks":
		in.NextUploadOn = time.Now().AddDate(0, 0, 14).Unix()
	default:
		in.Period = "everymonth"
		in.NextUploadOn = time.Now().AddDate(0, 1, 0).Unix()
	}
}
