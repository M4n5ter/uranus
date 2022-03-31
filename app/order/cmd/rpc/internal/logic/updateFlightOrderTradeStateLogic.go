package logic

import (
	"context"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFlightOrderTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFlightOrderTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFlightOrderTradeStateLogic {
	return &UpdateFlightOrderTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新机票订单状态
func (l *UpdateFlightOrderTradeStateLogic) UpdateFlightOrderTradeState(in *pb.UpdateFlightOrderTradeStateReq) (*pb.UpdateFlightOrderTradeStateResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateFlightOrderTradeStateResp{}, nil
}
