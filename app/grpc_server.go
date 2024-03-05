package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/webitel/engine/discovery"
	errors "github.com/webitel/engine/model"
	"github.com/webitel/logger/model"

	proto "github.com/webitel/logger/api/native"
	"github.com/webitel/wlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	RequestContextName          = "grpc_ctx"
	APP_SERVICE_TTL             = time.Second * 30
	APP_DEREGESTER_CRITICAL_TTL = time.Minute * 2
)

//func ServeRequests(app *App, config *model.ConsulConfig, errChan chan errors.AppError) {
//	// * Build grpc server
//	server, appErr := buildGrpc(app)
//	if appErr != nil {
//		errChan <- appErr
//		return
//	}
//	//  * Open tcp connection
//	listener, err := net.Listen("tcp", config.PublicAddress)
//	if err != nil {
//		errChan <- errors.NewInternalError("api.grpc_server.serve_requests.listen.error", err.Error())
//		return
//	}
//	appErr = connectConsul(config)
//	if appErr != nil {
//		errChan <- appErr
//		return
//	}
//	err = server.Serve(listener)
//	if err != nil {
//		errChan <- errors.NewInternalError("api.grpc_server.serve_requests.serve.error", err.Error())
//		return
//	}
//}

type AppServer struct {
	server   *grpc.Server
	listener net.Listener
	config   *model.ConsulConfig
	exitChan chan errors.AppError
}

func BuildServer(app *App, config *model.ConsulConfig, exitChan chan errors.AppError) (*AppServer, errors.AppError) {
	// * Build grpc server
	server, appErr := buildGrpc(app)
	if appErr != nil {
		return nil, appErr
	}
	//  * Open tcp connection
	listener, err := net.Listen("tcp", config.PublicAddress)
	if err != nil {
		return nil, errors.NewInternalError("api.grpc_server.serve_requests.listen.error", err.Error())
	}

	return &AppServer{
		server:   server,
		listener: listener,
		exitChan: exitChan,
		config:   config,
	}, nil
}

func (a *AppServer) Start() {
	appErr := connectConsul(a.config)
	if appErr != nil {
		a.exitChan <- appErr
		return
	}
	err := a.server.Serve(a.listener)
	if err != nil {
		a.exitChan <- errors.NewInternalError("api.grpc_server.serve_requests.serve.error", err.Error())
		return
	}
}

func buildGrpc(app *App) (*grpc.Server, errors.AppError) {

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

	h, err := handler(reqCtx, req)

	if err != nil {
		wlog.Error(fmt.Sprintf("[%s] method %s duration %s, error: %v", ip, info.FullMethod, time.Since(start), err.Error()))

		switch err.(type) {
		case errors.AppError:
			e := err.(errors.AppError)
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

func connectConsul(config *model.ConsulConfig) errors.AppError {
	if config.Id == "" {
		errors.NewBadRequestError("api.grpc_server.build_consul.service_id.error", "service id is empty! (set it by '-id' flag)")
	}
	consul, err := discovery.NewConsul(config.Id, config.Address, func() (bool, error) {
		return true, nil
	})
	ip, port, err := net.SplitHostPort(config.PublicAddress)
	if err != nil {
		return errors.NewBadRequestError("api.grpc_server.build_consul.parse_address.error", "unable to parse address")
	}
	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return errors.NewBadRequestError("api.grpc_server.build_consul.parse_address.error", "unable to parse grpc port")
	}
	err = consul.RegisterService(model.SERVICE_NAME, ip, parsedPort, APP_SERVICE_TTL, APP_DEREGESTER_CRITICAL_TTL)
	if err != nil {
		return errors.NewInternalError("api.grpc_server.build_consul.register_in_consul.error", err.Error())
	}
	return nil
}
