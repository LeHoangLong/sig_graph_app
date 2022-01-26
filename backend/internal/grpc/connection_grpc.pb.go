// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package message

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

// ConnectionGrpcClient is the client API for ConnectionGrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConnectionGrpcClient interface {
	GetConnectionProfile(ctx context.Context, in *GetConnectionProfileRequest, opts ...grpc.CallOption) (*ConnectionProfile, error)
	SaveConnectionProfile(ctx context.Context, in *SaveConnectionProfileRequest, opts ...grpc.CallOption) (*ConnectionProfile, error)
	Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectResponse, error)
}

type connectionGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewConnectionGrpcClient(cc grpc.ClientConnInterface) ConnectionGrpcClient {
	return &connectionGrpcClient{cc}
}

func (c *connectionGrpcClient) GetConnectionProfile(ctx context.Context, in *GetConnectionProfileRequest, opts ...grpc.CallOption) (*ConnectionProfile, error) {
	out := new(ConnectionProfile)
	err := c.cc.Invoke(ctx, "/dashboard.ConnectionGrpc/GetConnectionProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectionGrpcClient) SaveConnectionProfile(ctx context.Context, in *SaveConnectionProfileRequest, opts ...grpc.CallOption) (*ConnectionProfile, error) {
	out := new(ConnectionProfile)
	err := c.cc.Invoke(ctx, "/dashboard.ConnectionGrpc/SaveConnectionProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectionGrpcClient) Connect(ctx context.Context, in *ConnectRequest, opts ...grpc.CallOption) (*ConnectResponse, error) {
	out := new(ConnectResponse)
	err := c.cc.Invoke(ctx, "/dashboard.ConnectionGrpc/Connect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectionGrpcServer is the server API for ConnectionGrpc service.
// All implementations must embed UnimplementedConnectionGrpcServer
// for forward compatibility
type ConnectionGrpcServer interface {
	GetConnectionProfile(context.Context, *GetConnectionProfileRequest) (*ConnectionProfile, error)
	SaveConnectionProfile(context.Context, *SaveConnectionProfileRequest) (*ConnectionProfile, error)
	Connect(context.Context, *ConnectRequest) (*ConnectResponse, error)
	mustEmbedUnimplementedConnectionGrpcServer()
}

// UnimplementedConnectionGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedConnectionGrpcServer struct {
}

func (UnimplementedConnectionGrpcServer) GetConnectionProfile(context.Context, *GetConnectionProfileRequest) (*ConnectionProfile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnectionProfile not implemented")
}
func (UnimplementedConnectionGrpcServer) SaveConnectionProfile(context.Context, *SaveConnectionProfileRequest) (*ConnectionProfile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveConnectionProfile not implemented")
}
func (UnimplementedConnectionGrpcServer) Connect(context.Context, *ConnectRequest) (*ConnectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedConnectionGrpcServer) mustEmbedUnimplementedConnectionGrpcServer() {}

// UnsafeConnectionGrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConnectionGrpcServer will
// result in compilation errors.
type UnsafeConnectionGrpcServer interface {
	mustEmbedUnimplementedConnectionGrpcServer()
}

func RegisterConnectionGrpcServer(s grpc.ServiceRegistrar, srv ConnectionGrpcServer) {
	s.RegisterService(&ConnectionGrpc_ServiceDesc, srv)
}

func _ConnectionGrpc_GetConnectionProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConnectionProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionGrpcServer).GetConnectionProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.ConnectionGrpc/GetConnectionProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionGrpcServer).GetConnectionProfile(ctx, req.(*GetConnectionProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectionGrpc_SaveConnectionProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveConnectionProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionGrpcServer).SaveConnectionProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.ConnectionGrpc/SaveConnectionProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionGrpcServer).SaveConnectionProfile(ctx, req.(*SaveConnectionProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConnectionGrpc_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionGrpcServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dashboard.ConnectionGrpc/Connect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionGrpcServer).Connect(ctx, req.(*ConnectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ConnectionGrpc_ServiceDesc is the grpc.ServiceDesc for ConnectionGrpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConnectionGrpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dashboard.ConnectionGrpc",
	HandlerType: (*ConnectionGrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConnectionProfile",
			Handler:    _ConnectionGrpc_GetConnectionProfile_Handler,
		},
		{
			MethodName: "SaveConnectionProfile",
			Handler:    _ConnectionGrpc_SaveConnectionProfile_Handler,
		},
		{
			MethodName: "Connect",
			Handler:    _ConnectionGrpc_Connect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "connection.proto",
}
