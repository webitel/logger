package grpc

import "github.com/webitel/logger/internal/handler/grpc/errors"

var (
	ConversionError = errors.NewInternalError("application.conversion.error", "an internal error occurred while converting request/response")
)
