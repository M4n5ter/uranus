package logic

import (
	"context"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlightOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlightOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlightOrderDetailLogic {
	return &FlightOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 机票订单详情
func (l *FlightOrderDetailLogic) FlightOrderDetail(in *pb.FlightOrderDetailReq) (*pb.FlightOrderDetailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.FlightOrderDetailResp{}, nil
}
