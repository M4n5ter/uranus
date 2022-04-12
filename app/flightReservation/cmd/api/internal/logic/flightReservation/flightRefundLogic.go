package flightReservation

import (
	"context"
	"uranus/app/flightReservation/cmd/api/internal/svc"
	"uranus/app/flightReservation/cmd/api/internal/types"
	"uranus/app/flightReservation/cmd/rpc/flightreservation"
	"uranus/common/ctxdata"

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

	userID := ctxdata.GetUidFromCtx(l.ctx)

	_, err = l.svcCtx.FlightReservationRpcClient.RefundAirTickets(l.ctx, &flightreservation.RefundAirTicketsReq{
		UserID:  userID,
		OrderSn: req.OrderSn,
	})
	if err != nil {
		return nil, err
	}

	return &types.FlightRefundResp{Message: "退票成功"}, nil
}
