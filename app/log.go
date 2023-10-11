package app

import (
	"context"

	"github.com/webitel/logger/model"

	"time"

	proto "github.com/webitel/protos/logger"

	errors "github.com/webitel/engine/model"
)

// region PERFORM ACTIONS
func (a *App) SearchLogsByUserId(ctx context.Context, in *proto.SearchLogByUserIdRequest) (*proto.Logs, errors.AppError) {
	var (
		res  proto.Logs
		rows []*proto.Log
		//filters []model.Filter
	)
	filters := model.FilterArray{
		Filters:    []*model.FilterBunch{},
		Connection: model.AND,
	}

	searchOpt := ExtractSearchOptions(in)
	// region APPLYING FILTERS

	if len(in.GetObjectId()) != 0 {
		newFilterBunch := model.FilterBunch{
			ConnectionType: model.OR,
		}
		for _, v := range in.GetObjectId() {
			newFilterBunch.Bunch = append(newFilterBunch.Bunch, &model.Filter{
				Column:         "object_config.object_id",
				Value:          v,
				ComparisonType: model.Equal,
			})

		}
		filters.Filters = append(filters.Filters, &newFilterBunch)

	}

	// REQUIRED !
	requiredFilterBunch := model.FilterBunch{
		ConnectionType: model.AND,
	}
	requiredFilterBunch.Bunch = append(requiredFilterBunch.Bunch, &model.Filter{
		Column:         "user_id",
		Value:          in.GetUserId(),
		ComparisonType: model.Equal,
	})

	filters.Filters = append(filters.Filters, extractDefaultFiltersFromLogSearch(in)...)
	filters.Filters = append(filters.Filters, &requiredFilterBunch)
	// endregion
	// region PERFORM
	modelLogs, appErr := a.storage.Log().Get(
		ctx,
		searchOpt,
		filters,
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
		res  proto.Logs
		rows []*proto.Log
		//notDefaultFilters []model.Filter
	)
	filters := model.FilterArray{
		Connection: model.AND,
	}
	searchOpt := ExtractSearchOptions(in)
	newFilterBunch := model.FilterBunch{
		ConnectionType: model.OR,
	}
	// region APPLYING FILTERS
	for _, v := range in.GetUserId() {
		newFilterBunch.Bunch = append(newFilterBunch.Bunch, &model.Filter{
			Column:         "user_id",
			Value:          v,
			ComparisonType: model.Equal,
		})
	}
	if len(newFilterBunch.Bunch) != 0 {
		filters.Filters = append(filters.Filters, &newFilterBunch)
	}
	// REQUIRED !
	requiredFilterBunch := model.FilterBunch{
		ConnectionType: model.AND,
	}
	requiredFilterBunch.Bunch = append(requiredFilterBunch.Bunch, &model.Filter{
		Column:         "log.config_id",
		Value:          in.GetConfigId(),
		ComparisonType: 0,
	})
	filters.Filters = append(filters.Filters, extractDefaultFiltersFromLogSearch(in)...)
	filters.Filters = append(filters.Filters, &requiredFilterBunch)
	// endregion
	// region PERFORM
	modelLogs, appErr := a.storage.Log().Get(
		ctx,
		searchOpt,
		filters,
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

func (a *App) SearchLogsByRecordId(ctx context.Context, in *proto.SearchLogByRecordIdRequest) (*proto.Logs, errors.AppError) {
	var (
		res  proto.Logs
		rows []*proto.Log
		//notDefaultFilters []model.Filter
	)
	filters := model.FilterArray{
		Connection: model.AND,
	}
	searchOpt := ExtractSearchOptions(in)
	userFilterBunch := model.FilterBunch{
		ConnectionType: model.OR,
	}
	// region APPLYING FILTERS

	// region NOT REQUIRED !
	// multiselect user filter
	// [OR] connection between multiselect
	for _, v := range in.GetUserId() {
		userFilterBunch.Bunch = append(userFilterBunch.Bunch, &model.Filter{
			Column:         "user_id",
			Value:          v,
			ComparisonType: model.Equal,
		})
	}
	if len(userFilterBunch.Bunch) != 0 {
		filters.Filters = append(filters.Filters, &userFilterBunch)
	}
	// endregion

	// region REQUIRED !
	// [AND] connection between required filters
	requiredFilterBunch := model.FilterBunch{
		ConnectionType: model.AND,
	}
	// object filter
	requiredFilterBunch.Bunch = append(requiredFilterBunch.Bunch, &model.Filter{
		Column:         "log.object_name",
		Value:          in.GetObject().String(),
		ComparisonType: model.Equal,
	})
	// record id filter
	requiredFilterBunch.Bunch = append(requiredFilterBunch.Bunch, &model.Filter{
		Column:         "log.record_id",
		Value:          in.GetRecordId(),
		ComparisonType: model.Equal,
	})

	filters.Filters = append(filters.Filters, extractDefaultFiltersFromLogSearch(in)...)
	filters.Filters = append(filters.Filters, &requiredFilterBunch)
	// endregion

	// endregion

	// region PERFORM
	modelLogs, appErr := a.storage.Log().Get(
		ctx,
		searchOpt,
		filters,
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
	searchResult, err := a.storage.Config().Get(ctx, nil, nil, model.FilterBunch{
		Bunch: []*model.Filter{
			{
				Column:         "wbt_class.name",
				Value:          objectName,
				ComparisonType: model.Like,
			},
			{
				Column:         "object_config.domain_id",
				Value:          domainId,
				ComparisonType: model.Equal,
			},
		},
		ConnectionType: model.AND,
	})
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
	if s := m.Record.Id.Int32(); s != 0 {
		log.Record = &proto.Lookup{
			Id:   s,
			Name: m.Record.Name.String(),
		}
	}
	//if s := m.Record.Name.String(); s != 0 {
	//	log.Record = &proto.Lookup{
	//		Id:  s,
	//
	//	}
	//}
	return log, nil
}

func convertRabbitMessageToModel(m *model.RabbitMessage, configId int) (*model.Log, errors.AppError) {
	log := &model.Log{
		Action:   m.Action,
		Date:     (model.NullTime)(time.Unix(m.Date, 0)),
		UserIp:   m.UserIp,
		NewState: m.NewState,
		ConfigId: configId,
		Object:   model.Lookup{Name: model.NewNullString(m.Schema)},
	}
	userId, err := model.NewNullInt(m.UserId)
	if err != nil {
		return nil, errors.NewInternalError("app.log.convert_rabbit_message.convert_to_null_user.error", err.Error())
	}
	log.User = model.Lookup{Id: userId}
	recordId, err := model.NewNullInt(m.RecordId)
	if err != nil {
		return nil, errors.NewInternalError("app.log.convert_rabbit_message.convert_to_null_record.error", err.Error())
	}
	log.Record = model.Lookup{Id: recordId}

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
	GetAction() []proto.Action
}

func extractDefaultFiltersFromLogSearch(in LogSearch) []*model.FilterBunch {
	var result []*model.FilterBunch
	andBunch := &model.FilterBunch{
		ConnectionType: model.AND,
	}
	orBunch := &model.FilterBunch{
		ConnectionType: model.OR,
	}
	if in.GetDateFrom() != 0 {
		andBunch.Bunch = append(andBunch.Bunch, &model.Filter{
			Column:         "log.date",
			Value:          time.Unix(in.GetDateFrom()/1000, 0).UTC(),
			ComparisonType: model.GreaterThanOrEqual,
		})
	}
	for _, v := range in.GetAction() {
		if v != proto.Action_default_no_action {
			orBunch.Bunch = append(orBunch.Bunch, &model.Filter{
				Column:         "log.action",
				Value:          v.String(),
				ComparisonType: model.ILike,
			})
		}
	}

	if in.GetDateTo() != 0 {
		andBunch.Bunch = append(andBunch.Bunch, &model.Filter{
			Column:         "log.date",
			Value:          time.Unix(in.GetDateTo()/1000, 0).UTC(),
			ComparisonType: model.LessThanOrEqual,
		})
	}
	if in.GetUserIp() != "" {
		andBunch.Bunch = append(andBunch.Bunch, &model.Filter{
			Column:         "log.user_ip",
			Value:          in.GetUserIp(),
			ComparisonType: model.Equal,
		})
	}
	if len(orBunch.Bunch) != 0 {
		result = append(result, orBunch)
	}

	if len(andBunch.Bunch) != 0 {
		result = append(result, andBunch)
	}

	return result
}

// endregion
