package grpc

import (
	"context"
	"github.com/webitel/logger/internal/handler/grpc/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"time"

	"github.com/webitel/logger/internal/model"

	proto "github.com/webitel/logger/api/logger"
)

type LogManager interface {
	SearchLogs(ctx context.Context, searchOpt *model.SearchOptions, filters *model.LogFilters) ([]*model.Log, model.AppError)
	DeleteLogs(ctx context.Context, configId int, earlierThan time.Time) (int, model.AppError)
}

type LoggerService struct {
	app LogManager
	proto.UnimplementedLoggerServiceServer
}

func NewLoggerService(app LogManager) (*LoggerService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_logger_service.args_check.app_nil", "app is nil")
	}
	return &LoggerService{app: app}, nil
}

func (s *LoggerService) SearchLogByRecordId(ctx context.Context, in *proto.SearchLogByRecordIdRequest) (*proto.Logs, error) {
	var (
		res proto.Logs
	)
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("record.id", int(in.GetRecordId())), attribute.String("object", in.GetObject().String()))
	// common log filters
	filters := extractDefaultFiltersFromLogSearch(in)
	// record id (required)
	if param := in.GetRecordId(); param != 0 {
		filters.User = []int64{int64(param)}
	} else {
		return nil, model.NewBadRequestError("app.api_log.search_log_by_record.checks_args.error", "record id required")
	}
	// specific filters
	if param := in.GetUserId(); len(param) != 0 {
		filters.User = param
	}
	if param := in.GetObject(); param != 0 {
		filters.Object = []int64{int64(param)}
	}
	searchOpts := ExtractSearchOptions(in)
	resModels, err := s.app.SearchLogs(ctx, searchOpts, filters)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items, err = utils.CalculateListResultMetadata[*model.Log, *proto.Log](searchOpts, resModels, convertLogModelToMessage)
	if err != nil {
		return nil, err
	}
	res.Page = in.GetPage()
	return &res, nil
}

func (s *LoggerService) SearchLogByUserId(ctx context.Context, in *proto.SearchLogByUserIdRequest) (*proto.Logs, error) {
	var (
		res proto.Logs
	)
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("user.id", int(in.GetUserId())))

	// common log filters
	filters := extractDefaultFiltersFromLogSearch(in)
	// user id (required)
	if param := in.GetUserId(); param != 0 {
		filters.User = []int64{int64(param)}
	} else {
		return nil, model.NewBadRequestError("app.api_log.search_log_by_user.checks_args.error", "user id required")
	}
	// specific filters
	if param := in.GetObjectId(); len(param) != 0 {
		filters.Object = param
	}
	if param := in.GetUserId(); param != 0 {
		filters.User = []int64{int64(param)}
	}

	// perform
	searchOpts := ExtractSearchOptions(in)
	resModels, err := s.app.SearchLogs(ctx, searchOpts, filters)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items, err = utils.CalculateListResultMetadata[*model.Log, *proto.Log](searchOpts, resModels, convertLogModelToMessage)
	if err != nil {
		return nil, err
	}
	res.Page = in.GetPage()
	return &res, nil
}

func (s *LoggerService) SearchLogByConfigId(ctx context.Context, in *proto.SearchLogByConfigIdRequest) (*proto.Logs, error) {
	var (
		res proto.Logs
	)
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	// common log filters
	filters := extractDefaultFiltersFromLogSearch(in)
	// config id (required)
	if param := in.GetConfigId(); param != 0 {
		filters.ConfigId = []int64{int64(param)}
	} else {
		return nil, model.NewBadRequestError("app.api_log.search_log_by_user.checks_args.error", "user id required")
	}
	// specific filters
	if param := in.GetUserId(); len(param) != 0 {
		filters.User = param
	}

	searchOpts := ExtractSearchOptions(in)
	resModels, err := s.app.SearchLogs(ctx, searchOpts, filters)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items, err = utils.CalculateListResultMetadata[*model.Log, *proto.Log](searchOpts, resModels, convertLogModelToMessage)
	if err != nil {
		return nil, err
	}
	res.Page = in.GetPage()
	return &res, nil
}

func (s *LoggerService) DeleteConfigLogs(ctx context.Context, request *proto.DeleteConfigLogsRequest) (*proto.DeleteConfigLogsResponse, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(request.GetConfigId())))

	var olderThan time.Time
	if request.GetOlderThan() != 0 {
		olderThan = time.UnixMilli(request.OlderThan)
	}
	processed, err := s.app.DeleteLogs(ctx, int(request.ConfigId), olderThan)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteConfigLogsResponse{Processed: int64(processed)}, nil
}

// region utility

// Fills the
// DateFrom, DateTo, UserIp, Actions
func extractDefaultFiltersFromLogSearch(in LogSearcher) *model.LogFilters {
	filters := &model.LogFilters{}
	if param := in.GetDateFrom(); param != 0 {
		t := time.Unix(param/1000, 0).UTC()
		filters.DateFrom = &t
	}
	if param := in.GetDateTo(); param != 0 {
		t := time.Unix(param/1000, 0).UTC()
		filters.DateTo = &t
	}
	if param := in.GetUserIp(); param != "" {
		filters.UserIp = append(filters.UserIp, param)
	}
	if param := in.GetAction(); len(param) != 0 {
		for _, action := range param {
			filters.Action = append(filters.Action, action.String())
		}
	}
	return filters
}

func (s *LoggerService) Marshal(m *model.Log) (*proto.Log, model.AppError) {
	log := &proto.Log{
		Id:     int32(m.Id),
		Action: m.Action,

		UserIp:   m.UserIp,
		NewState: string(m.NewState),
		ConfigId: int32(m.ConfigId),
	}
	if !m.User.IsZero() {
		log.User = &proto.Lookup{
			Id:   int32(m.User.Id.Int()),
			Name: m.User.Name.String(),
		}
	}
	if !m.Object.IsZero() {
		log.Object = &proto.Lookup{
			Id:   int32(m.Object.Id.Int()),
			Name: m.Object.Name.String(),
		}
	}
	if !m.Date.IsZero() {
		log.Date = m.Date.ToMilliseconds()
	}
	if s := m.Record.Id.Int32(); s != 0 {
		log.Record = &proto.Lookup{
			Id:   s,
			Name: m.Record.Name.String(),
		}
	}
	return log, nil
}

func (s *LoggerService) Unmarshal(m *model.Log) (*proto.Log, model.AppError) {
	log := &proto.Log{
		Id:     int32(m.Id),
		Action: m.Action,

		UserIp:   m.UserIp,
		NewState: string(m.NewState),
		ConfigId: int32(m.ConfigId),
	}
	if !m.User.IsZero() {
		log.User = &proto.Lookup{
			Id:   int32(m.User.Id.Int()),
			Name: m.User.Name.String(),
		}
	}
	if !m.Object.IsZero() {
		log.Object = &proto.Lookup{
			Id:   int32(m.Object.Id.Int()),
			Name: m.Object.Name.String(),
		}
	}
	if !m.Date.IsZero() {
		log.Date = m.Date.ToMilliseconds()
	}
	if s := m.Record.Id.Int32(); s != 0 {
		log.Record = &proto.Lookup{
			Id:   s,
			Name: m.Record.Name.String(),
		}
	}
	return log, nil
}

func convertLogModelToMessageBulk(m []*model.Log) ([]*proto.Log, model.AppError) {
	var rows []*proto.Log
	for _, v := range m {
		protoLog, err := convertLogModelToMessage(v)
		if err != nil {
			return nil, err
		}
		rows = append(rows, protoLog)
	}
	return rows, nil
}

func GroupAttributesAndBindToSpan(ctx context.Context, rootName string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	group := model.NewAttributeGroup(rootName, attrs...)
	span.SetAttributes(group.Unparse()...)
}

func GroupIncomingAttributesAndBindToSpan(ctx context.Context, attrs ...attribute.KeyValue) {
	GroupAttributesAndBindToSpan(ctx, "in", attrs...)
}

func GroupOutgoingAttributesAndBindToSpan(ctx context.Context, attrs ...attribute.KeyValue) {
	GroupAttributesAndBindToSpan(ctx, "out", attrs...)
}

type LogSearcher interface {
	GetDateFrom() int64
	GetDateTo() int64
	GetUserIp() string
	GetAction() []proto.Action
}

// endregion
