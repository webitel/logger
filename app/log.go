package app

import (
	"context"
	"webitel_logger/model"
	"webitel_logger/proto"

	errors "github.com/webitel/engine/model"
)

func (a *App) GetLogsByObjectId(ctx context.Context, domainId, objectId int) ([]*proto.Log, errors.AppError) {
	var result []*proto.Log
	rows, appErr := a.storage.Log().GetByObjectId(ctx, objectId, domainId)
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

func (a *App) GetLogsByUserId(ctx context.Context, userId int) ([]*proto.Log, errors.AppError) {
	var result []*proto.Log
	rows, appErr := a.storage.Log().GetByUserId(ctx, userId)
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

func convertLogModelToMessage(m *model.Log) (*proto.Log, errors.AppError) {
	return &proto.Log{
		Id:       int32(m.Id),
		Action:   m.Action,
		Date:     m.Date,
		UserId:   int32(m.UserId),
		UserIp:   m.UserIp,
		ObjectId: int32(m.ObjectId),
		NewState: m.NewState,
		DomainId: int32(m.DomainId),
	}, nil
}
