// Code generated by goctl. DO NOT EDIT!
// Source: order.proto

package order

import (
	"context"

	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateFlightOrderReq            = pb.CreateFlightOrderReq
	CreateFlightOrderResp           = pb.CreateFlightOrderResp
	FlightOrder                     = pb.FlightOrder
	FlightOrderDetailReq            = pb.FlightOrderDetailReq
	FlightOrderDetailResp           = pb.FlightOrderDetailResp
	UpdateFlightOrderTradeStateReq  = pb.UpdateFlightOrderTradeStateReq
	UpdateFlightOrderTradeStateResp = pb.UpdateFlightOrderTradeStateResp
	UserFlightOrderListReq          = pb.UserFlightOrderListReq
	UserFlightOrderListResp         = pb.UserFlightOrderListResp

	Order interface {
		// 机票下订单
		CreateFlightOrder(ctx context.Context, in *CreateFlightOrderReq, opts ...grpc.CallOption) (*CreateFlightOrderResp, error)
		// 机票订单详情
		FlightOrderDetail(ctx context.Context, in *FlightOrderDetailReq, opts ...grpc.CallOption) (*FlightOrderDetailResp, error)
		// 更新机票订单状态
		UpdateFlightOrderTradeState(ctx context.Context, in *UpdateFlightOrderTradeStateReq, opts ...grpc.CallOption) (*UpdateFlightOrderTradeStateResp, error)
		// 用户机票订单
		UserFlightOrderList(ctx context.Context, in *UserFlightOrderListReq, opts ...grpc.CallOption) (*UserFlightOrderListResp, error)
	}

	defaultOrder struct {
		cli zrpc.Client
	}
)

func NewOrder(cli zrpc.Client) Order {
	return &defaultOrder{
		cli: cli,
	}
}

// 机票下订单
func (m *defaultOrder) CreateFlightOrder(ctx context.Context, in *CreateFlightOrderReq, opts ...grpc.CallOption) (*CreateFlightOrderResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.CreateFlightOrder(ctx, in, opts...)
}

// 机票订单详情
func (m *defaultOrder) FlightOrderDetail(ctx context.Context, in *FlightOrderDetailReq, opts ...grpc.CallOption) (*FlightOrderDetailResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.FlightOrderDetail(ctx, in, opts...)
}

// 更新机票订单状态
func (m *defaultOrder) UpdateFlightOrderTradeState(ctx context.Context, in *UpdateFlightOrderTradeStateReq, opts ...grpc.CallOption) (*UpdateFlightOrderTradeStateResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.UpdateFlightOrderTradeState(ctx, in, opts...)
}

// 用户机票订单
func (m *defaultOrder) UserFlightOrderList(ctx context.Context, in *UserFlightOrderListReq, opts ...grpc.CallOption) (*UserFlightOrderListResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.UserFlightOrderList(ctx, in, opts...)
}
