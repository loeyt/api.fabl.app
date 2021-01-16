// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	CurrentAccount(ctx context.Context, in *CurrentAccountRequest, opts ...grpc.CallOption) (*CurrentAccountResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) CurrentAccount(ctx context.Context, in *CurrentAccountRequest, opts ...grpc.CallOption) (*CurrentAccountResponse, error) {
	out := new(CurrentAccountResponse)
	err := c.cc.Invoke(ctx, "/fabl.v1.AccountService/CurrentAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/fabl.v1.AccountService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	out := new(LogoutResponse)
	err := c.cc.Invoke(ctx, "/fabl.v1.AccountService/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations must embed UnimplementedAccountServiceServer
// for forward compatibility
type AccountServiceServer interface {
	CurrentAccount(context.Context, *CurrentAccountRequest) (*CurrentAccountResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	mustEmbedUnimplementedAccountServiceServer()
}

// UnimplementedAccountServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccountServiceServer struct {
}

func (UnimplementedAccountServiceServer) CurrentAccount(context.Context, *CurrentAccountRequest) (*CurrentAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CurrentAccount not implemented")
}
func (UnimplementedAccountServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAccountServiceServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_CurrentAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CurrentAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).CurrentAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fabl.v1.AccountService/CurrentAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).CurrentAccount(ctx, req.(*CurrentAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fabl.v1.AccountService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/fabl.v1.AccountService/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fabl.v1.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CurrentAccount",
			Handler:    _AccountService_CurrentAccount_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _AccountService_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _AccountService_Logout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/account_service.proto",
}
