package api

import (
	"webitel_logger/app"
	"webitel_logger/proto"
	"webitel_logger/storage"

	errors "github.com/webitel/engine/model"
	"google.golang.org/grpc"
)

func BuildGrpc(store storage.Storage, app *app.App) (*grpc.Server, errors.AppError) {

	grpcServer := grpc.NewServer()
	// * Creating services
	l, appErr := NewLoggerService(app)
	if appErr != nil {
		return nil, appErr
	}
	c, appErr := NewConfigService(app)
	if appErr != nil {
		return nil, appErr
	}

	// * register logger service
	proto.RegisterLoggerServiceServer(grpcServer, l)
	// * register config service
	proto.RegisterConfigServiceServer(grpcServer, c)

	return grpcServer, nil
}
