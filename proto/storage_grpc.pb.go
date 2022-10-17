// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: proto/storage.proto

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

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageClient interface {
	InitBranch(ctx context.Context, in *InitBranchRequest, opts ...grpc.CallOption) (*InitBranchResponse, error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) InitBranch(ctx context.Context, in *InitBranchRequest, opts ...grpc.CallOption) (*InitBranchResponse, error) {
	out := new(InitBranchResponse)
	err := c.cc.Invoke(ctx, "/gok.Storage/InitBranch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServer is the server API for Storage service.
// All implementations must embed UnimplementedStorageServer
// for forward compatibility
type StorageServer interface {
	InitBranch(context.Context, *InitBranchRequest) (*InitBranchResponse, error)
	mustEmbedUnimplementedStorageServer()
}

// UnimplementedStorageServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (UnimplementedStorageServer) InitBranch(context.Context, *InitBranchRequest) (*InitBranchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitBranch not implemented")
}
func (UnimplementedStorageServer) mustEmbedUnimplementedStorageServer() {}

// UnsafeStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServer will
// result in compilation errors.
type UnsafeStorageServer interface {
	mustEmbedUnimplementedStorageServer()
}

func RegisterStorageServer(s grpc.ServiceRegistrar, srv StorageServer) {
	s.RegisterService(&Storage_ServiceDesc, srv)
}

func _Storage_InitBranch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitBranchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).InitBranch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gok.Storage/InitBranch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).InitBranch(ctx, req.(*InitBranchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gok.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitBranch",
			Handler:    _Storage_InitBranch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/storage.proto",
}
