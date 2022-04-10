// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: app/stock/cmd/rpc/pb/stock.proto

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

// StockClient is the client API for Stock service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StockClient interface {
	//通过 ticketID 加库存
	AddStockByTicketID(ctx context.Context, in *AddStockByTicketIDReq, opts ...grpc.CallOption) (*AddStockResp, error)
	//通过 ticketID 加库存 rollback
	AddStockByTicketIDRollBack(ctx context.Context, in *AddStockByTicketIDReq, opts ...grpc.CallOption) (*AddStockResp, error)
	//通过 spaceID 加库存
	AddStockBySpaceID(ctx context.Context, in *AddStockBySpaceIDReq, opts ...grpc.CallOption) (*AddStockResp, error)
	//通过 spaceID 加库存 rollback
	AddStockBySpaceIDRollBack(ctx context.Context, in *AddStockBySpaceIDReq, opts ...grpc.CallOption) (*AddStockResp, error)
	//通过 ticketID 扣库存
	DeductStockByTicketID(ctx context.Context, in *DeductStockByTicketIDReq, opts ...grpc.CallOption) (*DeductStockResp, error)
	//通过 ticketID 扣库存 rollback
	DeductStockByTicketIDRollBack(ctx context.Context, in *DeductStockByTicketIDReq, opts ...grpc.CallOption) (*DeductStockResp, error)
	//通过 spaceID 扣库存
	DeductStockBySpaceID(ctx context.Context, in *DeductStockBySpaceIDReq, opts ...grpc.CallOption) (*DeductStockResp, error)
	//通过 spaceID 扣库存 rollback
	DeductStockBySpaceIDRollBack(ctx context.Context, in *DeductStockBySpaceIDReq, opts ...grpc.CallOption) (*DeductStockResp, error)
}

type stockClient struct {
	cc grpc.ClientConnInterface
}

func NewStockClient(cc grpc.ClientConnInterface) StockClient {
	return &stockClient{cc}
}

func (c *stockClient) AddStockByTicketID(ctx context.Context, in *AddStockByTicketIDReq, opts ...grpc.CallOption) (*AddStockResp, error) {
	out := new(AddStockResp)
	err := c.cc.Invoke(ctx, "/pb.stock/AddStockByTicketID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) AddStockByTicketIDRollBack(ctx context.Context, in *AddStockByTicketIDReq, opts ...grpc.CallOption) (*AddStockResp, error) {
	out := new(AddStockResp)
	err := c.cc.Invoke(ctx, "/pb.stock/AddStockByTicketIDRollBack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) AddStockBySpaceID(ctx context.Context, in *AddStockBySpaceIDReq, opts ...grpc.CallOption) (*AddStockResp, error) {
	out := new(AddStockResp)
	err := c.cc.Invoke(ctx, "/pb.stock/AddStockBySpaceID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) AddStockBySpaceIDRollBack(ctx context.Context, in *AddStockBySpaceIDReq, opts ...grpc.CallOption) (*AddStockResp, error) {
	out := new(AddStockResp)
	err := c.cc.Invoke(ctx, "/pb.stock/AddStockBySpaceIDRollBack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) DeductStockByTicketID(ctx context.Context, in *DeductStockByTicketIDReq, opts ...grpc.CallOption) (*DeductStockResp, error) {
	out := new(DeductStockResp)
	err := c.cc.Invoke(ctx, "/pb.stock/DeductStockByTicketID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) DeductStockByTicketIDRollBack(ctx context.Context, in *DeductStockByTicketIDReq, opts ...grpc.CallOption) (*DeductStockResp, error) {
	out := new(DeductStockResp)
	err := c.cc.Invoke(ctx, "/pb.stock/DeductStockByTicketIDRollBack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) DeductStockBySpaceID(ctx context.Context, in *DeductStockBySpaceIDReq, opts ...grpc.CallOption) (*DeductStockResp, error) {
	out := new(DeductStockResp)
	err := c.cc.Invoke(ctx, "/pb.stock/DeductStockBySpaceID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) DeductStockBySpaceIDRollBack(ctx context.Context, in *DeductStockBySpaceIDReq, opts ...grpc.CallOption) (*DeductStockResp, error) {
	out := new(DeductStockResp)
	err := c.cc.Invoke(ctx, "/pb.stock/DeductStockBySpaceIDRollBack", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StockServer is the server API for Stock service.
// All implementations must embed UnimplementedStockServer
// for forward compatibility
type StockServer interface {
	//通过 ticketID 加库存
	AddStockByTicketID(context.Context, *AddStockByTicketIDReq) (*AddStockResp, error)
	//通过 ticketID 加库存 rollback
	AddStockByTicketIDRollBack(context.Context, *AddStockByTicketIDReq) (*AddStockResp, error)
	//通过 spaceID 加库存
	AddStockBySpaceID(context.Context, *AddStockBySpaceIDReq) (*AddStockResp, error)
	//通过 spaceID 加库存 rollback
	AddStockBySpaceIDRollBack(context.Context, *AddStockBySpaceIDReq) (*AddStockResp, error)
	//通过 ticketID 扣库存
	DeductStockByTicketID(context.Context, *DeductStockByTicketIDReq) (*DeductStockResp, error)
	//通过 ticketID 扣库存 rollback
	DeductStockByTicketIDRollBack(context.Context, *DeductStockByTicketIDReq) (*DeductStockResp, error)
	//通过 spaceID 扣库存
	DeductStockBySpaceID(context.Context, *DeductStockBySpaceIDReq) (*DeductStockResp, error)
	//通过 spaceID 扣库存 rollback
	DeductStockBySpaceIDRollBack(context.Context, *DeductStockBySpaceIDReq) (*DeductStockResp, error)
	mustEmbedUnimplementedStockServer()
}

// UnimplementedStockServer must be embedded to have forward compatible implementations.
type UnimplementedStockServer struct {
}

func (UnimplementedStockServer) AddStockByTicketID(context.Context, *AddStockByTicketIDReq) (*AddStockResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddStockByTicketID not implemented")
}
func (UnimplementedStockServer) AddStockByTicketIDRollBack(context.Context, *AddStockByTicketIDReq) (*AddStockResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddStockByTicketIDRollBack not implemented")
}
func (UnimplementedStockServer) AddStockBySpaceID(context.Context, *AddStockBySpaceIDReq) (*AddStockResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddStockBySpaceID not implemented")
}
func (UnimplementedStockServer) AddStockBySpaceIDRollBack(context.Context, *AddStockBySpaceIDReq) (*AddStockResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddStockBySpaceIDRollBack not implemented")
}
func (UnimplementedStockServer) DeductStockByTicketID(context.Context, *DeductStockByTicketIDReq) (*DeductStockResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeductStockByTicketID not implemented")
}
func (UnimplementedStockServer) DeductStockByTicketIDRollBack(context.Context, *DeductStockByTicketIDReq) (*DeductStockResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeductStockByTicketIDRollBack not implemented")
}
func (UnimplementedStockServer) DeductStockBySpaceID(context.Context, *DeductStockBySpaceIDReq) (*DeductStockResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeductStockBySpaceID not implemented")
}
func (UnimplementedStockServer) DeductStockBySpaceIDRollBack(context.Context, *DeductStockBySpaceIDReq) (*DeductStockResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeductStockBySpaceIDRollBack not implemented")
}
func (UnimplementedStockServer) mustEmbedUnimplementedStockServer() {}

// UnsafeStockServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StockServer will
// result in compilation errors.
type UnsafeStockServer interface {
	mustEmbedUnimplementedStockServer()
}

func RegisterStockServer(s grpc.ServiceRegistrar, srv StockServer) {
	s.RegisterService(&Stock_ServiceDesc, srv)
}

func _Stock_AddStockByTicketID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStockByTicketIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).AddStockByTicketID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.stock/AddStockByTicketID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).AddStockByTicketID(ctx, req.(*AddStockByTicketIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_AddStockByTicketIDRollBack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStockByTicketIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).AddStockByTicketIDRollBack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.stock/AddStockByTicketIDRollBack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).AddStockByTicketIDRollBack(ctx, req.(*AddStockByTicketIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_AddStockBySpaceID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStockBySpaceIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).AddStockBySpaceID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.stock/AddStockBySpaceID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).AddStockBySpaceID(ctx, req.(*AddStockBySpaceIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_AddStockBySpaceIDRollBack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddStockBySpaceIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).AddStockBySpaceIDRollBack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.stock/AddStockBySpaceIDRollBack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).AddStockBySpaceIDRollBack(ctx, req.(*AddStockBySpaceIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_DeductStockByTicketID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeductStockByTicketIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).DeductStockByTicketID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.stock/DeductStockByTicketID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).DeductStockByTicketID(ctx, req.(*DeductStockByTicketIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_DeductStockByTicketIDRollBack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeductStockByTicketIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).DeductStockByTicketIDRollBack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.stock/DeductStockByTicketIDRollBack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).DeductStockByTicketIDRollBack(ctx, req.(*DeductStockByTicketIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_DeductStockBySpaceID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeductStockBySpaceIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).DeductStockBySpaceID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.stock/DeductStockBySpaceID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).DeductStockBySpaceID(ctx, req.(*DeductStockBySpaceIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_DeductStockBySpaceIDRollBack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeductStockBySpaceIDReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).DeductStockBySpaceIDRollBack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.stock/DeductStockBySpaceIDRollBack",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).DeductStockBySpaceIDRollBack(ctx, req.(*DeductStockBySpaceIDReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Stock_ServiceDesc is the grpc.ServiceDesc for Stock service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stock_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.stock",
	HandlerType: (*StockServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddStockByTicketID",
			Handler:    _Stock_AddStockByTicketID_Handler,
		},
		{
			MethodName: "AddStockByTicketIDRollBack",
			Handler:    _Stock_AddStockByTicketIDRollBack_Handler,
		},
		{
			MethodName: "AddStockBySpaceID",
			Handler:    _Stock_AddStockBySpaceID_Handler,
		},
		{
			MethodName: "AddStockBySpaceIDRollBack",
			Handler:    _Stock_AddStockBySpaceIDRollBack_Handler,
		},
		{
			MethodName: "DeductStockByTicketID",
			Handler:    _Stock_DeductStockByTicketID_Handler,
		},
		{
			MethodName: "DeductStockByTicketIDRollBack",
			Handler:    _Stock_DeductStockByTicketIDRollBack_Handler,
		},
		{
			MethodName: "DeductStockBySpaceID",
			Handler:    _Stock_DeductStockBySpaceID_Handler,
		},
		{
			MethodName: "DeductStockBySpaceIDRollBack",
			Handler:    _Stock_DeductStockBySpaceIDRollBack_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "app/stock/cmd/rpc/pb/stock.proto",
}
