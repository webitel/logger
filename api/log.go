package api

import (
	"context"
	"github.com/webitel/engine/auth_manager"
	"github.com/webitel/logger/app"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/proto"

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

func (s *LoggerService) SearchLogByUserId(ctx context.Context, in *proto.SearchLogByUserIdRequest) (*proto.Logs, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	//
	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	var result proto.Logs
	rows, err := s.app.SearchLogsByUserId(ctx, in)
	if err != nil {
		return nil, err
	}
	if int32(len(rows)-1) == in.Size {
		result.Next = true
	}
	result.Items = rows
	result.Page = in.GetPage()
	return &result, nil
}

func (s *LoggerService) SearchLogByConfigId(ctx context.Context, in *proto.SearchLogByConfigIdRequest) (*proto.Logs, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	//
	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	var result proto.Logs
	rows, err := s.app.SearchLogsByConfigId(ctx, in)
	if err != nil {
		return nil, err
	}
	if int32(len(rows)-1) == in.Size {
		result.Next = true
	}
	result.Items = rows
	result.Page = in.GetPage()
	return &result, nil
}
