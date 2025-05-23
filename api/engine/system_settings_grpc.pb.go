// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: system_settings.proto

package engine

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SystemSettingService_CreateSystemSetting_FullMethodName          = "/engine.SystemSettingService/CreateSystemSetting"
	SystemSettingService_SearchSystemSetting_FullMethodName          = "/engine.SystemSettingService/SearchSystemSetting"
	SystemSettingService_SearchAvailableSystemSetting_FullMethodName = "/engine.SystemSettingService/SearchAvailableSystemSetting"
	SystemSettingService_ReadSystemSetting_FullMethodName            = "/engine.SystemSettingService/ReadSystemSetting"
	SystemSettingService_UpdateSystemSetting_FullMethodName          = "/engine.SystemSettingService/UpdateSystemSetting"
	SystemSettingService_PatchSystemSetting_FullMethodName           = "/engine.SystemSettingService/PatchSystemSetting"
	SystemSettingService_DeleteSystemSetting_FullMethodName          = "/engine.SystemSettingService/DeleteSystemSetting"
)

// SystemSettingServiceClient is the client API for SystemSettingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SystemSettingServiceClient interface {
	CreateSystemSetting(ctx context.Context, in *CreateSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error)
	SearchSystemSetting(ctx context.Context, in *SearchSystemSettingRequest, opts ...grpc.CallOption) (*ListSystemSetting, error)
	SearchAvailableSystemSetting(ctx context.Context, in *SearchAvailableSystemSettingRequest, opts ...grpc.CallOption) (*ListAvailableSystemSetting, error)
	ReadSystemSetting(ctx context.Context, in *ReadSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error)
	UpdateSystemSetting(ctx context.Context, in *UpdateSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error)
	PatchSystemSetting(ctx context.Context, in *PatchSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error)
	DeleteSystemSetting(ctx context.Context, in *DeleteSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error)
}

type systemSettingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSystemSettingServiceClient(cc grpc.ClientConnInterface) SystemSettingServiceClient {
	return &systemSettingServiceClient{cc}
}

func (c *systemSettingServiceClient) CreateSystemSetting(ctx context.Context, in *CreateSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SystemSetting)
	err := c.cc.Invoke(ctx, SystemSettingService_CreateSystemSetting_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemSettingServiceClient) SearchSystemSetting(ctx context.Context, in *SearchSystemSettingRequest, opts ...grpc.CallOption) (*ListSystemSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSystemSetting)
	err := c.cc.Invoke(ctx, SystemSettingService_SearchSystemSetting_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemSettingServiceClient) SearchAvailableSystemSetting(ctx context.Context, in *SearchAvailableSystemSettingRequest, opts ...grpc.CallOption) (*ListAvailableSystemSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAvailableSystemSetting)
	err := c.cc.Invoke(ctx, SystemSettingService_SearchAvailableSystemSetting_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemSettingServiceClient) ReadSystemSetting(ctx context.Context, in *ReadSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SystemSetting)
	err := c.cc.Invoke(ctx, SystemSettingService_ReadSystemSetting_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemSettingServiceClient) UpdateSystemSetting(ctx context.Context, in *UpdateSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SystemSetting)
	err := c.cc.Invoke(ctx, SystemSettingService_UpdateSystemSetting_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemSettingServiceClient) PatchSystemSetting(ctx context.Context, in *PatchSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SystemSetting)
	err := c.cc.Invoke(ctx, SystemSettingService_PatchSystemSetting_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *systemSettingServiceClient) DeleteSystemSetting(ctx context.Context, in *DeleteSystemSettingRequest, opts ...grpc.CallOption) (*SystemSetting, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SystemSetting)
	err := c.cc.Invoke(ctx, SystemSettingService_DeleteSystemSetting_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SystemSettingServiceServer is the server API for SystemSettingService service.
// All implementations must embed UnimplementedSystemSettingServiceServer
// for forward compatibility.
type SystemSettingServiceServer interface {
	CreateSystemSetting(context.Context, *CreateSystemSettingRequest) (*SystemSetting, error)
	SearchSystemSetting(context.Context, *SearchSystemSettingRequest) (*ListSystemSetting, error)
	SearchAvailableSystemSetting(context.Context, *SearchAvailableSystemSettingRequest) (*ListAvailableSystemSetting, error)
	ReadSystemSetting(context.Context, *ReadSystemSettingRequest) (*SystemSetting, error)
	UpdateSystemSetting(context.Context, *UpdateSystemSettingRequest) (*SystemSetting, error)
	PatchSystemSetting(context.Context, *PatchSystemSettingRequest) (*SystemSetting, error)
	DeleteSystemSetting(context.Context, *DeleteSystemSettingRequest) (*SystemSetting, error)
	mustEmbedUnimplementedSystemSettingServiceServer()
}

// UnimplementedSystemSettingServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSystemSettingServiceServer struct{}

func (UnimplementedSystemSettingServiceServer) CreateSystemSetting(context.Context, *CreateSystemSettingRequest) (*SystemSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSystemSetting not implemented")
}
func (UnimplementedSystemSettingServiceServer) SearchSystemSetting(context.Context, *SearchSystemSettingRequest) (*ListSystemSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchSystemSetting not implemented")
}
func (UnimplementedSystemSettingServiceServer) SearchAvailableSystemSetting(context.Context, *SearchAvailableSystemSettingRequest) (*ListAvailableSystemSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchAvailableSystemSetting not implemented")
}
func (UnimplementedSystemSettingServiceServer) ReadSystemSetting(context.Context, *ReadSystemSettingRequest) (*SystemSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadSystemSetting not implemented")
}
func (UnimplementedSystemSettingServiceServer) UpdateSystemSetting(context.Context, *UpdateSystemSettingRequest) (*SystemSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSystemSetting not implemented")
}
func (UnimplementedSystemSettingServiceServer) PatchSystemSetting(context.Context, *PatchSystemSettingRequest) (*SystemSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PatchSystemSetting not implemented")
}
func (UnimplementedSystemSettingServiceServer) DeleteSystemSetting(context.Context, *DeleteSystemSettingRequest) (*SystemSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSystemSetting not implemented")
}
func (UnimplementedSystemSettingServiceServer) mustEmbedUnimplementedSystemSettingServiceServer() {}
func (UnimplementedSystemSettingServiceServer) testEmbeddedByValue()                              {}

// UnsafeSystemSettingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SystemSettingServiceServer will
// result in compilation errors.
type UnsafeSystemSettingServiceServer interface {
	mustEmbedUnimplementedSystemSettingServiceServer()
}

func RegisterSystemSettingServiceServer(s grpc.ServiceRegistrar, srv SystemSettingServiceServer) {
	// If the following call pancis, it indicates UnimplementedSystemSettingServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SystemSettingService_ServiceDesc, srv)
}

func _SystemSettingService_CreateSystemSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSystemSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemSettingServiceServer).CreateSystemSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SystemSettingService_CreateSystemSetting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemSettingServiceServer).CreateSystemSetting(ctx, req.(*CreateSystemSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SystemSettingService_SearchSystemSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchSystemSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemSettingServiceServer).SearchSystemSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SystemSettingService_SearchSystemSetting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemSettingServiceServer).SearchSystemSetting(ctx, req.(*SearchSystemSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SystemSettingService_SearchAvailableSystemSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchAvailableSystemSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemSettingServiceServer).SearchAvailableSystemSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SystemSettingService_SearchAvailableSystemSetting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemSettingServiceServer).SearchAvailableSystemSetting(ctx, req.(*SearchAvailableSystemSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SystemSettingService_ReadSystemSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadSystemSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemSettingServiceServer).ReadSystemSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SystemSettingService_ReadSystemSetting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemSettingServiceServer).ReadSystemSetting(ctx, req.(*ReadSystemSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SystemSettingService_UpdateSystemSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSystemSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemSettingServiceServer).UpdateSystemSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SystemSettingService_UpdateSystemSetting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemSettingServiceServer).UpdateSystemSetting(ctx, req.(*UpdateSystemSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SystemSettingService_PatchSystemSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PatchSystemSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemSettingServiceServer).PatchSystemSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SystemSettingService_PatchSystemSetting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemSettingServiceServer).PatchSystemSetting(ctx, req.(*PatchSystemSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SystemSettingService_DeleteSystemSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSystemSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SystemSettingServiceServer).DeleteSystemSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SystemSettingService_DeleteSystemSetting_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SystemSettingServiceServer).DeleteSystemSetting(ctx, req.(*DeleteSystemSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SystemSettingService_ServiceDesc is the grpc.ServiceDesc for SystemSettingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SystemSettingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "engine.SystemSettingService",
	HandlerType: (*SystemSettingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSystemSetting",
			Handler:    _SystemSettingService_CreateSystemSetting_Handler,
		},
		{
			MethodName: "SearchSystemSetting",
			Handler:    _SystemSettingService_SearchSystemSetting_Handler,
		},
		{
			MethodName: "SearchAvailableSystemSetting",
			Handler:    _SystemSettingService_SearchAvailableSystemSetting_Handler,
		},
		{
			MethodName: "ReadSystemSetting",
			Handler:    _SystemSettingService_ReadSystemSetting_Handler,
		},
		{
			MethodName: "UpdateSystemSetting",
			Handler:    _SystemSettingService_UpdateSystemSetting_Handler,
		},
		{
			MethodName: "PatchSystemSetting",
			Handler:    _SystemSettingService_PatchSystemSetting_Handler,
		},
		{
			MethodName: "DeleteSystemSetting",
			Handler:    _SystemSettingService_DeleteSystemSetting_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "system_settings.proto",
}
