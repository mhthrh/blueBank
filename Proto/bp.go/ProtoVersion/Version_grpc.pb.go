// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: Proto/Version.proto

package ProtoVersion

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

// VersionServicesClient is the client API for VersionServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VersionServicesClient interface {
	GetVersion(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error)
}

type versionServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewVersionServicesClient(cc grpc.ClientConnInterface) VersionServicesClient {
	return &versionServicesClient{cc}
}

func (c *versionServicesClient) GetVersion(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/VersionServices/GetVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VersionServicesServer is the server API for VersionServices service.
// All implementations must embed UnimplementedVersionServicesServer
// for forward compatibility
type VersionServicesServer interface {
	GetVersion(context.Context, *VersionRequest) (*VersionResponse, error)
	mustEmbedUnimplementedVersionServicesServer()
}

// UnimplementedVersionServicesServer must be embedded to have forward compatible implementations.
type UnimplementedVersionServicesServer struct {
}

func (UnimplementedVersionServicesServer) GetVersion(context.Context, *VersionRequest) (*VersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVersion not implemented")
}
func (UnimplementedVersionServicesServer) mustEmbedUnimplementedVersionServicesServer() {}

// UnsafeVersionServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VersionServicesServer will
// result in compilation errors.
type UnsafeVersionServicesServer interface {
	mustEmbedUnimplementedVersionServicesServer()
}

func RegisterVersionServicesServer(s grpc.ServiceRegistrar, srv VersionServicesServer) {
	s.RegisterService(&VersionServices_ServiceDesc, srv)
}

func _VersionServices_GetVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VersionServicesServer).GetVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/VersionServices/GetVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VersionServicesServer).GetVersion(ctx, req.(*VersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// VersionServices_ServiceDesc is the grpc.ServiceDesc for VersionServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VersionServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "VersionServices",
	HandlerType: (*VersionServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVersion",
			Handler:    _VersionServices_GetVersion_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/Version.proto",
}
