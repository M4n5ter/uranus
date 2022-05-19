package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireTransferFlightsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuireTransferFlightsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireTransferFlightsLogic {
	return &QuireTransferFlightsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuireTransferFlightsLogic) QuireTransferFlights(req *types.QuireTransferFlightsReq) (resp *types.QuireTransferFlightsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
