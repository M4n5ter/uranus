// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.2
// source: app/flightReservation/cmd/rpc/pb/flightReservation.proto

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

// FlightReservationClient is the client API for FlightReservation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FlightReservationClient interface {
	//BookAirTickets 给定： 用户的平台唯一id 航班号 出发日期 是否为头等舱/商务舱 起飞地点/时间 降落地点/时间 来预定机票
	BookAirTickets(ctx context.Context, in *BookAirTicketsReq, opts ...grpc.CallOption) (*BookAirTicketsResp, error)
	//RefundAirTickets 给定：用户的平台唯一id 用户拥有的对应票id 来退订机票
	RefundAirTickets(ctx context.Context, in *RefundAirTicketsReq, opts ...grpc.CallOption) (*RefundAirTicketsResp, error)
	//ChangeAirTickets 给定：用户的平台唯一id 用户拥有的对应票id 目标舱位id
	ChangeAirTickets(ctx context.Context, in *ChangeAirTicketsReq, opts ...grpc.CallOption) (*ChangeAirTicketsResp, error)
}

type flightReservationClient struct {
	cc grpc.ClientConnInterface
}

func NewFlightReservationClient(cc grpc.ClientConnInterface) FlightReservationClient {
	return &flightReservationClient{cc}
}

func (c *flightReservationClient) BookAirTickets(ctx context.Context, in *BookAirTicketsReq, opts ...grpc.CallOption) (*BookAirTicketsResp, error) {
	out := new(BookAirTicketsResp)
	err := c.cc.Invoke(ctx, "/pb.flightReservation/BookAirTickets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *flightReservationClient) RefundAirTickets(ctx context.Context, in *RefundAirTicketsReq, opts ...grpc.CallOption) (*RefundAirTicketsResp, error) {
	out := new(RefundAirTicketsResp)
	err := c.cc.Invoke(ctx, "/pb.flightReservation/RefundAirTickets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *flightReservationClient) ChangeAirTickets(ctx context.Context, in *ChangeAirTicketsReq, opts ...grpc.CallOption) (*ChangeAirTicketsResp, error) {
	out := new(ChangeAirTicketsResp)
	err := c.cc.Invoke(ctx, "/pb.flightReservation/ChangeAirTickets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FlightReservationServer is the server API for FlightReservation service.
// All implementations must embed UnimplementedFlightReservationServer
// for forward compatibility
type FlightReservationServer interface {
	//BookAirTickets 给定： 用户的平台唯一id 航班号 出发日期 是否为头等舱/商务舱 起飞地点/时间 降落地点/时间 来预定机票
	BookAirTickets(context.Context, *BookAirTicketsReq) (*BookAirTicketsResp, error)
	//RefundAirTickets 给定：用户的平台唯一id 用户拥有的对应票id 来退订机票
	RefundAirTickets(context.Context, *RefundAirTicketsReq) (*RefundAirTicketsResp, error)
	//ChangeAirTickets 给定：用户的平台唯一id 用户拥有的对应票id 目标舱位id
	ChangeAirTickets(context.Context, *ChangeAirTicketsReq) (*ChangeAirTicketsResp, error)
	mustEmbedUnimplementedFlightReservationServer()
}

// UnimplementedFlightReservationServer must be embedded to have forward compatible implementations.
type UnimplementedFlightReservationServer struct {
}

func (UnimplementedFlightReservationServer) BookAirTickets(context.Context, *BookAirTicketsReq) (*BookAirTicketsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BookAirTickets not implemented")
}
func (UnimplementedFlightReservationServer) RefundAirTickets(context.Context, *RefundAirTicketsReq) (*RefundAirTicketsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefundAirTickets not implemented")
}
func (UnimplementedFlightReservationServer) ChangeAirTickets(context.Context, *ChangeAirTicketsReq) (*ChangeAirTicketsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangeAirTickets not implemented")
}
func (UnimplementedFlightReservationServer) mustEmbedUnimplementedFlightReservationServer() {}

// UnsafeFlightReservationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FlightReservationServer will
// result in compilation errors.
type UnsafeFlightReservationServer interface {
	mustEmbedUnimplementedFlightReservationServer()
}

func RegisterFlightReservationServer(s grpc.ServiceRegistrar, srv FlightReservationServer) {
	s.RegisterService(&FlightReservation_ServiceDesc, srv)
}

func _FlightReservation_BookAirTickets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookAirTicketsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlightReservationServer).BookAirTickets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.flightReservation/BookAirTickets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlightReservationServer).BookAirTickets(ctx, req.(*BookAirTicketsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FlightReservation_RefundAirTickets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefundAirTicketsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlightReservationServer).RefundAirTickets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.flightReservation/RefundAirTickets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlightReservationServer).RefundAirTickets(ctx, req.(*RefundAirTicketsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FlightReservation_ChangeAirTickets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangeAirTicketsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FlightReservationServer).ChangeAirTickets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.flightReservation/ChangeAirTickets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FlightReservationServer).ChangeAirTickets(ctx, req.(*ChangeAirTicketsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// FlightReservation_ServiceDesc is the grpc.ServiceDesc for FlightReservation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FlightReservation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.flightReservation",
	HandlerType: (*FlightReservationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BookAirTickets",
			Handler:    _FlightReservation_BookAirTickets_Handler,
		},
		{
			MethodName: "RefundAirTickets",
			Handler:    _FlightReservation_RefundAirTickets_Handler,
		},
		{
			MethodName: "ChangeAirTickets",
			Handler:    _FlightReservation_ChangeAirTickets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "app/flightReservation/cmd/rpc/pb/flightReservation.proto",
}
