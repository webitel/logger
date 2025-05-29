package grpc

import (
	"context"
	"errors"
	proto "github.com/webitel/logger/api/logger"
	"github.com/webitel/logger/internal/handler/grpc/utils"
	"github.com/webitel/logger/internal/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"strings"
)

type ConfigManager interface {
	SearchConfig(ctx context.Context, rbac *model.RbacOptions, searchOpt *model.SearchOptions) ([]*model.Config, error)
	UpdateConfig(ctx context.Context, in *model.Config, fields []string) (*model.Config, error)
	CreateConfig(ctx context.Context, in *model.Config) (*model.Config, error)
	DeleteConfig(ctx context.Context, ids []int) error

	GetConfigById(ctx context.Context, rbac *model.RbacOptions, id int) (*model.Config, error)
	GetConfigByObjectId(ctx context.Context, objectId int, domainId int) (*model.Config, error)
	CheckConfigStatus(ctx context.Context, objectName string, domainId int) (bool, error)
	GetSystemObjects(ctx context.Context, in *proto.ReadSystemObjectsRequest) ([]*model.SystemObject, error)
}

type ConfigService struct {
	app ConfigManager
	proto.UnimplementedConfigServiceServer
}

func NewConfigService(app ConfigManager) (*ConfigService, error) {
	if app == nil {
		return nil, errors.New("app is nil")
	}
	return &ConfigService{app: app}, nil
}

// ReadConfig selects config by id
func (s *ConfigService) ReadConfig(ctx context.Context, in *proto.ReadConfigRequest) (*proto.Config, error) {
	var rbac *model.RbacOptions
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	config, err := s.app.GetConfigById(ctx, rbac, int(in.GetConfigId()))
	if err != nil {
		return nil, err
	}
	message, err := s.Marshal(config)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}
	return message[0], nil
}

func (s *ConfigService) ReadSystemObjects(ctx context.Context, request *proto.ReadSystemObjectsRequest) (*proto.SystemObjects, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Bool("include_existing", request.GetIncludeExisting()))
	objects, err := s.app.GetSystemObjects(ctx, request)
	if err != nil {
		return nil, err
	}
	var resultObjects proto.SystemObjects
	for _, object := range objects {
		var obj proto.Lookup
		id := object.GetId()
		if id != nil {
			obj.Id = int32(*id)
		}
		name := object.GetName()
		if name != nil {
			obj.Name = *name
		}
		resultObjects.Items = append(resultObjects.Items, &obj)
	}
	return &resultObjects, nil
}

// ReadConfigByObjectId used for internal purpose
func (s *ConfigService) ReadConfigByObjectId(ctx context.Context, in *proto.ReadConfigByObjectIdRequest) (*proto.Config, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("domain.id", int(in.GetDomainId())))
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("object.id", int(in.GetObjectId())))
	config, err := s.app.GetConfigByObjectId(ctx, int(in.GetObjectId()), 0)
	if err != nil {
		return nil, err
	}
	message, err := s.Marshal(config)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}
	return message[0], nil
}

// ReadConfigByObjectId used for internal purpose with client, checks if config enabled
func (s *ConfigService) CheckConfigStatus(ctx context.Context, in *proto.CheckConfigStatusRequest) (*proto.ConfigStatus, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.Int("domain.id", int(in.GetDomainId())))
	GroupAttributesAndBindToSpan(ctx, "in", attribute.String("object.name", in.GetObjectName()))
	isEnabled, err := s.app.CheckConfigStatus(ctx, in.GetObjectName(), int(in.GetDomainId()))
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
	searchOpts := ExtractSearchOptions(in)
	configs, err := s.app.SearchConfig(ctx, rbac, searchOpts)
	if err != nil {
		return nil, err
	}
	messages, err := s.Marshal(configs...)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}

	res.Items, res.Next = utils.ResolvePaging(searchOpts.GetSize(), messages)
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
	config, err := s.app.UpdateConfig(ctx, mod, nil)
	if err != nil {
		return nil, err
	}
	message, err := s.Marshal(config)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}
	return message[0], nil

}

// PatchConfig updates existing config
func (s *ConfigService) PatchConfig(ctx context.Context, in *proto.PatchConfigRequest) (*proto.Config, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	updatedConfigModel, err := ConvertPatchConfigMessageToModel(in)
	if err != nil {
		return nil, err
	}
	config, err := s.app.UpdateConfig(ctx, updatedConfigModel, in.GetFields())
	if err != nil {
		return nil, err
	}
	message, err := s.Marshal(config)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}

	return message[0], nil
}

// CreateConfig inserts new config
func (s *ConfigService) CreateConfig(ctx context.Context, in *proto.CreateConfigRequest) (*proto.Config, error) {
	model, err := s.ConvertCreateConfigMessageToModel(in)
	if err != nil {
		return nil, err
	}
	config, err := s.app.CreateConfig(ctx, model)
	if err != nil {
		return nil, err
	}
	messages, err := s.Marshal(config)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}
	GroupOutgoingAttributesAndBindToSpan(ctx, attribute.Int("config.id", config.Id))

	return messages[0], nil
}

// DeleteConfig deletes config by id
func (s *ConfigService) DeleteConfig(ctx context.Context, in *proto.DeleteConfigRequest) (*proto.Empty, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	appErr := s.app.DeleteConfig(ctx, []int{int(in.GetConfigId())})
	if appErr != nil {
		return nil, appErr
	}
	return &proto.Empty{}, nil
}

func ConvertUpdateConfigMessageToModel(in *proto.UpdateConfigRequest) (*model.Config, error) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		Description: &in.Description,
		Storage:     utils.UnmarshalLookup(in.GetStorage(), &model.Storage{}),
	}
	return config, nil
}

func ConvertPatchConfigMessageToModel(in *proto.PatchConfigRequest) (*model.Config, error) {
	config := &model.Config{
		Id:          int(in.GetConfigId()),
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		Description: &in.Description,
		Storage:     utils.UnmarshalLookup(in.GetStorage(), &model.Storage{}),
	}
	return config, nil
}

func (s *ConfigService) ConvertCreateConfigMessageToModel(in *proto.CreateConfigRequest) (*model.Config, error) {
	if in.GetDaysToStore() <= 0 {
		return nil, errors.New("days to store should be greater than one")
	}
	config := &model.Config{
		Enabled:     in.GetEnabled(),
		DaysToStore: int(in.GetDaysToStore()),
		Period:      int(in.GetPeriod()),
		Description: &in.Description,
		Storage:     &model.Storage{},
		Object:      utils.UnmarshalLookup(in.GetObject(), &model.Object{}),
	}
	return config, nil
}

func (s *ConfigService) Marshal(models ...*model.Config) ([]*proto.Config, error) {
	var res []*proto.Config
	for _, in := range models {
		conf := &proto.Config{
			Id:          int32(in.Id),
			Enabled:     in.Enabled,
			DaysToStore: int32(in.DaysToStore),
			Period:      int32(in.Period),
			Storage:     utils.MarshalLookup(in.Storage),
			Object:      utils.MarshalLookup(in.Object),
		}
		if in.Description != nil {
			conf.Description = *in.Description
		}
		if in.LogsSize != nil {
			conf.LogsSize = *in.LogsSize
		}
		if in.LogsCount != nil {
			conf.LogsCount = int64(*in.LogsCount)
		}
		res = append(res, conf)
	}

	return res, nil
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
