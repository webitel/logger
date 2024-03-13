package app

import (
	"context"
	"io"

	strg "github.com/webitel/logger/api/storage"

	"github.com/webitel/logger/model"

	"time"

	proto "github.com/webitel/logger/api/native"

	errors "github.com/webitel/engine/model"
)

// region PERFORM ACTIONS

func (a *App) SearchLogs(ctx context.Context, searchOpt *model.SearchOptions, filters *model.LogFilters) ([]*model.Log, errors.AppError) {
	// region PERFORM
	modelLogs, appErr := a.storage.Log().Get(
		ctx,
		searchOpt,
		filters.ExtractFilters(),
	)
	if appErr != nil {
		if IsErrNoRows(appErr) {
			return modelLogs, nil
		} else {
			return nil, appErr
		}
	}
	// endregion
	return modelLogs, nil
}

func (a *App) UploadFile(ctx context.Context, domainId int64, uuid string, storageId int, sFile io.Reader, metadata model.File) (*model.File, errors.AppError) {
	stream, err := a.file.UploadFile(ctx)
	if err != nil {
		return nil, errors.NewInternalError("app.log.upload_file.request_stream.error", err.Error())
	}

	err = stream.Send(&strg.UploadFileRequest{
		Data: &strg.UploadFileRequest_Metadata_{
			Metadata: &strg.UploadFileRequest_Metadata{
				DomainId:  domainId,
				Name:      metadata.Name,
				MimeType:  metadata.MimeType,
				Uuid:      uuid,
				ProfileId: int64(storageId),
			},
		},
	})

	if err != nil {
		return nil, errors.NewInternalError("app.log.upload_file.send_metadata.error", err.Error())
	}

	defer stream.CloseSend()

	buf := make([]byte, 4*1024)
	var n int
	for {
		n, err = sFile.Read(buf)
		buf = buf[:n]
		if err != nil {
			break
		}
		err = stream.Send(&strg.UploadFileRequest{
			Data: &strg.UploadFileRequest_Chunk{
				Chunk: buf,
			},
		})
		if err != nil {
			break
		}
	}

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		return nil, errors.NewInternalError("app.log.upload_file.send_stream.error", err.Error())
	}

	var res *strg.UploadFileResponse
	res, err = stream.CloseAndRecv()
	if err != nil {
		return nil, errors.NewInternalError("app.log.upload_file.close_stream.error", err.Error())
	}

	metadata.Id = int(res.FileId)
	metadata.Size = res.Size
	metadata.Url = res.FileUrl
	metadata.PublicUrl = res.Server + res.FileUrl

	return &metadata, nil
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
	err = a.storage.Log().Insert(ctx, model, domainId)
	if err != nil {
		return err
	}

	return nil

}

func (a *App) InsertLogByRabbitMessageBulk(ctx context.Context, rabbitMessages []*model.RabbitMessage, domainId int64, objectName string) errors.AppError {
	//searchResult, err := a.storage.Config().Get(ctx, nil, nil, model.FilterBunch{
	//	Bunch: []*model.Filter{
	//		{
	//			Column:         "wbt_class.name",
	//			Value:          objectName,
	//			ComparisonType: model.Like,
	//		},
	//		{
	//			Column:         "object_config.domain_id",
	//			Value:          domainId,
	//			ComparisonType: model.Equal,
	//		},
	//	},
	//	ConnectionType: model.AND,
	//})
	//if err != nil {
	//	return err
	//}
	//config := searchResult[0]

	logs, err := convertRabbitMessageToModelBulk(rabbitMessages, 0)
	if err != nil {
		return err
	}
	err = a.storage.Log().InsertMany(ctx, *logs, int(domainId))
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

func convertLogModelToMessageBulk(m []*model.Log) ([]*proto.Log, errors.AppError) {
	var rows []*proto.Log
	for _, v := range m {
		protoLog, err := convertLogModelToMessage(v)
		if err != nil {
			return nil, err
		}
		rows = append(rows, protoLog)
	}
	return rows, nil
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

type LogSearcher interface {
	GetDateFrom() int64
	GetDateTo() int64
	GetUserIp() string
	GetAction() []proto.Action
}
