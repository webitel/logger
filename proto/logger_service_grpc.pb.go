// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: proto/logger_service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	LoggerService_GetLogsByUserId_FullMethodName   = "/webitel_logger.LoggerService/GetLogsByUserId"
	LoggerService_GetLogsByObjectId_FullMethodName = "/webitel_logger.LoggerService/GetLogsByObjectId"
	LoggerService_GetLogsByConfigId_FullMethodName = "/webitel_logger.LoggerService/GetLogsByConfigId"
)

// LoggerServiceClient is the client API for LoggerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoggerServiceClient interface {
	GetLogsByUserId(ctx context.Context, in *GetLogsByUserIdRequest, opts ...grpc.CallOption) (*Logs, error)
	GetLogsByObjectId(ctx context.Context, in *GetLogsByObjectIdRequest, opts ...grpc.CallOption) (*Logs, error)
	GetLogsByConfigId(ctx context.Context, in *GetLogsByConfigIdRequest, opts ...grpc.CallOption) (*Logs, error)
}

type loggerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLoggerServiceClient(cc grpc.ClientConnInterface) LoggerServiceClient {
	return &loggerServiceClient{cc}
}

func (c *loggerServiceClient) GetLogsByUserId(ctx context.Context, in *GetLogsByUserIdRequest, opts ...grpc.CallOption) (*Logs, error) {
	out := new(Logs)
	err := c.cc.Invoke(ctx, LoggerService_GetLogsByUserId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggerServiceClient) GetLogsByObjectId(ctx context.Context, in *GetLogsByObjectIdRequest, opts ...grpc.CallOption) (*Logs, error) {
	out := new(Logs)
	err := c.cc.Invoke(ctx, LoggerService_GetLogsByObjectId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *loggerServiceClient) GetLogsByConfigId(ctx context.Context, in *GetLogsByConfigIdRequest, opts ...grpc.CallOption) (*Logs, error) {
	out := new(Logs)
	err := c.cc.Invoke(ctx, LoggerService_GetLogsByConfigId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoggerServiceServer is the server API for LoggerService service.
// All implementations must embed UnimplementedLoggerServiceServer
// for forward compatibility
type LoggerServiceServer interface {
	GetLogsByUserId(context.Context, *GetLogsByUserIdRequest) (*Logs, error)
	GetLogsByObjectId(context.Context, *GetLogsByObjectIdRequest) (*Logs, error)
	GetLogsByConfigId(context.Context, *GetLogsByConfigIdRequest) (*Logs, error)
	mustEmbedUnimplementedLoggerServiceServer()
}

// UnimplementedLoggerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLoggerServiceServer struct {
}

func (UnimplementedLoggerServiceServer) GetLogsByUserId(context.Context, *GetLogsByUserIdRequest) (*Logs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLogsByUserId not implemented")
}
func (UnimplementedLoggerServiceServer) GetLogsByObjectId(context.Context, *GetLogsByObjectIdRequest) (*Logs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLogsByObjectId not implemented")
}
func (UnimplementedLoggerServiceServer) GetLogsByConfigId(context.Context, *GetLogsByConfigIdRequest) (*Logs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLogsByConfigId not implemented")
}
func (UnimplementedLoggerServiceServer) mustEmbedUnimplementedLoggerServiceServer() {}

// UnsafeLoggerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoggerServiceServer will
// result in compilation errors.
type UnsafeLoggerServiceServer interface {
	mustEmbedUnimplementedLoggerServiceServer()
}

func RegisterLoggerServiceServer(s grpc.ServiceRegistrar, srv LoggerServiceServer) {
	s.RegisterService(&LoggerService_ServiceDesc, srv)
}

func _LoggerService_GetLogsByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogsByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggerServiceServer).GetLogsByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoggerService_GetLogsByUserId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggerServiceServer).GetLogsByUserId(ctx, req.(*GetLogsByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggerService_GetLogsByObjectId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogsByObjectIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggerServiceServer).GetLogsByObjectId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoggerService_GetLogsByObjectId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggerServiceServer).GetLogsByObjectId(ctx, req.(*GetLogsByObjectIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LoggerService_GetLogsByConfigId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogsByConfigIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoggerServiceServer).GetLogsByConfigId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoggerService_GetLogsByConfigId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoggerServiceServer).GetLogsByConfigId(ctx, req.(*GetLogsByConfigIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LoggerService_ServiceDesc is the grpc.ServiceDesc for LoggerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoggerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel_logger.LoggerService",
	HandlerType: (*LoggerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLogsByUserId",
			Handler:    _LoggerService_GetLogsByUserId_Handler,
		},
		{
			MethodName: "GetLogsByObjectId",
			Handler:    _LoggerService_GetLogsByObjectId_Handler,
		},
		{
			MethodName: "GetLogsByConfigId",
			Handler:    _LoggerService_GetLogsByConfigId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/logger_service.proto",
}
