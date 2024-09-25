package app

import (
	"context"
	authmodel "github.com/webitel/logger/auth/model"
	"go.opentelemetry.io/otel/attribute"
	"time"

	"github.com/webitel/logger/model"

	proto "buf.build/gen/go/webitel/logger/protocolbuffers/go"
)

type LoggerService struct {
	app *App
}

func NewLoggerService(app *App) (*LoggerService, model.AppError) {
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
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Read
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
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
	res.Next, res.Items, err = CalculateListResultMetadata[*model.Log, *proto.Log](searchOpts, resModels, convertLogModelToMessage)
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
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Read
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

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
	res.Next, res.Items, err = CalculateListResultMetadata[*model.Log, *proto.Log](searchOpts, resModels, convertLogModelToMessage)
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
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := model.ScopeLog
	if !session.HasObacAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

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
	res.Next, res.Items, err = CalculateListResultMetadata[*model.Log, *proto.Log](searchOpts, resModels, convertLogModelToMessage)
	if err != nil {
		return nil, err
	}
	res.Page = in.GetPage()
	return &res, nil
}

func (s *LoggerService) DeleteConfigLogs(ctx context.Context, request *proto.DeleteConfigLogsRequest) (*proto.DeleteConfigLogsResponse, error) {
	GroupIncomingAttributesAndBindToSpan(ctx, attribute.Int("config.id", int(request.GetConfigId())))
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	accessMode := authmodel.Edit
	secondaryMode := authmodel.Delete
	scope := model.ScopeLog
	// Edit or Delete permissions allow this operation
	if !session.HasObacAccess(scope, accessMode) && !session.HasObacAccess(scope, secondaryMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}
	// RBAC check
	if session.UseRbacAccess(scope, accessMode) && session.UseRbacAccess(scope, secondaryMode) {
		rbacAccess, err := s.app.storage.Config().CheckAccess(ctx, session.GetDomainId(), request.GetConfigId(), session.GetAclRoles(), uint8(accessMode))
		if err != nil {
			return nil, err
		}
		secondaryRbacAccess, err := s.app.storage.Config().CheckAccess(ctx, session.GetDomainId(), request.GetConfigId(), session.GetAclRoles(), uint8(secondaryMode))
		if err != nil {
			return nil, err
		}
		if !rbacAccess && !secondaryRbacAccess {
			return nil, s.app.MakeScopeError(session, scope, accessMode)
		}
	}

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
