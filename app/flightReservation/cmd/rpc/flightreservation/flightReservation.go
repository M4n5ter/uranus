// Code generated by goctl. DO NOT EDIT!
// Source: flightReservation.proto

package flightreservation

import (
	"context"

	"uranus/app/flightReservation/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	BookAirTicketsReq    = pb.BookAirTicketsReq
	BookAirTicketsResp   = pb.BookAirTicketsResp
	ChangeAirTicketsReq  = pb.ChangeAirTicketsReq
	ChangeAirTicketsResp = pb.ChangeAirTicketsResp
	RefundAirTicketsReq  = pb.RefundAirTicketsReq
	RefundAirTicketsResp = pb.RefundAirTicketsResp

	FlightReservation interface {
		//  BookAirTickets 给定： 用户的平台唯一id 航班号 出发日期 是否为头等舱/商务舱 起飞地点/时间 降落地点/时间 来预定机票
		BookAirTickets(ctx context.Context, in *BookAirTicketsReq, opts ...grpc.CallOption) (*BookAirTicketsResp, error)
		//  RefundAirTickets 给定：用户的平台唯一id 用户拥有的对应订单号 来退订机票
		RefundAirTickets(ctx context.Context, in *RefundAirTicketsReq, opts ...grpc.CallOption) (*RefundAirTicketsResp, error)
		//  ChangeAirTickets 给定：用户的平台唯一id 用户拥有的对应订单号 目标舱位id
		ChangeAirTickets(ctx context.Context, in *ChangeAirTicketsReq, opts ...grpc.CallOption) (*ChangeAirTicketsResp, error)
	}

	defaultFlightReservation struct {
		cli zrpc.Client
	}
)

func NewFlightReservation(cli zrpc.Client) FlightReservation {
	return &defaultFlightReservation{
		cli: cli,
	}
}

//  BookAirTickets 给定： 用户的平台唯一id 航班号 出发日期 是否为头等舱/商务舱 起飞地点/时间 降落地点/时间 来预定机票
func (m *defaultFlightReservation) BookAirTickets(ctx context.Context, in *BookAirTicketsReq, opts ...grpc.CallOption) (*BookAirTicketsResp, error) {
	client := pb.NewFlightReservationClient(m.cli.Conn())
	return client.BookAirTickets(ctx, in, opts...)
}

//  RefundAirTickets 给定：用户的平台唯一id 用户拥有的对应订单号 来退订机票
func (m *defaultFlightReservation) RefundAirTickets(ctx context.Context, in *RefundAirTicketsReq, opts ...grpc.CallOption) (*RefundAirTicketsResp, error) {
	client := pb.NewFlightReservationClient(m.cli.Conn())
	return client.RefundAirTickets(ctx, in, opts...)
}

//  ChangeAirTickets 给定：用户的平台唯一id 用户拥有的对应订单号 目标舱位id
func (m *defaultFlightReservation) ChangeAirTickets(ctx context.Context, in *ChangeAirTicketsReq, opts ...grpc.CallOption) (*ChangeAirTicketsResp, error) {
	client := pb.NewFlightReservationClient(m.cli.Conn())
	return client.ChangeAirTickets(ctx, in, opts...)
}
