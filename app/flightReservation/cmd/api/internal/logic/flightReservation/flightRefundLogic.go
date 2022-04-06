package flightReservation

import (
	"context"

	"uranus/app/flightReservation/cmd/api/internal/svc"
	"uranus/app/flightReservation/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlightRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlightRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) FlightRefundLogic {
	return FlightRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlightRefundLogic) FlightRefund(req *types.FlightRefundReq) (resp *types.FlightRefundResp, err error) {
	// todo: add your logic here and delete this line

	return
}
