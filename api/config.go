package api

import (
	"context"
	"github.com/webitel/engine/auth_manager"
	"webitel_logger/app"
	"webitel_logger/proto"

	errors "github.com/webitel/engine/model"
)

type ConfigService struct {
	proto.UnimplementedConfigServiceServer
	app *app.App
}

func NewConfigService(app *app.App) (*ConfigService, errors.AppError) {
	if app == nil {
		return nil, errors.NewInternalError("api.config.new_config_service.args_check.app_nill", "app is nil")
	}
	return &ConfigService{app: app}, nil
}

func (s *ConfigService) GetConfigByObjectId(ctx context.Context, in *proto.GetConfigByObjectIdRequest) (*proto.Config, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	permission := session.GetPermission()
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}

	return s.app.GetConfigByObjectId(ctx, int(in.GetDomainId()), int(in.GetObjectId()))
}

func (s *ConfigService) GetAllConfigs(ctx context.Context, in *proto.GetAllConfigsRequest) (*proto.Configs, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	permission := session.GetPermission()
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	var out proto.Configs
	opt, err := app.ExtractSearchOptions(in)
	if err != nil {
		return nil, err
	}
	res, err := s.app.GetAllConfigs(ctx, opt, int(in.GetDomainId()))
	if err != nil {
		return nil, err
	}
	out.Configs = res
	return &out, nil
}

func (s *ConfigService) UpdateConfig(ctx context.Context, in *proto.UpdateConfigRequest) (*proto.Config, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	permission := session.GetPermission()
	if !permission.CanUpdate() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	return s.app.UpdateConfig(ctx, in)
}
