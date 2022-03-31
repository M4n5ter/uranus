package logic

import (
	"context"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFlightOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFlightOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFlightOrderLogic {
	return &CreateFlightOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 机票下订单
func (l *CreateFlightOrderLogic) CreateFlightOrder(in *pb.CreateFlightOrderReq) (*pb.CreateFlightOrderResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateFlightOrderResp{}, nil
}
