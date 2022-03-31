package logic

import (
	"context"

	"uranus/app/payment/cmd/rpc/internal/svc"
	"uranus/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateLogic {
	return &UpdateTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateTradeState 更新交易状态
func (l *UpdateTradeStateLogic) UpdateTradeState(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateTradeStateResp{}, nil
}
