package api

import (
	"context"
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

func (s *ConfigService) GetByObjectId(ctx context.Context, in *proto.Object) (*proto.Config, error) {
	return s.app.GetConfigByObjectId(ctx, int(in.GetObjectId()))
}

func (s *ConfigService) Update(ctx context.Context, in *proto.Config) (*proto.Config, error) {
	return s.app.UpdateConfig(ctx, in)
}
