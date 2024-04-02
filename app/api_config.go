package app

import (
	"context"

	proto "buf.build/gen/go/webitel/logger/protocolbuffers/go"
	"github.com/webitel/engine/auth_manager"
	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"
)

type ConfigService struct {
	app *App
}

func NewConfigService(app *App) (*ConfigService, errors.AppError) {
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
	resModel, err := s.app.GetConfigById(ctx, rbac, int(in.GetConfigId()))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

func (s *ConfigService) ReadSystemObjects(ctx context.Context, request *proto.ReadSystemObjectsRequest) (*proto.SystemObjects, error) {

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

// ReadConfigByObjectId used for internal purpose
func (s *ConfigService) ReadConfigByObjectId(ctx context.Context, in *proto.ReadConfigByObjectIdRequest) (*proto.Config, error) {
	resModel, err := s.app.GetConfigByObjectId(ctx, int(in.GetDomainId()), int(in.GetObjectId()))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

// ReadConfigByObjectId used for internal purpose with client, checks if config enabled
func (s *ConfigService) CheckConfigStatus(ctx context.Context, in *proto.CheckConfigStatusRequest) (*proto.ConfigStatus, error) {
	isEnabled, err := s.app.CheckConfigStatus(ctx, in.GetObjectName(), in.GetDomainId())
	if err != nil {
		return nil, err
	}
	return &proto.ConfigStatus{IsEnabled: isEnabled}, nil

}

// SearchConfig selects all configs by domainId
func (s *ConfigService) SearchConfig(ctx context.Context, in *proto.SearchConfigRequest) (*proto.Configs, error) {
	var (
		rbac *model.RbacOptions
		res  proto.Configs
	)
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
	resModels, err := s.app.GetAllConfigs(ctx, rbac, ExtractSearchOptions(in), session.DomainId)
	if err != nil {
		if IsErrNoRows(err) {
			return &res, nil
		} else {
			return nil, err
		}
	}
	res.Next, res.Items, err = CalculateListResultMetadata[*model.Config, *proto.Config](in, resModels, ConvertConfigModelToMessage)
	if err != nil {
		return nil, err
	}
	res.Page = in.GetPage()

	return &res, nil
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
	mod, err := ConvertUpdateConfigMessageToModel(in, session.DomainId)
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.UpdateConfig(ctx, mod, int(session.UserId))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)

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
	updatedConfigModel, err := ConvertPatchConfigMessageToModel(in, session.DomainId)
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.PatchUpdateConfig(ctx, updatedConfigModel, in.GetFields(), int(session.UserId))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
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
	model, err := ConvertCreateConfigMessageToModel(in, session.GetDomainId())
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.InsertConfig(ctx, model, int(session.UserId))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
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
