// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: file.proto

package storage

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
	FileService_UploadFile_FullMethodName           = "/storage.FileService/UploadFile"
	FileService_SafeUploadFile_FullMethodName       = "/storage.FileService/SafeUploadFile"
	FileService_DownloadFile_FullMethodName         = "/storage.FileService/DownloadFile"
	FileService_UploadFileUrl_FullMethodName        = "/storage.FileService/UploadFileUrl"
	FileService_GenerateFileLink_FullMethodName     = "/storage.FileService/GenerateFileLink"
	FileService_BulkGenerateFileLink_FullMethodName = "/storage.FileService/BulkGenerateFileLink"
	FileService_DeleteFiles_FullMethodName          = "/storage.FileService/DeleteFiles"
	FileService_SearchFiles_FullMethodName          = "/storage.FileService/SearchFiles"
)

// FileServiceClient is the client API for FileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileServiceClient interface {
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadFileRequest, UploadFileResponse], error)
	SafeUploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[SafeUploadFileRequest, SafeUploadFileResponse], error)
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StreamFile], error)
	UploadFileUrl(ctx context.Context, in *UploadFileUrlRequest, opts ...grpc.CallOption) (*UploadFileUrlResponse, error)
	GenerateFileLink(ctx context.Context, in *GenerateFileLinkRequest, opts ...grpc.CallOption) (*GenerateFileLinkResponse, error)
	BulkGenerateFileLink(ctx context.Context, in *BulkGenerateFileLinkRequest, opts ...grpc.CallOption) (*BulkGenerateFileLinkResponse, error)
	DeleteFiles(ctx context.Context, in *DeleteFilesRequest, opts ...grpc.CallOption) (*DeleteFilesResponse, error)
	SearchFiles(ctx context.Context, in *SearchFilesRequest, opts ...grpc.CallOption) (*ListFile, error)
}

type fileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileServiceClient(cc grpc.ClientConnInterface) FileServiceClient {
	return &fileServiceClient{cc}
}

func (c *fileServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadFileRequest, UploadFileResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[0], FileService_UploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadFileRequest, UploadFileResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_UploadFileClient = grpc.ClientStreamingClient[UploadFileRequest, UploadFileResponse]

func (c *fileServiceClient) SafeUploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[SafeUploadFileRequest, SafeUploadFileResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[1], FileService_SafeUploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SafeUploadFileRequest, SafeUploadFileResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_SafeUploadFileClient = grpc.BidiStreamingClient[SafeUploadFileRequest, SafeUploadFileResponse]

func (c *fileServiceClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[StreamFile], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FileService_ServiceDesc.Streams[2], FileService_DownloadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadFileRequest, StreamFile]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_DownloadFileClient = grpc.ServerStreamingClient[StreamFile]

func (c *fileServiceClient) UploadFileUrl(ctx context.Context, in *UploadFileUrlRequest, opts ...grpc.CallOption) (*UploadFileUrlResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadFileUrlResponse)
	err := c.cc.Invoke(ctx, FileService_UploadFileUrl_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) GenerateFileLink(ctx context.Context, in *GenerateFileLinkRequest, opts ...grpc.CallOption) (*GenerateFileLinkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateFileLinkResponse)
	err := c.cc.Invoke(ctx, FileService_GenerateFileLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) BulkGenerateFileLink(ctx context.Context, in *BulkGenerateFileLinkRequest, opts ...grpc.CallOption) (*BulkGenerateFileLinkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BulkGenerateFileLinkResponse)
	err := c.cc.Invoke(ctx, FileService_BulkGenerateFileLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) DeleteFiles(ctx context.Context, in *DeleteFilesRequest, opts ...grpc.CallOption) (*DeleteFilesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteFilesResponse)
	err := c.cc.Invoke(ctx, FileService_DeleteFiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileServiceClient) SearchFiles(ctx context.Context, in *SearchFilesRequest, opts ...grpc.CallOption) (*ListFile, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListFile)
	err := c.cc.Invoke(ctx, FileService_SearchFiles_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileServiceServer is the server API for FileService service.
// All implementations must embed UnimplementedFileServiceServer
// for forward compatibility.
type FileServiceServer interface {
	UploadFile(grpc.ClientStreamingServer[UploadFileRequest, UploadFileResponse]) error
	SafeUploadFile(grpc.BidiStreamingServer[SafeUploadFileRequest, SafeUploadFileResponse]) error
	DownloadFile(*DownloadFileRequest, grpc.ServerStreamingServer[StreamFile]) error
	UploadFileUrl(context.Context, *UploadFileUrlRequest) (*UploadFileUrlResponse, error)
	GenerateFileLink(context.Context, *GenerateFileLinkRequest) (*GenerateFileLinkResponse, error)
	BulkGenerateFileLink(context.Context, *BulkGenerateFileLinkRequest) (*BulkGenerateFileLinkResponse, error)
	DeleteFiles(context.Context, *DeleteFilesRequest) (*DeleteFilesResponse, error)
	SearchFiles(context.Context, *SearchFilesRequest) (*ListFile, error)
	mustEmbedUnimplementedFileServiceServer()
}

// UnimplementedFileServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFileServiceServer struct{}

func (UnimplementedFileServiceServer) UploadFile(grpc.ClientStreamingServer[UploadFileRequest, UploadFileResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedFileServiceServer) SafeUploadFile(grpc.BidiStreamingServer[SafeUploadFileRequest, SafeUploadFileResponse]) error {
	return status.Errorf(codes.Unimplemented, "method SafeUploadFile not implemented")
}
func (UnimplementedFileServiceServer) DownloadFile(*DownloadFileRequest, grpc.ServerStreamingServer[StreamFile]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedFileServiceServer) UploadFileUrl(context.Context, *UploadFileUrlRequest) (*UploadFileUrlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFileUrl not implemented")
}
func (UnimplementedFileServiceServer) GenerateFileLink(context.Context, *GenerateFileLinkRequest) (*GenerateFileLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateFileLink not implemented")
}
func (UnimplementedFileServiceServer) BulkGenerateFileLink(context.Context, *BulkGenerateFileLinkRequest) (*BulkGenerateFileLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BulkGenerateFileLink not implemented")
}
func (UnimplementedFileServiceServer) DeleteFiles(context.Context, *DeleteFilesRequest) (*DeleteFilesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFiles not implemented")
}
func (UnimplementedFileServiceServer) SearchFiles(context.Context, *SearchFilesRequest) (*ListFile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchFiles not implemented")
}
func (UnimplementedFileServiceServer) mustEmbedUnimplementedFileServiceServer() {}
func (UnimplementedFileServiceServer) testEmbeddedByValue()                     {}

// UnsafeFileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileServiceServer will
// result in compilation errors.
type UnsafeFileServiceServer interface {
	mustEmbedUnimplementedFileServiceServer()
}

func RegisterFileServiceServer(s grpc.ServiceRegistrar, srv FileServiceServer) {
	// If the following call pancis, it indicates UnimplementedFileServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FileService_ServiceDesc, srv)
}

func _FileService_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServiceServer).UploadFile(&grpc.GenericServerStream[UploadFileRequest, UploadFileResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_UploadFileServer = grpc.ClientStreamingServer[UploadFileRequest, UploadFileResponse]

func _FileService_SafeUploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileServiceServer).SafeUploadFile(&grpc.GenericServerStream[SafeUploadFileRequest, SafeUploadFileResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_SafeUploadFileServer = grpc.BidiStreamingServer[SafeUploadFileRequest, SafeUploadFileResponse]

func _FileService_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileServiceServer).DownloadFile(m, &grpc.GenericServerStream[DownloadFileRequest, StreamFile]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FileService_DownloadFileServer = grpc.ServerStreamingServer[StreamFile]

func _FileService_UploadFileUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadFileUrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).UploadFileUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_UploadFileUrl_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).UploadFileUrl(ctx, req.(*UploadFileUrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_GenerateFileLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateFileLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).GenerateFileLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_GenerateFileLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).GenerateFileLink(ctx, req.(*GenerateFileLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_BulkGenerateFileLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkGenerateFileLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).BulkGenerateFileLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_BulkGenerateFileLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).BulkGenerateFileLink(ctx, req.(*BulkGenerateFileLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_DeleteFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).DeleteFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_DeleteFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).DeleteFiles(ctx, req.(*DeleteFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileService_SearchFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileServiceServer).SearchFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FileService_SearchFiles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileServiceServer).SearchFiles(ctx, req.(*SearchFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileService_ServiceDesc is the grpc.ServiceDesc for FileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storage.FileService",
	HandlerType: (*FileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UploadFileUrl",
			Handler:    _FileService_UploadFileUrl_Handler,
		},
		{
			MethodName: "GenerateFileLink",
			Handler:    _FileService_GenerateFileLink_Handler,
		},
		{
			MethodName: "BulkGenerateFileLink",
			Handler:    _FileService_BulkGenerateFileLink_Handler,
		},
		{
			MethodName: "DeleteFiles",
			Handler:    _FileService_DeleteFiles_Handler,
		},
		{
			MethodName: "SearchFiles",
			Handler:    _FileService_SearchFiles_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _FileService_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "SafeUploadFile",
			Handler:       _FileService_SafeUploadFile_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFile",
			Handler:       _FileService_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "file.proto",
}
