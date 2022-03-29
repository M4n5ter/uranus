package logic

import (
	"context"

	"uranus/app/flightReservation/cmd/rpc/internal/svc"
	"uranus/app/flightReservation/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundAirTicketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundAirTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundAirTicketsLogic {
	return &RefundAirTicketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RefundAirTickets 给定：用户的平台唯一id 用户拥有的对应票id 来退订机票
func (l *RefundAirTicketsLogic) RefundAirTickets(in *pb.RefundAirTicketsReq) (*pb.RefundAirTicketsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.RefundAirTicketsResp{}, nil
}
