package flightReservation

import (
	"context"

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

	return
}
