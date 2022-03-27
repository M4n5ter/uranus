package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DatePositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDatePositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) DatePositionLogic {
	return DatePositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DatePositionLogic) DatePosition(req *types.QuireBySetOutDateStartPositionEndPositionReq) (resp *types.QuireBySetOutDateStartPositionEndPositionResp, err error) {
	// todo: add your logic here and delete this line

	return
}
