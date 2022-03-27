package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DateNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDateNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) DateNumberLogic {
	return DateNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DateNumberLogic) DateNumber(req *types.QuireBySetOutDateAndFlightNumberReq) (resp *types.QuireBySetOutDateAndFlightNumberResp, err error) {
	// todo: add your logic here and delete this line

	return
}
