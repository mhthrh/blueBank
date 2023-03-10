// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: Proto/Gateway.proto

package ProtoGateway

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

// GatewayServicesClient is the client API for GatewayServices service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GatewayServicesClient interface {
	GatewayLogin(ctx context.Context, in *GatewayLoginRequest, opts ...grpc.CallOption) (*GatewayLoginResponse, error)
}

type gatewayServicesClient struct {
	cc grpc.ClientConnInterface
}

func NewGatewayServicesClient(cc grpc.ClientConnInterface) GatewayServicesClient {
	return &gatewayServicesClient{cc}
}

func (c *gatewayServicesClient) GatewayLogin(ctx context.Context, in *GatewayLoginRequest, opts ...grpc.CallOption) (*GatewayLoginResponse, error) {
	out := new(GatewayLoginResponse)
	err := c.cc.Invoke(ctx, "/GatewayServices/GatewayLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GatewayServicesServer is the server API for GatewayServices service.
// All implementations must embed UnimplementedGatewayServicesServer
// for forward compatibility
type GatewayServicesServer interface {
	GatewayLogin(context.Context, *GatewayLoginRequest) (*GatewayLoginResponse, error)
	mustEmbedUnimplementedGatewayServicesServer()
}

// UnimplementedGatewayServicesServer must be embedded to have forward compatible implementations.
type UnimplementedGatewayServicesServer struct {
}

func (UnimplementedGatewayServicesServer) GatewayLogin(context.Context, *GatewayLoginRequest) (*GatewayLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GatewayLogin not implemented")
}
func (UnimplementedGatewayServicesServer) mustEmbedUnimplementedGatewayServicesServer() {}

// UnsafeGatewayServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GatewayServicesServer will
// result in compilation errors.
type UnsafeGatewayServicesServer interface {
	mustEmbedUnimplementedGatewayServicesServer()
}

func RegisterGatewayServicesServer(s grpc.ServiceRegistrar, srv GatewayServicesServer) {
	s.RegisterService(&GatewayServices_ServiceDesc, srv)
}

func _GatewayServices_GatewayLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GatewayLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GatewayServicesServer).GatewayLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GatewayServices/GatewayLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GatewayServicesServer).GatewayLogin(ctx, req.(*GatewayLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GatewayServices_ServiceDesc is the grpc.ServiceDesc for GatewayServices service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GatewayServices_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "GatewayServices",
	HandlerType: (*GatewayServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GatewayLogin",
			Handler:    _GatewayServices_GatewayLogin_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/Gateway.proto",
}
