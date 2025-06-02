package app

import (
	"context"
	storageModel "github.com/webitel/logger/api/storage"
	"github.com/webitel/logger/internal/auth"
	"github.com/webitel/logger/internal/model"
	notifier "github.com/webitel/webitel-go-kit/pkg/watcher"
	grpc "google.golang.org/grpc"
	"io"
	"log/slog"
	"time"
)

const (
	LogsNotifierObject = "log"
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

func (a *App) CreateLog(ctx context.Context, log *model.Log, domainId int) error {
	var (
		err error
	)
	defer func() {
		notifyErr := a.watcherManager.Notify(LogsNotifierObject, notifier.EventTypeCreate, NewNotifierLogArgs(err == nil, log))
		if notifyErr != nil {
			slog.ErrorContext(ctx, notifyErr.Error())
		}
	}()
	err = a.storage.Log().Insert(ctx, log, domainId)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) DeleteLogs(ctx context.Context, configId int, earlierThan time.Time) (int, error) {
	if earlierThan.IsZero() {
		earlierThan = time.Now()
	}
	affected, err := a.storage.Log().Delete(ctx, earlierThan, configId)
	if err != nil {
		return 0, err
	}
	return affected, err
}

func (a *App) UploadFile(ctx context.Context, domainId int, uuid string, storageId int, sFile io.Reader, metadata model.File) (*model.File, error) {
	stream, err := a.file.UploadFile(ctx)
	if err != nil {
		return nil, err
	}
	err = stream.Send(&storageModel.UploadFileRequest{
		Data: &storageModel.UploadFileRequest_Metadata_{
			Metadata: &storageModel.UploadFileRequest_Metadata{
				DomainId:  int64(domainId),
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

	defer func(stream grpc.ClientStreamingClient[storageModel.UploadFileRequest, storageModel.UploadFileResponse]) {
		_ = stream.CloseSend()
	}(stream)

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

// endregion

// region UTIL FUNCTIONS

type NotifierLogArgs struct {
	Log                *model.Log
	OperationSucceeded bool
}

func (n *NotifierLogArgs) GetArgs() map[string]any {
	return map[string]any{
		"object":    n.Log,
		"objclass":  LogsNotifierObject,
		"succeeded": n.OperationSucceeded,
	}

}

func NewNotifierLogArgs(success bool, log *model.Log) *NotifierLogArgs {
	return &NotifierLogArgs{OperationSucceeded: success, Log: log}
}

// endregion
