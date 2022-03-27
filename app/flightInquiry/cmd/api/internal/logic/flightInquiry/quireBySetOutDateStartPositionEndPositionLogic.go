package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutDateStartPositionEndPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuireBySetOutDateStartPositionEndPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) QuireBySetOutDateStartPositionEndPositionLogic {
	return QuireBySetOutDateStartPositionEndPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuireBySetOutDateStartPositionEndPositionLogic) QuireBySetOutDateStartPositionEndPosition(req *types.QuireBySetOutDateStartPositionEndPositionReq) (resp *types.QuireBySetOutDateStartPositionEndPositionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
