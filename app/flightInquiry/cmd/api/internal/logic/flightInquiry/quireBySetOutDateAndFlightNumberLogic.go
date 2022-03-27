package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutDateAndFlightNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuireBySetOutDateAndFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) QuireBySetOutDateAndFlightNumberLogic {
	return QuireBySetOutDateAndFlightNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuireBySetOutDateAndFlightNumberLogic) QuireBySetOutDateAndFlightNumber(req *types.QuireBySetOutDateAndFlightNumberReq) (resp *types.QuireBySetOutDateAndFlightNumberResp, err error) {
	// todo: add your logic here and delete this line

	return
}
