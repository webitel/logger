package grpc

import (
	"context"
	"fmt"
	proto "github.com/webitel/logger/api/logger"
	autherrors "github.com/webitel/logger/internal/auth/errors"
	"github.com/webitel/logger/internal/handler/grpc/errors"
	otelgrpc "github.com/webitel/webitel-go-kit/tracing/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	RequestContextName = "grpc_ctx"
)

type Handler interface {
	LogManager
	ConfigManager
}

func Build(app Handler) (*grpc.Server, error) {
	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithMessageEvents(otelgrpc.SentEvents, otelgrpc.ReceivedEvents))), grpc.UnaryInterceptor(unaryInterceptor))
	// * Creating services
	l, appErr := NewLoggerService(app)
	if appErr != nil {
		return nil, appErr
	}
	c, appErr := NewConfigService(app)
	if appErr != nil {
		return nil, appErr
	}

	// * register logger service
	proto.RegisterLoggerServiceServer(grpcServer, l)
	// * register config service
	proto.RegisterConfigServiceServer(grpcServer, c)

	return grpcServer, nil
}

func unaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	var reqCtx context.Context
	var ip string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqCtx = context.WithValue(ctx, RequestContextName, md)
		ip = getClientIp(md)
	} else {
		ip = "<not found>"
		reqCtx = context.WithValue(ctx, RequestContextName, nil)
	}
	res, err := handler(reqCtx, req)
	if err != nil {
		switch typedErr := err.(type) {
		case errors.AppError:
			return nil, status.Error(httpCodeToGrpc(typedErr.GetStatusCode()), typedErr.ToJson())
		case *autherrors.AuthorizationError:
			return nil, status.Error(httpCodeToGrpc(typedErr.GetStatusCode()), typedErr.ToJson())
		default:
			slog.ErrorContext(ctx, typedErr.Error())
			return nil, errors.ErrInternal
		}
	}
	slog.DebugContext(ctx, fmt.Sprintf("[%s] method %s duration %s", ip, info.FullMethod, time.Since(start)))
	return res, err
}

func httpCodeToGrpc(c int) codes.Code {
	switch c {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusRequestTimeout:
		return codes.DeadlineExceeded
	case http.StatusConflict:
		return codes.Aborted
	case http.StatusGone:
		return codes.NotFound
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Unknown
	}
}

func getClientIp(info metadata.MD) string {
	ip := strings.Join(info.Get("x-real-ip"), ",")
	if ip == "" {
		ip = strings.Join(info.Get("x-forwarded-for"), ",")
	}

	return ip
}
