package app

import (
	"context"
	authmodel "github.com/webitel/logger/auth/model"

	proto "buf.build/gen/go/webitel/logger/protocolbuffers/go"
	"github.com/webitel/logger/model"
)

type ConfigService struct {
	app *App
}

func NewConfigService(app *App) (*ConfigService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_config_service.args_check.app_nil", "app is nil")
	}
	return &ConfigService{app: app}, nil
}

// ReadConfig selects config by id
func (s *ConfigService) ReadConfig(ctx context.Context, in *proto.ReadConfigRequest) (*proto.Config, error) {
	var rbac *model.RbacOptions
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, authmodel.Read) {
		return nil, s.app.MakeScopeError(session, scope, authmodel.Read)
	}
	// RBAC check
	if scope.IsRbacUsed() {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: authmodel.Read.Value(),
		}
	}
	resModel, err := s.app.GetConfigById(ctx, rbac, int(in.GetConfigId()))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

func (s *ConfigService) ReadSystemObjects(ctx context.Context, request *proto.ReadSystemObjectsRequest) (*proto.SystemObjects, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, authmodel.Read) {
		return nil, s.app.MakeScopeError(session, scope, authmodel.Read)
	}
	return s.app.GetSystemObjects(ctx, request, int(session.GetDomainId()))
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
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, authmodel.Read) {
		return nil, s.app.MakeScopeError(session, scope, authmodel.Read)
	}
	// RBAC check
	if scope.IsRbacUsed() {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: authmodel.Read.Value(),
		}
	}

	resModels, err := s.app.GetAllConfigs(ctx, rbac, ExtractSearchOptions(in), session.GetDomainId())
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
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	scope := session.GetScope(model.ScopeLog)
	accessMode := authmodel.Edit
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	// RBAC check
	if scope.IsRbacUsed() {
		access, err := s.app.ConfigCheckAccess(ctx, session.GetDomainId(), int64(in.GetConfigId()), session.GetAclRoles(), accessMode)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, s.app.MakeScopeError(session, scope, accessMode)
		}
	}
	mod, err := ConvertUpdateConfigMessageToModel(in, session.GetDomainId())
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.UpdateConfig(ctx, mod, session.GetUserId())
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)

}

// PatchConfig updates existing config
func (s *ConfigService) PatchConfig(ctx context.Context, in *proto.PatchConfigRequest) (*proto.Config, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, authmodel.Edit) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// RBAC check
	if scope.IsRbacUsed() {
		access, err := s.app.ConfigCheckAccess(ctx, session.GetDomainId(), int64(in.GetConfigId()), session.GetAclRoles(), accessMode)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, s.app.MakeScopeError(session, scope, accessMode)
		}
	}
	updatedConfigModel, err := ConvertPatchConfigMessageToModel(in, session.GetDomainId())
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.PatchUpdateConfig(ctx, updatedConfigModel, in.GetFields(), session.GetUserId())
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

// CreateConfig inserts new config
func (s *ConfigService) CreateConfig(ctx context.Context, in *proto.CreateConfigRequest) (*proto.Config, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	model, err := ConvertCreateConfigMessageToModel(in, session.GetDomainId())
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.InsertConfig(ctx, model, session.GetDomainId())
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

// DeleteConfig deletes config by id
func (s *ConfigService) DeleteConfig(ctx context.Context, in *proto.DeleteConfigRequest) (*proto.Empty, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, authmodel.Add) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	if scope.IsRbacUsed() {
		access, err := s.app.ConfigCheckAccess(ctx, session.GetDomainId(), int64(in.GetConfigId()), session.GetAclRoles(), accessMode)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, s.app.MakeScopeError(session, scope, accessMode)
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
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, authmodel.Add) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	var rbac *model.RbacOptions
	if scope.IsRbacUsed() {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: authmodel.Delete.Value(),
		}
	}
	appErr := s.app.DeleteConfigs(ctx, rbac, in.GetIds())
	if appErr != nil {
		return nil, appErr
	}
	return &proto.Empty{}, nil
}
