package app

import (
	"context"
	"github.com/webitel/engine/auth_manager"
	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"
	proto "github.com/webitel/protos/logger"
)

func NewSchemaVersionService(app *App) (*SchemaVersionService, errors.AppError) {
	if app == nil {
		return nil, errors.NewInternalError("api.config.new_config_service.args_check.app_nill", "app is nil")
	}
	return &SchemaVersionService{app: app}, nil
}

type SchemaVersionService struct {
	app *App
	proto.UnimplementedSchemaVersionsServiceServer
}

func (s *SchemaVersionService) Search(ctx context.Context, in *proto.SearchSchemaVersionRequest) (*proto.SchemaVersions, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	//
	permission := session.GetPermission(model.PERMISSION_SCOPE_SCHEMA)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	return s.app.GetFlowVersions(ctx, in)
}
