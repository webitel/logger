package api

import (
	"context"
	"github.com/webitel/engine/auth_manager"
	"github.com/webitel/logger/app"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/proto"

	errors "github.com/webitel/engine/model"
)

//var _ proto.ConfigServiceServer = ConfigService{}

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

// ReadConfig selects config by id
func (s *ConfigService) ReadConfig(ctx context.Context, in *proto.ReadConfigRequest) (*proto.Config, error) {
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
	return s.app.GetConfigById(ctx, rbac, int(in.GetConfigId()))
}

func (s *ConfigService) ReadSystemObjects(ctx context.Context, request *proto.ReadSystemObjectsRequest) (*proto.SystemObjects, error) {

	// region AUTHORIZATION
	session, err := s.app.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	permission := session.GetPermission(model.PERMISSION_SCOPE_LOG)
	if !permission.CanRead() {
		return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_READ)
	}
	return s.app.GetSystemObjects(ctx, request, int(session.DomainId))
}

// ReadConfigByObjectId used for internal purpose with client, checks if config enabled
func (s *ConfigService) ReadConfigByObjectId(ctx context.Context, in *proto.ReadConfigByObjectIdRequest) (*proto.Config, error) {
	return s.app.GetConfigByObjectId(ctx, int(in.GetDomainId()), int(in.GetObjectId()))
}

// SearchConfig selects all configs by domainId
func (s *ConfigService) SearchConfig(ctx context.Context, in *proto.SearchConfigRequest) (*proto.Configs, error) {
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
	if session.UseRBAC(auth_manager.PERMISSION_ACCESS_READ, permission) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: auth_manager.PERMISSION_ACCESS_READ.Value(),
		}
	}
	return s.app.GetAllConfigs(ctx, rbac, int(session.DomainId), in)
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

// PatchConfig updates existing config
func (s *ConfigService) PatchConfig(ctx context.Context, in *proto.PatchConfigRequest) (*proto.Config, error) {
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

// CreateConfig inserts new config
func (s *ConfigService) CreateConfig(ctx context.Context, in *proto.CreateConfigRequest) (*proto.Config, error) {
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

// DeleteConfig deletes config by id
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
		access, err := s.app.ConfigCheckAccess(ctx, session.DomainId, int64(in.GetConfigId()), session.GetAclRoles(), auth_manager.PERMISSION_ACCESS_DELETE)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, s.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_DELETE)
		}
	}
	appErr := s.app.DeleteConfig(ctx, in.GetConfigId())
	if appErr != nil {
		return nil, appErr
	}
	return &proto.Empty{}, nil
}

// DeleteConfigBulk deletes configs by array of ids
func (s *ConfigService) DeleteConfigBulk(ctx context.Context, in *proto.DeleteConfigBulkRequest) (*proto.Empty, error) {
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
