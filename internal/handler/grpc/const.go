package grpc

import "github.com/webitel/logger/internal/model"

var (
	InternalError   = model.NewInternalError("application.process_api.error", "an internal error occurred while processing request")
	ConversionError = model.NewInternalError("application.conversion.error", "an internal error occurred while converting request/response")
	PermissionError = model.NewInternalError("application.permission.error", "access denied for the requested action")
	NotFoundError   = model.NewInternalError("application.not_found.error", "the requested resource was not found")
)
