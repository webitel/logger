package grpc

import (
	"context"
	"fmt"
	proto "github.com/webitel/logger/api/logger"
	"github.com/webitel/logger/internal/model"
	otelgrpc "github.com/webitel/webitel-go-kit/tracing/grpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	span := trace.SpanFromContext(ctx)
	var reqCtx context.Context
	var ip string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqCtx = context.WithValue(ctx, RequestContextName, md)
		ip = getClientIp(md)
	} else {
		ip = "<not found>"
		reqCtx = context.WithValue(ctx, RequestContextName, nil)
	}
	span.SetAttributes(attribute.String("caller_user.ip", ip))
	h, err := handler(reqCtx, req)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("[%s] method %s duration %s, error: %v", ip, info.FullMethod, time.Since(start), err.Error()))
		span.RecordError(err)
		switch err.(type) {
		case model.AppError:
			e := err.(model.AppError)
			return h, status.Error(httpCodeToGrpc(e.GetStatusCode()), e.ToJson())
		default:
			return h, err
		}
	}
	slog.DebugContext(ctx, fmt.Sprintf("[%s] method %s duration %s", ip, info.FullMethod, time.Since(start)))
	return h, err
}

func httpCodeToGrpc(c int) codes.Code {
	switch c {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusAccepted:
		return codes.ResourceExhausted
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	default:
		return codes.Internal
	}
}

func getClientIp(info metadata.MD) string {
	ip := strings.Join(info.Get("x-real-ip"), ",")
	if ip == "" {
		ip = strings.Join(info.Get("x-forwarded-for"), ",")
	}

	return ip
}
