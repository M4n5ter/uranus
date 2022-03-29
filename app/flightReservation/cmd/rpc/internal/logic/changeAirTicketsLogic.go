package logic

import (
	"context"

	"uranus/app/flightReservation/cmd/rpc/internal/svc"
	"uranus/app/flightReservation/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAirTicketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeAirTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAirTicketsLogic {
	return &ChangeAirTicketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChangeAirTickets 给定：用户的平台唯一id 用户拥有的对应票id 目标舱位id
func (l *ChangeAirTicketsLogic) ChangeAirTickets(in *pb.ChangeAirTicketsReq) (*pb.ChangeAirTicketsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ChangeAirTicketsResp{}, nil
}
