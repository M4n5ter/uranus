package logic

import (
	"context"

	"uranus/app/flightReservation/cmd/rpc/internal/svc"
	"uranus/app/flightReservation/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookAirTicketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBookAirTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookAirTicketsLogic {
	return &BookAirTicketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BookAirTickets 给定：航班号 出发日期 是否为头等舱/商务舱 起飞地点/时间 降落地点/时间 来预定机票
func (l *BookAirTicketsLogic) BookAirTickets(in *pb.BookAirTicketsReq) (*pb.BookAirTicketsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BookAirTicketsResp{}, nil
}
