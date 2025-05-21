package app

import (
	"context"
	storageModel "github.com/webitel/logger/api/storage"
	"github.com/webitel/logger/internal/auth"
	"io"

	"github.com/webitel/logger/internal/model"

	"time"
)

// region PERFORM ACTIONS

func (a *App) SearchLogs(ctx context.Context, searchOpt *model.SearchOptions, filters *model.LogFilters) ([]*model.Log, error) {

	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Read)
	if err != nil {
		return nil, err
	}
	// OBAC check
	if !session.CheckObacAccess() {
		return nil, a.MakeScopeError(session.GetMainObjClassName())
	}
	// region PERFORM
	modelLogs, err := a.storage.Log().Select(
		ctx,
		searchOpt,
		filters.ExtractFilters(),
	)
	if err != nil {
		return nil, err
	}
	// endregion
	return modelLogs, nil
}

func (a *App) UploadFile(ctx context.Context, domainId int64, uuid string, storageId int, sFile io.Reader, metadata model.File) (*model.File, error) {
	stream, err := a.file.UploadFile(ctx)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	var res *storageModel.UploadFileResponse
	res, err = stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	metadata.Id = int(res.FileId)
	metadata.Size = res.Size
	metadata.Url = res.FileUrl
	metadata.PublicUrl = res.Server + res.FileUrl

	return &metadata, nil
}

func (a *App) DeleteLogs(ctx context.Context, configId int, earlierThan time.Time) (int, error) {
	session, err := a.AuthorizeFromContext(ctx, model.ScopeLog, auth.Edit)
	if err != nil {
		return 0, err
	}

	// OBAC check
	// Edit or Delete permissions allow this operation
	if !session.CheckObacAccess() {
		return 0, a.MakeScopeError(session.GetMainObjClassName())
	}
	// RBAC check
	if session.IsRbacCheckRequired() {
		rbacAccess, err := a.storage.Config().CheckAccess(ctx, session.GetDomainId(), int64(configId), session.GetRoles(), session.GetMainAccessMode().Value())
		if err != nil {
			return 0, err
		}
		if !rbacAccess {
			return 0, a.MakeScopeError(session.GetMainObjClassName())
		}
	}

	if earlierThan.IsZero() {
		earlierThan = time.Now()
	}
	return a.storage.Log().Delete(ctx, earlierThan, configId)
}

func (a *App) InsertLogByRabbitMessage(ctx context.Context, rabbitMessage *model.RabbitMessage, domainId, objectId int) error {

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

func (a *App) InsertLogByRabbitMessageBulk(ctx context.Context, rabbitMessages []*model.RabbitMessage, domainId int64) error {
	logs, err := convertRabbitMessageToModelBulk(rabbitMessages)
	if err != nil {
		return err
	}
	_, err = a.storage.Log().InsertBulk(ctx, logs, int(domainId))
	if err != nil {
		return err
	}
	return nil

}

// endregion

// region UTIL FUNCTIONS

func convertRabbitMessageToModel(m *model.RabbitMessage) (*model.Log, error) {
	log := &model.Log{
		Action:   m.Action,
		Date:     (model.NullTime)(time.Unix(m.Date, 0)),
		UserIp:   m.UserIp,
		NewState: m.NewState,
		Object:   model.Lookup{Name: model.NewNullString(m.Schema)},
	}
	userId, err := model.NewNullInt(m.UserId)
	if err != nil {
		return nil, err
	}
	log.User = model.Lookup{Id: userId}
	recordId, err := model.NewNullInt(m.RecordId)
	if err != nil {
		return nil, err
	}
	log.Record = model.Lookup{Id: recordId}

	return log, nil
}

func convertRabbitMessageToModelBulk(m []*model.RabbitMessage) ([]*model.Log, error) {
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
