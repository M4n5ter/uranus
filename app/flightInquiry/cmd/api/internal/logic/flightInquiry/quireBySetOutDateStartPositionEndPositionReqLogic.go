package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutDateStartPositionEndPositionReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuireBySetOutDateStartPositionEndPositionReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) QuireBySetOutDateStartPositionEndPositionReqLogic {
	return QuireBySetOutDateStartPositionEndPositionReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuireBySetOutDateStartPositionEndPositionReqLogic) QuireBySetOutDateStartPositionEndPositionReq(req *types.QuireBySetOutDateStartPositionEndPositionReq) (resp *types.QuireBySetOutDateStartPositionEndPositionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
