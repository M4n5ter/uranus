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

func (l *QuireBySetOutTimeStartPositionEndPositionLogic) QuireBySetOutTimeStartPositionEndPosition(in *pb.QuireBySetOutTimeStartPositionEndPositionReq) (*pb.QuireBySetOutTimeStartPositionEndPositionResp, error) {
	// todo: add your logic here and delete this line

	return &pb.QuireBySetOutTimeStartPositionEndPositionResp{}, nil
}
