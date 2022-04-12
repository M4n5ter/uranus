package flightReservation

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/flightReservation/cmd/rpc/flightreservation"
	"uranus/common/ctxdata"
	"uranus/common/timeTools"
	"uranus/common/xerr"

	"uranus/app/flightReservation/cmd/api/internal/svc"
	"uranus/app/flightReservation/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRInvalidInput = xerr.NewErrMsg("非法输入")
var ERRBookFail = xerr.NewErrMsg("预定机票失败")

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
	// 检查输入
	if len(req.FlightNumber) == 0 || len(req.DepartTime) == 0 || len(req.ArriveTime) == 0 || len(req.ArrivePosition) == 0 || len(req.DepartPosition) == 0 || len(req.SetOutDate) == 0 {
		return nil, errors.Wrapf(ERRInvalidInput, "")
	}

	sod, err := timeTools.String2TimeYMD(req.SetOutDate)
	if err != nil {
		return nil, errors.Wrapf(ERRInvalidInput, "invalid input SOD err: %v", err)
	}

	arriveTime, err := timeTools.String2TimeYMDhms(req.ArriveTime)
	if err != nil {
		return nil, errors.Wrapf(ERRInvalidInput, "invalid input arriveTime err: %v", err)
	}

	departTime, err := timeTools.String2TimeYMDhms(req.DepartTime)
	if err != nil {
		return nil, errors.Wrapf(ERRInvalidInput, "invalid input departTime err: %v", err)
	}

	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("用户不存在"), "userID:%d 小于 1", userID)
	}

	bookResp, err := l.svcCtx.FlightReservationRpcClient.BookAirTickets(l.ctx, &flightreservation.BookAirTicketsReq{
		UserID:         userID,
		FlightNumber:   req.FlightNumber,
		SetOutDate:     timestamppb.New(sod),
		IsFirstClass:   req.IsFirstClass,
		DepartPosition: req.DepartPosition,
		DepartTime:     timestamppb.New(departTime),
		ArrivePosition: req.ArrivePosition,
		ArriveTime:     timestamppb.New(arriveTime),
	})
	if err != nil {
		//return nil, errors.Wrapf(ERRBookFail, "Rpc ERR in l.svcCtx.FlightReservationRpcClient.BookAirTickets: err: %v, req: %+v", err, req)
		return nil, err
	}

	resp = &types.FlightReservationResp{}
	resp.OrderSn = bookResp.OrderSn
	return
}
