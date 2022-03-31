package logic

import (
	"context"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFlightOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserFlightOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFlightOrderListLogic {
	return &UserFlightOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户机票订单
func (l *UserFlightOrderListLogic) UserFlightOrderList(in *pb.UserFlightOrderListReq) (*pb.UserFlightOrderListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserFlightOrderListResp{}, nil
}
