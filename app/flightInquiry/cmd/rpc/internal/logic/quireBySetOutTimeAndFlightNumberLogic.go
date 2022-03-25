package logic

import (
	"context"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutTimeAndFlightNumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuireBySetOutTimeAndFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireBySetOutTimeAndFlightNumberLogic {
	return &QuireBySetOutTimeAndFlightNumberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuireBySetOutTimeAndFlightNumber 通过给定日期、航班号进行航班查询请求
func (l *QuireBySetOutTimeAndFlightNumberLogic) QuireBySetOutTimeAndFlightNumber(in *pb.QuireBySetOutTimeAndFlightNumberReq) (*pb.QuireBySetOutTimeAndFlightNumberResp, error) {

	return &pb.QuireBySetOutTimeAndFlightNumberResp{}, nil
}
