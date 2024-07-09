package app

import (
	storageModel "buf.build/gen/go/webitel/storage/protocolbuffers/go"
	"context"
	"io"

	"github.com/webitel/logger/model"

	"time"

	proto "buf.build/gen/go/webitel/logger/protocolbuffers/go"
)

// region PERFORM ACTIONS

func (a *App) SearchLogs(ctx context.Context, searchOpt *model.SearchOptions, filters *model.LogFilters) ([]*model.Log, model.AppError) {
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

func (a *App) UploadFile(ctx context.Context, domainId int64, uuid string, storageId int, sFile io.Reader, metadata model.File) (*model.File, model.AppError) {
	stream, err := a.file.UploadFile(ctx)
	if err != nil {
		return nil, model.NewInternalError("app.log.upload_file.request_stream.error", err.Error())
	}

	err = stream.Send(&storageModel.UploadFileRequest{
		Data: &storageModel.UploadFileRequest_Metadata_{
			Metadata: &storageModel.UploadFileRequest_Metadata{
				DomainId:  domainId,
				Name:      metadata.Name,
				MimeType:  metadata.MimeType,
				Uuid:      uuid,
				ProfileId: int64(storageId),
			},
		},
	})

	if err != nil {
		return nil, model.NewInternalError("app.log.upload_file.send_metadata.error", err.Error())
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
		err = stream.Send(&storageModel.UploadFileRequest{
			Data: &storageModel.UploadFileRequest_Chunk{
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
		return nil, model.NewInternalError("app.log.upload_file.send_stream.error", err.Error())
	}

	var res *storageModel.UploadFileResponse
	res, err = stream.CloseAndRecv()
	if err != nil {
		return nil, model.NewInternalError("app.log.upload_file.close_stream.error", err.Error())
	}

	metadata.Id = int(res.FileId)
	metadata.Size = res.Size
	metadata.Url = res.FileUrl
	metadata.PublicUrl = res.Server + res.FileUrl

	return &metadata, nil
}

func (a *App) InsertLogByRabbitMessage(ctx context.Context, rabbitMessage *model.RabbitMessage, domainId, objectId int) model.AppError {

	model, err := convertRabbitMessageToModel(rabbitMessage)
	if err != nil {
		return err
	}
	err = a.storage.Log().Insert(ctx, model, domainId)
	if err != nil {
		return err
	}

	return nil

}

func (a *App) InsertLogByRabbitMessageBulk(ctx context.Context, rabbitMessages []*model.RabbitMessage, domainId int64) model.AppError {
	logs, err := convertRabbitMessageToModelBulk(rabbitMessages)
	if err != nil {
		return err
	}
	err = a.storage.Log().InsertBulk(ctx, logs, int(domainId))
	if err != nil {
		return err
	}
	return nil

}

// endregion

// region UTIL FUNCTIONS
func convertLogModelToMessage(m *model.Log) (*proto.Log, model.AppError) {
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
	//if s := m.LogEntity.Name.String(); s != 0 {
	//	log.LogEntity = &proto.Lookup{
	//		Id:  s,
	//
	//	}
	//}
	return log, nil
}

func convertLogModelToMessageBulk(m []*model.Log) ([]*proto.Log, model.AppError) {
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

func convertRabbitMessageToModel(m *model.RabbitMessage) (*model.Log, model.AppError) {
	log := &model.Log{
		Action:   m.Action,
		Date:     (model.NullTime)(time.Unix(m.Date, 0)),
		UserIp:   m.UserIp,
		NewState: m.NewState,
		Object:   model.Lookup{Name: model.NewNullString(m.Schema)},
	}
	userId, err := model.NewNullInt(m.UserId)
	if err != nil {
		return nil, model.NewInternalError("app.log.convert_rabbit_message.convert_to_null_user.error", err.Error())
	}
	log.User = model.Lookup{Id: userId}
	recordId, err := model.NewNullInt(m.RecordId)
	if err != nil {
		return nil, model.NewInternalError("app.log.convert_rabbit_message.convert_to_null_record.error", err.Error())
	}
	log.Record = model.Lookup{Id: recordId}

	return log, nil
}

func convertRabbitMessageToModelBulk(m []*model.RabbitMessage) ([]*model.Log, model.AppError) {
	var logs []*model.Log
	for _, v := range m {
		log, err := convertRabbitMessageToModel(v)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

type LogSearcher interface {
	GetDateFrom() int64
	GetDateTo() int64
	GetUserIp() string
	GetAction() []proto.Action
}
