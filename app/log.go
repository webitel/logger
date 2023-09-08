package app

import (
	"context"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/proto"
	"strings"
	"time"

	errors "github.com/webitel/engine/model"
)

// region PERFORM ACTIONS
func (a *App) SearchLogsByUserId(ctx context.Context, in *proto.SearchLogByUserIdRequest) (*proto.Logs, errors.AppError) {
	var (
		res                proto.Logs
		rows               []*proto.Log
		notStandartFilters []model.Filter
	)

	searchOpt := ExtractSearchOptions(in)
	// region APPLYING FILTERS
	if x := in.GetObject(); x != nil {
		if v := x.GetId(); v != 0 {
			notStandartFilters = append(notStandartFilters, model.Filter{
				Column:         "wbt_class.id",
				Value:          v,
				ComparisonType: model.Equal,
			})
		} else if v := x.GetName(); v != "" {
			notStandartFilters = append(notStandartFilters, model.Filter{
				Column:         "wbt_class.name",
				Value:          strings.Replace(v, "*", "%", -1),
				ComparisonType: model.ILike,
			})
		}
	}
	// REQUIRED !
	notStandartFilters = append(notStandartFilters, model.Filter{
		Column:         "user_id",
		Value:          in.GetUserId(),
		ComparisonType: model.Equal,
	})

	// endregion
	// region PERFORM
	modelLogs, appErr := a.storage.Log().Get(
		ctx,
		searchOpt,
		append(extractDefaultFiltersFromLogSearch(in), notStandartFilters...)...,
	)
	res.Page = int32(searchOpt.Page)
	if appErr != nil {
		if IsErrNoRows(appErr) {
			return &res, nil
		} else {
			return nil, appErr
		}
	}

	// endregion
	//region CONVERT TO RESPONSE
	for _, v := range modelLogs {
		protoLog, err := convertLogModelToMessage(v)
		if err != nil {
			return nil, err
		}
		rows = append(rows, protoLog)
	}
	if len(rows)-1 == searchOpt.Size {
		res.Next = true
		res.Items = rows[0 : len(rows)-1]
	} else {
		res.Items = rows
	}
	// endregion
	return &res, nil
}

func (a *App) SearchLogsByConfigId(ctx context.Context, in *proto.SearchLogByConfigIdRequest) (*proto.Logs, errors.AppError) {
	var (
		res               proto.Logs
		rows              []*proto.Log
		notDefaultFilters []model.Filter
	)
	searchOpt := ExtractSearchOptions(in)
	// region APPLYING FILTERS
	if x := in.GetUser(); x != nil {
		if v := x.GetId(); v != 0 {
			notDefaultFilters = append(notDefaultFilters, model.Filter{
				Column:         "user_id",
				Value:          v,
				ComparisonType: model.Equal,
			})
		} else if v := x.GetName(); v != "" {
			notDefaultFilters = append(notDefaultFilters, model.Filter{
				Column:         "coalesce(wbt_user.name::varchar, wbt_user.username::varchar)",
				Value:          strings.Replace(v, "*", "%", -1),
				ComparisonType: model.ILike,
			})
		}
	}
	// REQUIRED !
	notDefaultFilters = append(notDefaultFilters, model.Filter{
		Column:         "log.config_id",
		Value:          in.GetConfigId(),
		ComparisonType: 0,
	})
	// endregion
	// region PERFORM
	modelLogs, appErr := a.storage.Log().Get(
		ctx,
		searchOpt,
		append(extractDefaultFiltersFromLogSearch(in), notDefaultFilters...)...,
	)
	res.Page = int32(searchOpt.Page)
	if appErr != nil {
		if IsErrNoRows(appErr) {
			return &res, nil
		} else {
			return nil, appErr
		}
	}
	// endregion
	//region CONVERT TO RESPONSE
	for _, v := range modelLogs {
		protoLog, err := convertLogModelToMessage(v)
		if err != nil {
			return nil, err
		}
		rows = append(rows, protoLog)
	}
	if len(rows)-1 == searchOpt.Size {
		res.Next = true
		res.Items = rows[0 : len(rows)-1]
	} else {
		res.Items = rows
	}
	// endregion
	return &res, nil
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
	err = a.storage.Log().Insert(ctx, model)
	if err != nil {
		return err
	}

	return nil

}

func (a *App) InsertLogByRabbitMessageBulk(ctx context.Context, rabbitMessages []*model.RabbitMessage, domainId int64, objectName string) errors.AppError {
	searchResult, err := a.storage.Config().Get(ctx, nil, nil,
		model.Filter{
			Column:         "wbt_class.name",
			Value:          objectName,
			ComparisonType: model.Like,
		},
		model.Filter{
			Column:         "object_config.domain_id",
			Value:          domainId,
			ComparisonType: model.Equal,
		},
	)
	if err != nil {
		return err
	}
	config := searchResult[0]

	logs, err := convertRabbitMessageToModelBulk(rabbitMessages, config.Id)
	if err != nil {
		return err
	}
	err = a.storage.Log().InsertMany(ctx, *logs)
	if err != nil {
		return err
	}
	return nil

}

// endregion

// region UTIL FUNCTIONS
func convertLogModelToMessage(m *model.Log) (*proto.Log, errors.AppError) {
	log := &proto.Log{
		Id:     int32(m.Id),
		Action: m.Action,

		UserIp:   m.UserIp,
		NewState: string(m.NewState),
		ConfigId: int32(m.ConfigId),
	}
	if !m.User.IsZero() {
		log.User = &proto.Lookup{
			Id:   int32(m.User.Id.Int()),
			Name: m.User.Name.String(),
		}
	}
	if !m.Object.IsZero() {
		log.Object = &proto.Lookup{
			Id:   int32(m.Object.Id.Int()),
			Name: m.Object.Name.String(),
		}
	}
	if !m.Date.IsZero() {
		log.Date = m.Date.ToMilliseconds()
	}
	if m.RecordId != 0 {
		log.Record = &proto.Lookup{
			Id:   int32(m.RecordId),
			Name: "",
		}
	}
	return log, nil
}

func convertRabbitMessageToModel(m *model.RabbitMessage, configId int) (*model.Log, errors.AppError) {
	log := &model.Log{
		Action:   m.Action,
		Date:     (model.NullTime)(time.Unix(m.Date, 0)),
		UserIp:   m.UserIp,
		NewState: m.NewState,
		RecordId: m.RecordId,
		ConfigId: configId,
		User:     model.Lookup{Id: model.NewNullInt(m.UserId)},
	}
	// log.User = m.UserId)
	//if err != nil {
	//	return nil, errors.NewBadRequestError("app.log.convert_rabbit_message_to_model.scan.error", err.Error())
	//}

	return log, nil
}

func convertRabbitMessageToModelBulk(m []*model.RabbitMessage, configId int) (*[]*model.Log, errors.AppError) {
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

type LogSearch interface {
	GetDateFrom() int64
	GetDateTo() int64
	GetUserIp() string
	GetAction() proto.Action
}

func extractDefaultFiltersFromLogSearch(in LogSearch) []model.Filter {
	var result []model.Filter

	if in.GetDateFrom() != 0 {
		result = append(result, model.Filter{
			Column:         "log.date",
			Value:          time.Unix(in.GetDateFrom()/1000, 0).UTC(),
			ComparisonType: model.GreaterThanOrEqual,
		})
	}

	if in.GetAction() != proto.Action_default_no_action {
		result = append(result, model.Filter{
			Column:         "log.action",
			Value:          in.GetAction().String(),
			ComparisonType: model.ILike,
		})
	}

	if in.GetDateTo() != 0 {
		result = append(result, model.Filter{
			Column:         "log.date",
			Value:          time.Unix(in.GetDateTo()/1000, 0).UTC(),
			ComparisonType: model.LessThanOrEqual,
		})
	}
	if in.GetUserIp() != "" {
		result = append(result, model.Filter{
			Column:         "log.user_ip",
			Value:          in.GetUserIp(),
			ComparisonType: model.Equal,
		})
	}

	return result
}

// endregion
