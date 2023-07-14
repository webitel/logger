package api

import (
	"context"
	"webitel_logger/app"
	"webitel_logger/proto"

	errors "github.com/webitel/engine/model"
)

type LoggerService struct {
	proto.UnimplementedLoggerServiceServer
	app *app.App
}

func NewLoggerService(app *app.App) (*LoggerService, errors.AppError) {
	if app == nil {
		return nil, errors.NewInternalError("api.config.new_logger_service.args_check.app_nill", "app is nil")
	}
	return &LoggerService{app: app}, nil
}

func (s *LoggerService) GetByUserId(ctx context.Context, in *proto.User) (*proto.Logs, error) {
	var result *proto.Logs

	rows, err := s.app.GetLogsByUserId(ctx, int(in.GetUserId()))
	if err != nil {
		return nil, err
	}
	result.Logs = rows
	return result, nil
}

func (s *LoggerService) GetByObjectId(ctx context.Context, in *proto.Object) (*proto.Logs, error) {
	var result *proto.Logs

	rows, err := s.app.GetLogsByObjectId(ctx, int(in.GetDomainId()), int(in.GetObjectId()))
	if err != nil {
		return nil, err
	}
	result.Logs = rows
	return result, nil
}

// func (s *Server) Update(ctx context.Context, in *proto.Config) (*proto.Config, error) {
// 	return nil, nil
// }

// func (s *Server) GetByObjectId(ctx context.Context, in *proto.Object) (*proto.Config, error) {
// 	return nil, nil
// }
