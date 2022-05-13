package logic

import (
	"context"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireTransferFlightsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuireTransferFlightsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireTransferFlightsLogic {
	return &QuireTransferFlightsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuireTransferFlights 中转航班查询
func (l *QuireTransferFlightsLogic) QuireTransferFlights(in *pb.QuireTransferFlightsReq) (*pb.QuireTransferFlightsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.QuireTransferFlightsResp{}, nil
}
