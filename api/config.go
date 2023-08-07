package api

import (
	"context"
	"github.com/webitel/engine/auth_manager"
	"logger/app"
	"logger/model"
	"logger/proto"

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

// GetConfigById selects config by id
func (s *ConfigService) GetConfigById(ctx context.Context, in *proto.GetConfigByIdRequest) (*proto.Config, error) {
	var rbac *model.RbacOptions
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	if session.UseRBAC(auth_manager.PERMISSION_ACCESS_READ, permission) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: auth_manager.PERMISSION_ACCESS_READ.Value(),
		}
	}
	return s.app.GetConfigById(ctx, rbac, int(in.GetId()))
}

// For internal purpose when check is config enabled
func (s *ConfigService) GetConfigByObjectId(ctx context.Context, in *proto.GetConfigByObjectIdRequest) (*proto.Config, error) {
	return s.app.GetConfigByObjectId(ctx, int(in.GetDomainId()), int(in.GetObjectId()))
}

// GetAllConfigs selects all configs by domainId
func (s *ConfigService) GetAllConfigs(ctx context.Context, in *proto.GetAllConfigsRequest) (*proto.Configs, error) {
	var rbac *model.RbacOptions
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	//
	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	var out proto.Configs
	opt, err := app.ExtractSearchOptions(in)
	if err != nil {
		return nil, err
	}
	if session.UseRBAC(auth_manager.PERMISSION_ACCESS_READ, permission) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: auth_manager.PERMISSION_ACCESS_READ.Value(),
		}
	}
	rows, err := s.app.GetAllConfigs(ctx, opt, rbac, int(session.DomainId))
	if err != nil {
		return nil, err
	}
	if int32(len(rows)-1) == in.Size {
		out.Next = true
	}
	out.Items = rows
	out.Page = int32(opt.Page)
	return &out, nil
}

// UpdateConfig updates existing config
func (s *ConfigService) UpdateConfig(ctx context.Context, in *proto.UpdateConfigRequest) (*proto.Config, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	if !permission.CanUpdate() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_UPDATE)
	}
	if session.UseRBAC(auth_manager.PERMISSION_ACCESS_UPDATE, permission) {
		access, err := s.app.ConfigCheckAccess(ctx, session.DomainId, int64(in.GetConfigId()), session.GetAclRoles(), auth_manager.PERMISSION_ACCESS_UPDATE)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_UPDATE)
		}
	}
	return s.app.UpdateConfig(ctx, in, int(session.DomainId), int(session.UserId))
}

// UpdateConfig updates existing config
func (s *ConfigService) PatchUpdateConfig(ctx context.Context, in *proto.PatchUpdateConfigRequest) (*proto.Config, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	if !permission.CanUpdate() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_UPDATE)
	}
	if session.UseRBAC(auth_manager.PERMISSION_ACCESS_UPDATE, permission) {
		access, err := s.app.ConfigCheckAccess(ctx, session.DomainId, int64(in.GetConfigId()), session.GetAclRoles(), auth_manager.PERMISSION_ACCESS_UPDATE)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_UPDATE)
		}
	}
	return s.app.PatchUpdateConfig(ctx, in, int(session.DomainId), int(session.UserId))
}

// InsertConfig inserts new config
func (s *ConfigService) InsertConfig(ctx context.Context, in *proto.InsertConfigRequest) (*proto.Config, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	if !permission.CanCreate() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_CREATE)
	}
	return s.app.InsertConfig(ctx, in, int(session.DomainId), int(session.UserId))
}

// Delete deletes config
func (s *ConfigService) DeleteConfig(ctx context.Context, in *proto.DeleteConfigRequest) (*proto.Empty, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanDelete() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_DELETE)
	}
	if session.UseRBAC(auth_manager.PERMISSION_ACCESS_DELETE, permission) {
		access, err := s.app.ConfigCheckAccess(ctx, session.DomainId, int64(in.GetId()), session.GetAclRoles(), auth_manager.PERMISSION_ACCESS_DELETE)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_DELETE)
		}
	}
	appErr := s.app.DeleteConfig(ctx, in.GetId())
	if appErr != nil {
		return nil, appErr
	}
	return &proto.Empty{}, nil
}

// InsertConfig inserts new config
func (s *ConfigService) DeleteConfigs(ctx context.Context, in *proto.DeleteConfigsRequest) (*proto.Empty, error) {
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	if !permission.CanDelete() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_DELETE)
	}
	var rbac *model.RbacOptions
	if session.UseRBAC(auth_manager.PERMISSION_ACCESS_DELETE, permission) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: auth_manager.PERMISSION_ACCESS_DELETE.Value(),
		}
	}
	appErr := s.app.DeleteConfigs(ctx, rbac, in.GetIds())
	if appErr != nil {
		return nil, appErr
	}
	return &proto.Empty{}, nil
}
