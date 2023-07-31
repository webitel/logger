package api

import (
	"context"
	"github.com/webitel/engine/auth_manager"
	"webitel_logger/app"
	"webitel_logger/model"
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
	opt, err := app.ExtractSearchOptions(in)
	if err != nil {
		return nil, err
	}
	rows, err := s.app.GetLogsByUserId(ctx, opt, int(in.GetUserId()))
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

//func (s *LoggerService) GetLogsByObjectId(ctx context.Context, in *proto.GetLogsByObjectIdRequest) (*proto.Logs, error) {
//	session, err := s.app.GetSessionFromCtx(ctx)
//	if err != nil {
//		return nil, err
//	}
//	//
//	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
//	if !permission.CanRead() {
//		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
//	}
//	var result proto.Logs
//	opt, err := app.ExtractSearchOptions(in)
//	if err != nil {
//		return nil, err
//	}
//	rows, err := s.app.GetLogsByObjectI(ctx, opt, int(session.DomainId), int(in.GetObjectId()))
//	if err != nil {
//		return nil, err
//	}
//	if int32(len(rows)-1) == in.Size {
//		result.Next = true
//	}
//	result.Items = rows
//	result.Page = in.GetPage()
//	return &result, nil
//}

func (s *LoggerService) GetLogsByConfigId(ctx context.Context, in *proto.GetLogsByConfigIdRequest) (*proto.Logs, error) {
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
	opt, err := app.ExtractSearchOptions(in)
	if err != nil {
		return nil, err
	}
	rows, err := s.app.GetLogsByConfigId(ctx, opt, int(in.GetConfigId()))
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
