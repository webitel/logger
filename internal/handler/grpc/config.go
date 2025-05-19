package grpc

import (
	"context"
	proto "github.com/webitel/logger/api/logger"
	"github.com/webitel/logger/internal/handler/grpc/utils"
	"github.com/webitel/logger/internal/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strings"
	"time"
)

type ConfigManager interface {
	SearchConfig(ctx context.Context, rbac *model.RbacOptions, searchOpt *model.SearchOptions) ([]*model.Config, model.AppError)
	UpdateConfig(ctx context.Context, in *model.Config) (*model.Config, model.AppError)
	InsertConfig(ctx context.Context, in *model.Config) (*model.Config, model.AppError)
	DeleteConfig(ctx context.Context, ids []int32) model.AppError

	GetConfigById(ctx context.Context, rbac *model.RbacOptions, id int) (*model.Config, model.AppError)
	GetConfigByObjectId(ctx context.Context, objectId int, domainId int) (*model.Config, model.AppError)
	CheckConfigStatus(ctx context.Context, objectName string, domainId int) (bool, model.AppError)
	GetSystemObjects(ctx context.Context, in *proto.ReadSystemObjectsRequest) (*proto.SystemObjects, error)
}

type ConfigService struct {
	app ConfigManager
	proto.UnimplementedConfigServiceServer
}

func NewConfigService(app ConfigManager) (*ConfigService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_config_service.args_check.app_nil", "app is nil")
	}
	return &ConfigService{app: app}, nil
}

// ReadConfig selects config by id
func (s *ConfigService) ReadConfig(ctx context.Context, in *proto.ReadConfigRequest) (*proto.Config, error) {
	var rbac *model.RbacOptions
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	resModel, err := s.app.GetConfigById(ctx, rbac, int(in.GetConfigId()))
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

func (s *ConfigService) ReadSystemObjects(ctx context.Context, request *proto.ReadSystemObjectsRequest) (*proto.SystemObjects, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Bool("include_existing", request.GetIncludeExisting()))
	return s.app.GetSystemObjects(ctx, request)
}

// ReadConfigByObjectId used for internal purpose
func (s *ConfigService) ReadConfigByObjectId(ctx context.Context, in *proto.ReadConfigByObjectIdRequest) (*proto.Config, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("domain.id", int(in.GetDomainId())))
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("object.id", int(in.GetObjectId())))
	resModel, err := s.app.GetConfigByObjectId(ctx, int(in.GetObjectId()), 0)
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
	isEnabled, err := s.app.CheckConfigStatus(ctx, in.GetObjectName(), 0)
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
	searchOpts := ExtractSearchOptions(in)
	resModels, err := s.app.SearchConfig(ctx, rbac, searchOpts)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items, err = utils.CalculateListResultMetadata[*model.Config, *proto.Config](searchOpts, resModels, ConvertConfigModelToMessage)
	if err != nil {
		return nil, err
	}
	res.Page = in.GetPage()

	return &res, nil
}

// UpdateConfig updates existing config
func (s *ConfigService) UpdateConfig(ctx context.Context, in *proto.UpdateConfigRequest) (*proto.Config, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))

	mod, err := ConvertUpdateConfigMessageToModel(in)
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.UpdateConfig(ctx, mod)
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)

}

// PatchConfig updates existing config
func (s *ConfigService) PatchConfig(ctx context.Context, in *proto.PatchConfigRequest) (*proto.Config, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	updatedConfigModel, err := ConvertPatchConfigMessageToModel(in)
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.UpdateConfig(ctx, updatedConfigModel)
	if err != nil {
		return nil, err
	}
	return ConvertConfigModelToMessage(resModel)
}

// CreateConfig inserts new config
func (s *ConfigService) CreateConfig(ctx context.Context, in *proto.CreateConfigRequest) (*proto.Config, error) {
	model, err := ConvertCreateConfigMessageToModel(in)
	if err != nil {
		return nil, err
	}
	resModel, err := s.app.InsertConfig(ctx, model)
	if err != nil {
		return nil, err
	}
	GroupOutgoingAttributesAndBindToSpan(ctx, attribute.Int("config.id", resModel.Id))
	return ConvertConfigModelToMessage(resModel)
}

// DeleteConfig deletes config by id
func (s *ConfigService) DeleteConfig(ctx context.Context, in *proto.DeleteConfigRequest) (*proto.Empty, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	appErr := s.app.DeleteConfig(ctx, []int32{in.GetConfigId()})
	if appErr != nil {
		return nil, appErr
	}
	return &proto.Empty{}, nil
}

func ConvertUpdateConfigMessageToModel(in *proto.UpdateConfigRequest) (*model.Config, model.AppError) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		Description: *model.NewNullString(in.GetDescription()),
	}

	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, model.NewInternalError("app.config.convert_update_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func ConvertPatchConfigMessageToModel(in *proto.PatchConfigRequest) (*model.Config, model.AppError) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		Description: *model.NewNullString(in.GetDescription()),
	}
	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, model.NewInternalError("app.config.convert_patch_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func ConvertCreateConfigMessageToModel(in *proto.CreateConfigRequest) (*model.Config, model.AppError) {

	if in.GetDaysToStore() <= 0 {
		return nil, model.NewBadRequestError("app.config.convert_create_config_message.bad_args", "days to store should be greater than one")
	}
	config := &model.Config{

		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		Description: *model.NewNullString(in.GetDescription()),
	}
	objectId, err := model.NewNullInt(in.GetObject().GetId())
	if err != nil {
		return nil, model.NewInternalError("app.config.convert_create_config_message.convert_object_id.fail", err.Error())
	}
	config.Object.Id = objectId

	if v := in.GetStorage().GetId(); v != 0 {
		storageId, err := model.NewNullInt(in.GetStorage().GetId())
		if err != nil {
			return nil, model.NewInternalError("app.config.convert_create_config_message.convert_storage_id.fail", err.Error())
		}
		config.Storage.Id = storageId
	}

	return config, nil
}

func ConvertConfigModelToMessage(in *model.Config) (*proto.Config, model.AppError) {

	conf := &proto.Config{
		Id:          int32(in.Id),
		Enabled:     in.Enabled,
		DaysToStore: int32(in.DaysToStore),
		Period:      int32(in.Period),
		Description: in.Description.String(),
		LogsSize:    in.LogsSize.String(),
		LogsCount:   in.LogsCount.Int64(),
	}
	if !in.Object.IsZero() {
		conf.Object = &proto.Lookup{
			Id:   int32(in.Object.Id.Int()),
			Name: in.Object.Name.String(),
		}
	}
	if !in.Storage.IsZero() {
		conf.Storage = &proto.Lookup{
			Id:   int32(in.Storage.Id.Int()),
			Name: in.Storage.Name.String(),
		}
	}
	return conf, nil
}

func calculateNextPeriodFromDate(period int, from time.Time) *model.NullTime {
	return model.NewNullTime(from.Add(time.Hour * 24 * time.Duration(period)))

}

func ExtractSearchOptions(t model.Searcher) *model.SearchOptions {
	var res model.SearchOptions
	if t.GetSort() != "" {
		res.Sort = model.ConvertSort(t.GetSort())
	}
	if t.GetSize() <= 0 || t.GetSize() > model.MaxPageSize {
		res.Size = model.DefaultPageSize
	} else {
		res.Size = int(t.GetSize())
	}
	if t.GetPage() <= 0 {
		res.Page = model.DefaultPage
	} else {
		res.Page = int(t.GetPage())
	}
	if t.GetQ() != "" {
		//	if input := strings.Replace(t.GetQ(), "*", "%", -1); input == "" {
		res.Search = strings.Replace(t.GetQ(), "*", "%", -1)
		//	}

	}
	if s := t.GetFields(); len(s) != 0 {
		res.Fields = s
	}
	return &res
}
