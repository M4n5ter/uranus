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

func (l *QuireBySetOutTimeAndFlightNumberLogic) QuireBySetOutTimeAndFlightNumber(in *pb.QuireBySetOutTimeAndFlightNumberReq) (*pb.QuireBySetOutTimeAndFlightNumberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.QuireBySetOutTimeAndFlightNumberResp{}, nil
}
