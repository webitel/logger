package app

import (
	proto "buf.build/gen/go/webitel/logger/protocolbuffers/go"
	"context"
	authmodel "github.com/webitel/logger/auth/model"
	"github.com/webitel/logger/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	scope := model.ScopeLog
	accessMode := authmodel.Read
	// OBAC check
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	// RBAC check
	if session.UseRbacAccess(scope, accessMode) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: accessMode.Value(),
		}
	}
	resModel, err := s.app.GetConfigById(ctx, rbac, int(in.GetConfigId()))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

func (s *ConfigService) ReadSystemObjects(ctx context.Context, request *proto.ReadSystemObjectsRequest) (*proto.SystemObjects, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Bool("include_existing", request.GetIncludeExisting()))
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	scope := model.ScopeLog
	accessMode := authmodel.Read
	// OBAC check
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, authmodel.Read)
	}
	return s.app.GetSystemObjects(ctx, request, int(session.GetDomainId()))
}

// ReadConfigByObjectId used for internal purpose
func (s *ConfigService) ReadConfigByObjectId(ctx context.Context, in *proto.ReadConfigByObjectIdRequest) (*proto.Config, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("domain.id", int(in.GetDomainId())))
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("object.id", int(in.GetObjectId())))
	resModel, err := s.app.GetConfigByObjectId(ctx, int(in.GetDomainId()), int(in.GetObjectId()))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

// ReadConfigByObjectId used for internal purpose with client, checks if config enabled
func (s *ConfigService) CheckConfigStatus(ctx context.Context, in *proto.CheckConfigStatusRequest) (*proto.ConfigStatus, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("domain.id", int(in.GetDomainId())))
	GroupAttributesAndBindToSpan(ctx, "in", attribute.String("object.name", in.GetObjectName()))
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
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.String("q", in.GetQ()))
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	scope := model.ScopeLog
	accessMode := authmodel.Read
	// OBAC check
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	// RBAC check
	if session.UseRbacAccess(scope, accessMode) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: accessMode.Value(),
		}
	}
	searchOpts := ExtractSearchOptions(in)
	resModels, err := s.app.GetAllConfigs(ctx, rbac, searchOpts, session.GetDomainId())
	if err != nil {
		if IsErrNoRows(err) {
			return &res, nil
		} else {
			return nil, err
		}
	}
	res.Next, res.Items, err = CalculateListResultMetadata[*model.Config, *proto.Config](searchOpts, resModels, ConvertConfigModelToMessage)
	if err != nil {
		return nil, err
	}
	res.Page = in.GetPage()

	return &res, nil
}

// UpdateConfig updates existing config
func (s *ConfigService) UpdateConfig(ctx context.Context, in *proto.UpdateConfigRequest) (*proto.Config, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	scope := model.ScopeLog
	accessMode := authmodel.Edit
	// OBAC check
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	// RBAC check
	if session.UseRbacAccess(scope, accessMode) {
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
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// RBAC check
	if session.UseRbacAccess(scope, accessMode) {
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
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	model, err := ConvertCreateConfigMessageToModel(in, session.GetDomainId())
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.InsertConfig(ctx, model, session.GetUserId())
	if err != nil {
		return nil, err
	}
	GroupOutgoingAttributesAndBindToSpan(ctx, attribute.Int("config.id", resModel.Id))
	return ConvertConfigModelToMessage(resModel)
}

// DeleteConfig deletes config by id
func (s *ConfigService) DeleteConfig(ctx context.Context, in *proto.DeleteConfigRequest) (*proto.Empty, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	accessMode := authmodel.Delete
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	if session.UseRbacAccess(scope, accessMode) {
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
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	var rbac *model.RbacOptions
	if session.UseRbacAccess(scope, accessMode) {
		rbac = &model.RbacOptions{
			Groups: session.GetAclRoles(),
			Access: accessMode.Value(),
		}
	}
	appErr := s.app.DeleteConfigs(ctx, rbac, in.GetIds())
	if appErr != nil {
		return nil, appErr
	}
	return &proto.Empty{}, nil
}
