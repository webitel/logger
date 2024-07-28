package app

import (
	"context"
	"fmt"
	"github.com/webitel/logger/registry"
	"github.com/webitel/logger/registry/consul"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/webitel/logger/model"

	proto_grpc "buf.build/gen/go/webitel/logger/grpc/go/_gogrpc"
	"github.com/webitel/wlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	RequestContextName       = "grpc_ctx"
	AppServiceTtl            = time.Second * 30
	AppDeregesterCriticalTtl = time.Minute * 2
)

type AppServer struct {
	server   *grpc.Server
	listener net.Listener
	config   *model.ConsulConfig
	exitChan chan model.AppError
	registry registry.ServiceRegistrator
}

func BuildServer(app *App, config *model.ConsulConfig, exitChan chan model.AppError) (*AppServer, model.AppError) {
	// * Build grpc server
	server, appErr := buildGrpc(app)
	if appErr != nil {
		return nil, appErr
	}
	//  * Open tcp connection
	listener, err := net.Listen("tcp", config.PublicAddress)
	if err != nil {
		return nil, model.NewInternalError("api.grpc_server.serve_requests.listen.error", err.Error())
	}
	reg, appErr := consul.NewConsulRegistry(config)
	if appErr != nil {
		return nil, appErr
	}
	return &AppServer{
		server:   server,
		listener: listener,
		exitChan: exitChan,
		config:   config,
		registry: reg,
	}, nil
}

func (a *AppServer) Start() {
	appErr := a.registry.Register()
	if appErr != nil {
		a.exitChan <- appErr
		return
	}
	err := a.server.Serve(a.listener)
	if err != nil {
		a.exitChan <- model.NewInternalError("api.grpc_server.serve_requests.serve.error", err.Error())
		return
	}
}

func (a *AppServer) Stop() model.AppError {
	appErr := a.registry.Deregister()
	if appErr != nil {
		return appErr
	}
	a.stopGrpcServer()
	return nil
}

func (a *AppServer) stopGrpcServer() model.AppError {
	a.server.Stop()
	wlog.Info("grpc: server stopped")
	return nil
}

func buildGrpc(app *App) (*grpc.Server, model.AppError) {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
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
	proto_grpc.RegisterLoggerServiceServer(grpcServer, l)
	// * register config service
	proto_grpc.RegisterConfigServiceServer(grpcServer, c)

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

	h, err := handler(reqCtx, req)

	if err != nil {
		wlog.Error(fmt.Sprintf("[%s] method %s duration %s, error: %v", ip, info.FullMethod, time.Since(start), err.Error()))

		switch err.(type) {
		case model.AppError:
			e := err.(model.AppError)
			return h, status.Error(httpCodeToGrpc(e.GetStatusCode()), e.ToJson())
		default:
			return h, err
		}
	} else {
		wlog.Debug(fmt.Sprintf("[%s] method %s duration %s", ip, info.FullMethod, time.Since(start)))
	}

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
