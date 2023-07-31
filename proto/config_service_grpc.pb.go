// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: config_service.proto

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

// ConfigServiceClient is the client API for ConfigService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigServiceClient interface {
	UpdateConfig(ctx context.Context, in *UpdateConfigRequest, opts ...grpc.CallOption) (*Config, error)
	InsertConfig(ctx context.Context, in *InsertConfigRequest, opts ...grpc.CallOption) (*Config, error)
	GetConfigByObjectId(ctx context.Context, in *GetConfigByObjectIdRequest, opts ...grpc.CallOption) (*Config, error)
	GetConfigById(ctx context.Context, in *GetConfigByIdRequest, opts ...grpc.CallOption) (*Config, error)
	GetAllConfigs(ctx context.Context, in *GetAllConfigsRequest, opts ...grpc.CallOption) (*Configs, error)
}

type configServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigServiceClient(cc grpc.ClientConnInterface) ConfigServiceClient {
	return &configServiceClient{cc}
}

func (c *configServiceClient) UpdateConfig(ctx context.Context, in *UpdateConfigRequest, opts ...grpc.CallOption) (*Config, error) {
	out := new(Config)
	err := c.cc.Invoke(ctx, "/webitel_logger.ConfigService/UpdateConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceClient) InsertConfig(ctx context.Context, in *InsertConfigRequest, opts ...grpc.CallOption) (*Config, error) {
	out := new(Config)
	err := c.cc.Invoke(ctx, "/webitel_logger.ConfigService/InsertConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceClient) GetConfigByObjectId(ctx context.Context, in *GetConfigByObjectIdRequest, opts ...grpc.CallOption) (*Config, error) {
	out := new(Config)
	err := c.cc.Invoke(ctx, "/webitel_logger.ConfigService/GetConfigByObjectId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceClient) GetConfigById(ctx context.Context, in *GetConfigByIdRequest, opts ...grpc.CallOption) (*Config, error) {
	out := new(Config)
	err := c.cc.Invoke(ctx, "/webitel_logger.ConfigService/GetConfigById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configServiceClient) GetAllConfigs(ctx context.Context, in *GetAllConfigsRequest, opts ...grpc.CallOption) (*Configs, error) {
	out := new(Configs)
	err := c.cc.Invoke(ctx, "/webitel_logger.ConfigService/GetAllConfigs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConfigServiceServer is the server API for ConfigService service.
// All implementations must embed UnimplementedConfigServiceServer
// for forward compatibility
type ConfigServiceServer interface {
	UpdateConfig(context.Context, *UpdateConfigRequest) (*Config, error)
	InsertConfig(context.Context, *InsertConfigRequest) (*Config, error)
	GetConfigByObjectId(context.Context, *GetConfigByObjectIdRequest) (*Config, error)
	GetConfigById(context.Context, *GetConfigByIdRequest) (*Config, error)
	GetAllConfigs(context.Context, *GetAllConfigsRequest) (*Configs, error)
	mustEmbedUnimplementedConfigServiceServer()
}

// UnimplementedConfigServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConfigServiceServer struct {
}

func (UnimplementedConfigServiceServer) UpdateConfig(context.Context, *UpdateConfigRequest) (*Config, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateConfig not implemented")
}
func (UnimplementedConfigServiceServer) InsertConfig(context.Context, *InsertConfigRequest) (*Config, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertConfig not implemented")
}
func (UnimplementedConfigServiceServer) GetConfigByObjectId(context.Context, *GetConfigByObjectIdRequest) (*Config, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfigByObjectId not implemented")
}
func (UnimplementedConfigServiceServer) GetConfigById(context.Context, *GetConfigByIdRequest) (*Config, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfigById not implemented")
}
func (UnimplementedConfigServiceServer) GetAllConfigs(context.Context, *GetAllConfigsRequest) (*Configs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllConfigs not implemented")
}
func (UnimplementedConfigServiceServer) mustEmbedUnimplementedConfigServiceServer() {}

// UnsafeConfigServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigServiceServer will
// result in compilation errors.
type UnsafeConfigServiceServer interface {
	mustEmbedUnimplementedConfigServiceServer()
}

func RegisterConfigServiceServer(s grpc.ServiceRegistrar, srv ConfigServiceServer) {
	s.RegisterService(&ConfigService_ServiceDesc, srv)
}

func _ConfigService_UpdateConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceServer).UpdateConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webitel_logger.ConfigService/UpdateConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceServer).UpdateConfig(ctx, req.(*UpdateConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigService_InsertConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceServer).InsertConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webitel_logger.ConfigService/InsertConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceServer).InsertConfig(ctx, req.(*InsertConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigService_GetConfigByObjectId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigByObjectIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceServer).GetConfigByObjectId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webitel_logger.ConfigService/GetConfigByObjectId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceServer).GetConfigByObjectId(ctx, req.(*GetConfigByObjectIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigService_GetConfigById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceServer).GetConfigById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webitel_logger.ConfigService/GetConfigById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceServer).GetConfigById(ctx, req.(*GetConfigByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigService_GetAllConfigs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllConfigsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigServiceServer).GetAllConfigs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/webitel_logger.ConfigService/GetAllConfigs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigServiceServer).GetAllConfigs(ctx, req.(*GetAllConfigsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConfigService_ServiceDesc is the grpc.ServiceDesc for ConfigService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConfigService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "webitel_logger.ConfigService",
	HandlerType: (*ConfigServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateConfig",
			Handler:    _ConfigService_UpdateConfig_Handler,
		},
		{
			MethodName: "InsertConfig",
			Handler:    _ConfigService_InsertConfig_Handler,
		},
		{
			MethodName: "GetConfigByObjectId",
			Handler:    _ConfigService_GetConfigByObjectId_Handler,
		},
		{
			MethodName: "GetConfigById",
			Handler:    _ConfigService_GetConfigById_Handler,
		},
		{
			MethodName: "GetAllConfigs",
			Handler:    _ConfigService_GetAllConfigs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "config_service.proto",
}
