package flightReservation

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/flightReservation/cmd/rpc/flightreservation"
	"uranus/common/ctxdata"
	"uranus/common/xerr"

	"uranus/app/flightReservation/cmd/api/internal/svc"
	"uranus/app/flightReservation/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlightChangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlightChangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) FlightChangeLogic {
	return FlightChangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlightChangeLogic) FlightChange(req *types.FlightChangeReq) (resp *types.FlightChangeResp, err error) {
	if len(req.OrderSn) == 0 || req.SpaceID < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "")
	}

	userID := ctxdata.GetUidFromCtx(l.ctx)
	rpcResp, err := l.svcCtx.FlightReservationRpcClient.ChangeAirTickets(l.ctx, &flightreservation.ChangeAirTicketsReq{
		UserID:  userID,
		OrderSn: req.OrderSn,
		SpaceID: req.SpaceID,
	})
	if err != nil {
		return nil, err
	}

	return &types.FlightChangeResp{NewOrderSn: rpcResp.OrderSn}, nil
}
