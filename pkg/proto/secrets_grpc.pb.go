// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: proto/secrets.proto

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Secrets_GetUserSecret_FullMethodName    = "/proto.Secrets/GetUserSecret"
	Secrets_GetUserSecrets_FullMethodName   = "/proto.Secrets/GetUserSecrets"
	Secrets_SaveUserSecret_FullMethodName   = "/proto.Secrets/SaveUserSecret"
	Secrets_DeleteUserSecret_FullMethodName = "/proto.Secrets/DeleteUserSecret"
)

// SecretsClient is the client API for Secrets service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretsClient interface {
	GetUserSecret(ctx context.Context, in *GetUserSecretRequest, opts ...grpc.CallOption) (*GetUserSecretResponse, error)
	GetUserSecrets(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetUserSecretsResponse, error)
	SaveUserSecret(ctx context.Context, in *SaveUserSecretRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteUserSecret(ctx context.Context, in *DeleteUserSecretRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type secretsClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretsClient(cc grpc.ClientConnInterface) SecretsClient {
	return &secretsClient{cc}
}

func (c *secretsClient) GetUserSecret(ctx context.Context, in *GetUserSecretRequest, opts ...grpc.CallOption) (*GetUserSecretResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserSecretResponse)
	err := c.cc.Invoke(ctx, Secrets_GetUserSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretsClient) GetUserSecrets(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GetUserSecretsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserSecretsResponse)
	err := c.cc.Invoke(ctx, Secrets_GetUserSecrets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretsClient) SaveUserSecret(ctx context.Context, in *SaveUserSecretRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Secrets_SaveUserSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretsClient) DeleteUserSecret(ctx context.Context, in *DeleteUserSecretRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Secrets_DeleteUserSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretsServer is the server API for Secrets service.
// All implementations must embed UnimplementedSecretsServer
// for forward compatibility.
type SecretsServer interface {
	GetUserSecret(context.Context, *GetUserSecretRequest) (*GetUserSecretResponse, error)
	GetUserSecrets(context.Context, *empty.Empty) (*GetUserSecretsResponse, error)
	SaveUserSecret(context.Context, *SaveUserSecretRequest) (*empty.Empty, error)
	DeleteUserSecret(context.Context, *DeleteUserSecretRequest) (*empty.Empty, error)
	mustEmbedUnimplementedSecretsServer()
}

// UnimplementedSecretsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSecretsServer struct{}

func (UnimplementedSecretsServer) GetUserSecret(context.Context, *GetUserSecretRequest) (*GetUserSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserSecret not implemented")
}
func (UnimplementedSecretsServer) GetUserSecrets(context.Context, *empty.Empty) (*GetUserSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserSecrets not implemented")
}
func (UnimplementedSecretsServer) SaveUserSecret(context.Context, *SaveUserSecretRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveUserSecret not implemented")
}
func (UnimplementedSecretsServer) DeleteUserSecret(context.Context, *DeleteUserSecretRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserSecret not implemented")
}
func (UnimplementedSecretsServer) mustEmbedUnimplementedSecretsServer() {}
func (UnimplementedSecretsServer) testEmbeddedByValue()                 {}

// UnsafeSecretsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretsServer will
// result in compilation errors.
type UnsafeSecretsServer interface {
	mustEmbedUnimplementedSecretsServer()
}

func RegisterSecretsServer(s grpc.ServiceRegistrar, srv SecretsServer) {
	// If the following call pancis, it indicates UnimplementedSecretsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Secrets_ServiceDesc, srv)
}

func _Secrets_GetUserSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretsServer).GetUserSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secrets_GetUserSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretsServer).GetUserSecret(ctx, req.(*GetUserSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Secrets_GetUserSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretsServer).GetUserSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secrets_GetUserSecrets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretsServer).GetUserSecrets(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Secrets_SaveUserSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveUserSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretsServer).SaveUserSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secrets_SaveUserSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretsServer).SaveUserSecret(ctx, req.(*SaveUserSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Secrets_DeleteUserSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretsServer).DeleteUserSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Secrets_DeleteUserSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretsServer).DeleteUserSecret(ctx, req.(*DeleteUserSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Secrets_ServiceDesc is the grpc.ServiceDesc for Secrets service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Secrets_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Secrets",
	HandlerType: (*SecretsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserSecret",
			Handler:    _Secrets_GetUserSecret_Handler,
		},
		{
			MethodName: "GetUserSecrets",
			Handler:    _Secrets_GetUserSecrets_Handler,
		},
		{
			MethodName: "SaveUserSecret",
			Handler:    _Secrets_SaveUserSecret_Handler,
		},
		{
			MethodName: "DeleteUserSecret",
			Handler:    _Secrets_DeleteUserSecret_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/secrets.proto",
}
