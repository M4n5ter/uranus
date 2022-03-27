package logic

import (
	"context"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutDateStartPositionEndPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuireBySetOutDateStartPositionEndPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireBySetOutDateStartPositionEndPositionLogic {
	return &QuireBySetOutDateStartPositionEndPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuireBySetOutDateStartPositionEndPosition 通过给定日期、出发地、目的地进行航班查询请求
func (l *QuireBySetOutDateStartPositionEndPositionLogic) QuireBySetOutDateStartPositionEndPosition(in *pb.QuireBySetOutDateStartPositionEndPositionReq) (*pb.QuireBySetOutDateStartPositionEndPositionResp, error) {
	// todo: add your logic here and delete this line

	return &pb.QuireBySetOutDateStartPositionEndPositionResp{}, nil
}
