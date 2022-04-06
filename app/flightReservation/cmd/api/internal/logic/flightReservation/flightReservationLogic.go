package flightReservation

import (
	"context"

	"uranus/app/flightReservation/cmd/api/internal/svc"
	"uranus/app/flightReservation/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlightReservationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFlightReservationLogic(ctx context.Context, svcCtx *svc.ServiceContext) FlightReservationLogic {
	return FlightReservationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FlightReservationLogic) FlightReservation(req *types.FlightReservationReq) (resp *types.FlightReservationResp, err error) {
	// todo: add your logic here and delete this line

	return
}
