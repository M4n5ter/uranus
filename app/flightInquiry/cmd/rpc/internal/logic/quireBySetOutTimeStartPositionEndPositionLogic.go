package logic

import (
	"context"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutTimeStartPositionEndPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuireBySetOutTimeStartPositionEndPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireBySetOutTimeStartPositionEndPositionLogic {
	return &QuireBySetOutTimeStartPositionEndPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuireBySetOutTimeStartPositionEndPosition 通过给定日期、出发地、目的地进行航班查询请求
func (l *QuireBySetOutTimeStartPositionEndPositionLogic) QuireBySetOutTimeStartPositionEndPosition(in *pb.QuireBySetOutTimeStartPositionEndPositionReq) (*pb.QuireBySetOutTimeStartPositionEndPositionResp, error) {

	return &pb.QuireBySetOutTimeStartPositionEndPositionResp{}, nil
}
