package app

import (
	"context"
	"time"
	"webitel_logger/model"
	"webitel_logger/proto"

	errors "github.com/webitel/engine/model"
)

func (a *App) GetLogsByObjectId(ctx context.Context, opt *model.SearchOptions, domainId, objectId int) ([]*proto.Log, errors.AppError) {
	var result []*proto.Log
	rows, appErr := a.storage.Log().GetByObjectId(ctx, opt, domainId, objectId)
	if appErr != nil {
		return nil, appErr
	}
	for _, v := range *rows {
		protoLog, err := convertLogModelToMessage(&v)
		if err != nil {
			return nil, err
		}
		result = append(result, protoLog)
	}
	return result, nil
}

func (a *App) GetLogsByUserId(ctx context.Context, opt *model.SearchOptions, userId int) ([]*proto.Log, errors.AppError) {
	var result []*proto.Log
	rows, appErr := a.storage.Log().GetByUserId(ctx, opt, userId)
	if appErr != nil {
		return nil, appErr
	}
	for _, v := range *rows {
		protoLog, err := convertLogModelToMessage(&v)
		if err != nil {
			return nil, err
		}
		result = append(result, protoLog)
	}
	return result, nil
}
func (a *App) GetLogsByConfigId(ctx context.Context, opt *model.SearchOptions, configId int) ([]*proto.Log, errors.AppError) {
	var result []*proto.Log
	rows, appErr := a.storage.Log().GetByConfigId(ctx, opt, configId)
	if appErr != nil {
		return nil, appErr
	}
	for _, v := range *rows {
		protoLog, err := convertLogModelToMessage(&v)
		if err != nil {
			return nil, err
		}
		result = append(result, protoLog)
	}
	return result, nil
}

func (a *App) InsertLogByRabbitMessage(ctx context.Context, rabbitMessage *model.RabbitMessage) errors.AppError {
	model, err := convertRabbitMessageToModel(rabbitMessage)
	if err != nil {
		return err
	}
	newModel, err := a.storage.Config().GetByObjectId(ctx /*opt,*/, rabbitMessage.DomainId, rabbitMessage.ObjectId)
	if err != nil {
		return err
	}
	model.ConfigId = newModel.Id
	_, err = a.storage.Log().Insert(ctx, model)
	if err != nil {
		return err
	}

	return nil

}

func convertLogModelToMessage(m *model.Log) (*proto.Log, errors.AppError) {
	return &proto.Log{
		Id:     int32(m.Id),
		Action: m.Action,
		Date:   m.Date.String(),
		User: &proto.Lookup{
			Id:   int32(m.User.Id),
			Name: m.User.Name,
		},
		UserIp:   m.UserIp,
		NewState: m.NewState,
		ConfigId: int32(m.ConfigId),
	}, nil
}
func convertRabbitMessageToModel(m *model.RabbitMessage) (*model.Log, errors.AppError) {
	log := &model.Log{
		Action:   m.Action,
		Date:     time.Unix(m.Date, 0),
		UserIp:   m.UserIp,
		NewState: string(m.NewState),
		RecordId: m.RecordId,
	}
	log.User.Id = m.UserId
	return log, nil
}
