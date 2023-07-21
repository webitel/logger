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
		return nil, errors.NewInternalError("api.config.new_logger_service.args_check.app_nil", "app is nil")
	}
	return &LoggerService{app: app}, nil
}

func (s *LoggerService) GetLogsByUserId(ctx context.Context, in *proto.GetLogsByUserIdRequest) (*proto.Logs, error) {
	var result *proto.Logs
	opt, err := app.ExtractSearchOptions(in)
	if err != nil {
		return nil, err
	}
	rows, err := s.app.GetLogsByUserId(ctx, opt, int(in.GetUserId()))
	if err != nil {
		return nil, err
	}
	result.Logs = rows
	return result, nil
}

func (s *LoggerService) GetLogsByObjectId(ctx context.Context, in *proto.GetLogsByObjectIdRequest) (*proto.Logs, error) {
	var result proto.Logs
	opt, err := app.ExtractSearchOptions(in)
	if err != nil {
		return nil, err
	}
	rows, err := s.app.GetLogsByObjectId(ctx, opt, int(in.GetDomainId()), int(in.GetObjectId()))
	if err != nil {
		return nil, err
	}
	result.Logs = rows
	return &result, nil
}
