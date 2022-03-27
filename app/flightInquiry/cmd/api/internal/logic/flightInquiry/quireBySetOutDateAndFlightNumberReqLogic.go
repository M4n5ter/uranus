package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutDateAndFlightNumberReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuireBySetOutDateAndFlightNumberReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) QuireBySetOutDateAndFlightNumberReqLogic {
	return QuireBySetOutDateAndFlightNumberReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuireBySetOutDateAndFlightNumberReqLogic) QuireBySetOutDateAndFlightNumberReq(req *types.QuireBySetOutDateAndFlightNumberReq) (resp *types.QuireBySetOutDateAndFlightNumberResp, err error) {
	// todo: add your logic here and delete this line

	return
}
