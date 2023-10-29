package app

import (
	"context"

	errors "github.com/webitel/engine/model"

	"github.com/webitel/engine/auth_manager"
	"github.com/webitel/logger/model"

	proto "github.com/webitel/protos/logger"
)

type LoggerService struct {
	proto.UnimplementedLoggerServiceServer
	app *App
}

func NewLoggerService(app *App) (*LoggerService, errors.AppError) {
	if app == nil {
		return nil, errors.NewInternalError("api.config.new_logger_service.args_check.app_nil", "app is nil")
	}
	return &LoggerService{app: app}, nil
}

func (s *LoggerService) SearchLogByRecordId(ctx context.Context, in *proto.SearchLogByRecordIdRequest) (*proto.Logs, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	//
	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	return s.app.SearchLogsByRecordId(ctx, in)
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
	return s.app.SearchLogsByUserId(ctx, in)
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

	return s.app.SearchLogsByConfigId(ctx, in)
}
