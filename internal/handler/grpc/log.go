package grpc

import (
	"context"
	deferr "errors"
	"github.com/webitel/logger/internal/handler/grpc/errors"
	"github.com/webitel/logger/internal/handler/grpc/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"time"

	"github.com/webitel/logger/internal/model"

	proto "github.com/webitel/logger/api/logger"
)

type LogManager interface {
	SearchLogs(ctx context.Context, searchOpt *model.SearchOptions, filters *model.LogFilters) ([]*model.Log, error)
	DeleteLogs(ctx context.Context, configId int, earlierThan time.Time) (int, error)
}

type LoggerService struct {
	app LogManager
	proto.UnimplementedLoggerServiceServer
}

func NewLoggerService(app LogManager) (*LoggerService, error) {
	if app == nil {
		return nil, deferr.New("app is nil")
	}
	return &LoggerService{app: app}, nil
}

func (s *LoggerService) SearchLogByRecordId(ctx context.Context, in *proto.SearchLogByRecordIdRequest) (*proto.Logs, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("record.id", int(in.GetRecordId())), attribute.String("object", in.GetObject().String()))
	// common log filters
	filters := extractDefaultFiltersFromLogSearch(in)
	// record id (required)
	if param := in.GetRecordId(); param != 0 {
		filters.User = []int64{int64(param)}
	} else {
		return nil, errors.NewBadRequestError("app.api_log.search_log_by_record.checks_args.error", "record id required")
	}
	// specific filters
	if param := in.GetUserId(); len(param) != 0 {
		filters.User = param
	}
	if param := in.GetObject(); param != 0 {
		filters.Object = []int64{int64(param)}
	}
	searchOpts := ExtractSearchOptions(in)
	models, err := s.app.SearchLogs(ctx, searchOpts, filters)
	if err != nil {
		return nil, err
	}
	messages, err := s.Marshal(models...)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}
	var res proto.Logs
	res.Items, res.Next = utils.ResolvePaging(searchOpts.GetSize(), messages)
	res.Page = in.GetPage()
	return &res, nil
}

func (s *LoggerService) SearchLogByUserId(ctx context.Context, in *proto.SearchLogByUserIdRequest) (*proto.Logs, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("user.id", int(in.GetUserId())))

	// common log filters
	filters := extractDefaultFiltersFromLogSearch(in)
	// user id (required)
	if param := in.GetUserId(); param != 0 {
		filters.User = []int64{int64(param)}
	} else {
		return nil, errors.NewBadRequestError("app.api_log.search_log_by_user.checks_args.error", "user id required")
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
	models, err := s.app.SearchLogs(ctx, searchOpts, filters)
	if err != nil {
		return nil, err
	}
	messages, err := s.Marshal(models...)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}
	var res proto.Logs
	res.Items, res.Next = utils.ResolvePaging(searchOpts.GetSize(), messages)
	res.Page = in.GetPage()
	return &res, nil
}

func (s *LoggerService) SearchLogByConfigId(ctx context.Context, in *proto.SearchLogByConfigIdRequest) (*proto.Logs, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(in.GetConfigId())))
	// common log filters
	filters := extractDefaultFiltersFromLogSearch(in)
	// config id (required)
	if param := in.GetConfigId(); param != 0 {
		filters.ConfigId = []int64{int64(param)}
	} else {
		return nil, errors.NewBadRequestError("app.api_log.search_log_by_user.checks_args.error", "user id required")
	}
	// specific filters
	if param := in.GetUserId(); len(param) != 0 {
		filters.User = param
	}

	searchOpts := ExtractSearchOptions(in)
	models, err := s.app.SearchLogs(ctx, searchOpts, filters)
	if err != nil {
		return nil, err
	}
	messages, err := s.Marshal(models...)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, ConversionError
	}
	var res proto.Logs
	res.Items, res.Next = utils.ResolvePaging(searchOpts.GetSize(), messages)
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

func (s *LoggerService) Marshal(models ...*model.Log) ([]*proto.Log, error) {

	var res []*proto.Log
	for _, m := range models {
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
		res = append(res, log)
	}

	return res, nil
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
