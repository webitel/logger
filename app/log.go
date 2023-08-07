package app

import (
	"context"
	"logger/model"
	"logger/proto"
	"time"

	errors "github.com/webitel/engine/model"
)

//func (a *App) GetLogsByObjectId(ctx context.Context, opt *model.SearchOptions, domainId, objectId int) ([]*proto.Log, errors.AppError) {
//	var result []*proto.Log
//	rows, appErr := a.storage.Log().GetByObjectId(ctx, opt, domainId, objectId)
//	if appErr != nil {
//		if IsErrNoRows(appErr) {
//			return result, nil
//		} else {
//			return nil, appErr
//		}
//	}
//	for _, v := range *rows {
//		protoLog, err := convertLogModelToMessage(&v)
//		if err != nil {
//			return nil, err
//		}
//		result = append(result, protoLog)
//	}
//	return result, nil
//}

func (a *App) GetLogsByUserId(ctx context.Context, opt *model.SearchOptions, userId int) ([]*proto.Log, errors.AppError) {
	var result []*proto.Log
	rows, appErr := a.storage.Log().Get(
		ctx,
		opt,
		model.Filter{
			Column:         "log.user_id",
			Value:          userId,
			ComparisonType: model.Equal,
		})
	if appErr != nil {
		if IsErrNoRows(appErr) {
			return result, nil
		} else {
			return nil, appErr
		}
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
	rows, appErr := a.storage.Log().Get(
		ctx,
		opt,
		model.Filter{
			Column:         "log.config_id",
			Value:          configId,
			ComparisonType: 0,
		},
	)
	if appErr != nil {
		if IsErrNoRows(appErr) {
			return result, nil
		} else {
			return nil, appErr
		}
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

func (a *App) InsertLogByRabbitMessage(ctx context.Context, rabbitMessage *model.RabbitMessage, domainId, objectId int) errors.AppError {

	config, err := a.storage.Config().GetByObjectId(ctx, domainId, objectId)
	if err != nil {
		return err
	}
	model, err := convertRabbitMessageToModel(rabbitMessage, config.Id)
	if err != nil {
		return err
	}
	_, err = a.storage.Log().Insert(ctx, model)
	if err != nil {
		return err
	}

	return nil

}

func (a *App) InsertRabbitLogs(ctx context.Context, rabbitMessages []*model.RabbitMessage, domainId, objectId int) errors.AppError {
	config, err := a.storage.Config().GetByObjectId(ctx, int(domainId), int(objectId))
	if err != nil {
		return err
	}
	logs, err := bulkConvertRabbitMessageToModel(rabbitMessages, config.Id)
	if err != nil {
		return err
	}
	err = a.storage.Log().InsertMany(ctx, *logs)
	if err != nil {
		return err
	}
	return nil

}

func convertLogModelToMessage(m *model.Log) (*proto.Log, errors.AppError) {
	log := &proto.Log{
		Id:     int32(m.Id),
		Action: m.Action,
		Date:   m.Date.String(),

		UserIp:   m.UserIp,
		NewState: m.NewState,
		ConfigId: int32(m.ConfigId),
	}
	if !m.User.IsZero() {
		log.User = &proto.Lookup{
			Id:   int32(m.User.Id.Int()),
			Name: m.User.Name.String(),
		}
	}
	return log, nil
}
func convertRabbitMessageToModel(m *model.RabbitMessage, configId int) (*model.Log, errors.AppError) {
	log := &model.Log{
		Action:   m.Action,
		Date:     (model.NullTime)(time.Unix(m.Date, 0)),
		UserIp:   m.UserIp,
		NewState: string(m.NewState),
		RecordId: m.RecordId,
		ConfigId: configId,
	}
	err := log.User.Id.Scan(m.UserId)
	if err != nil {
		return nil, errors.NewBadRequestError("app.log.convert_rabbit_message_to_model.scan.error", err.Error())
	}

	return log, nil
}

func bulkConvertRabbitMessageToModel(m []*model.RabbitMessage, configId int) (*[]*model.Log, errors.AppError) {
	var logs []*model.Log
	for _, v := range m {
		log, err := convertRabbitMessageToModel(v, configId)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return &logs, nil
}
