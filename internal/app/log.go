package app

import (
	"context"
	storageModel "github.com/webitel/logger/api/storage"
	authmodel "github.com/webitel/logger/internal/auth/model"
	"io"

	"github.com/webitel/logger/internal/model"

	"time"
)

// region PERFORM ACTIONS

func (a *App) SearchLogs(ctx context.Context, searchOpt *model.SearchOptions, filters *model.LogFilters) ([]*model.Log, model.AppError) {

	session, err := a.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Read
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, a.MakeScopeError(session, scope, accessMode)
	}
	// region PERFORM
	modelLogs, appErr := a.storage.Log().Select(
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

func (a *App) DeleteLogs(ctx context.Context, configId int, olderThan time.Time) (int, model.AppError) {
	session, err := a.AuthorizeFromContext(ctx)
	if err != nil {
		return 0, err
	}

	// OBAC check
	accessMode := authmodel.Edit
	secondaryMode := authmodel.Delete
	scope := model.ScopeLog
	// Edit or Delete permissions allow this operation
	if !session.HasObacAccess(scope, accessMode) && !session.HasObacAccess(scope, secondaryMode) {
		return 0, a.MakeScopeError(session, scope, accessMode)
	}
	// RBAC check
	if session.UseRbacAccess(scope, accessMode) && session.UseRbacAccess(scope, secondaryMode) {
		rbacAccess, err := a.storage.Config().CheckAccess(ctx, session.GetDomainId(), int64(configId), session.GetAclRoles(), uint8(accessMode))
		if err != nil {
			return 0, err
		}
		secondaryRbacAccess, err := a.storage.Config().CheckAccess(ctx, session.GetDomainId(), int64(configId), session.GetAclRoles(), uint8(secondaryMode))
		if err != nil {
			return 0, err
		}
		if !rbacAccess && !secondaryRbacAccess {
			return 0, a.MakeScopeError(session, scope, accessMode)
		}
	}

	if olderThan.IsZero() {
		olderThan = time.Now()
	}
	return a.storage.Log().Delete(ctx, olderThan, configId)
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
