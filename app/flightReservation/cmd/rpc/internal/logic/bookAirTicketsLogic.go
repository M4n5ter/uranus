package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/order/cmd/rpc/order"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightReservation/cmd/rpc/internal/svc"
	"uranus/app/flightReservation/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRNotFound = xerr.NewErrMsg("找不到对应机票信息")
var ERRDBERR = xerr.NewErrCode(xerr.DB_ERROR)
var ERRNotEnough = xerr.NewErrMsg("票库存不足")

type BookAirTicketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBookAirTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookAirTicketsLogic {
	return &BookAirTicketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BookAirTickets 给定：用户的平台唯一id 航班号 出发日期 是否为头等舱/商务舱 起飞地点/时间 降落地点/时间 来预定机票
func (l *BookAirTicketsLogic) BookAirTickets(in *pb.BookAirTicketsReq) (*pb.BookAirTicketsResp, error) {
	// 获取对应航班信息
	flightInfo, err := l.svcCtx.FlightInfosModel.FindOneByByNumberAndSetOutDateAndPosition(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightNumber, in.SetOutDate.AsTime(), in.DepartPosition, in.DepartTime.AsTime(), in.ArrivePosition, in.ArriveTime.AsTime())
	if err != nil {
		if err == commonModel.ErrNotFound {
			return nil, errors.Wrapf(ERRNotFound, "err not found:flightreservation-rpc.BookAirTickets.l.svcCtx.FlightInfosModel.FindOneByByNumberAndSetOutDateAndPosition,number: %s, sod: %v,departPosition: %s, departTime: %v, arrivePosition: %s, arriveTime: %v\n", in.FlightNumber, in.SetOutDate.AsTime(), in.DepartPosition, in.DepartTime.AsTime(), in.ArrivePosition, in.ArriveTime.AsTime())
		}
		return nil, errors.Wrapf(ERRDBERR, "DBERR in flightreservation-rpc.BookAirTickets.l.svcCtx.FlightInfosModel.FindOneByByNumberAndSetOutDateAndPosition:number: %s, sod: %v,departPosition: %s, departTime: %v, arrivePosition: %s, arriveTime: %v, err: %v\n", in.FlightNumber, in.SetOutDate.AsTime(), in.DepartPosition, in.DepartTime.AsTime(), in.ArrivePosition, in.ArriveTime.AsTime(), err)
	}

	var isFirstClass int64
	var resp string
	if in.IsFirstClass {
		isFirstClass = 1
	} else {
		isFirstClass = 0
	}

	// 查找对应的舱位信息
	space, err := l.svcCtx.SpacesModel.FindOneByFlightInfoIdIsFirstClass(flightInfo.Id, isFirstClass)
	if err != nil {
		if err == commonModel.ErrNotFound {
			return nil, errors.Wrapf(ERRNotFound, "err not found:flightreservation-rpc.BookAirTickets.l.svcCtx.SpacesModel.FindOneByFlightInfoIdIsFirstClass: flightInfoID:%d, isFirstClass:%d", flightInfo.Id, isFirstClass)
		}
		return nil, errors.Wrapf(ERRDBERR, "DBERR in flightreservation-rpc.BookAirTickets.l.svcCtx.SpacesModel.FindOneByFlightInfoIdIsFirstClass: flightInfoID:%d, isFirstClass:%d, err: %v", flightInfo.Id, isFirstClass, err)
	}
	if space.Surplus > 0 {
		// 有库存剩余时查找对应的票信息
		ticket, err := l.svcCtx.TicketsModel.FindOneBySpaceId(space.Id)
		if err != nil {
			if err == commonModel.ErrNotFound {
				return nil, errors.Wrapf(ERRNotFound, "err not found:flightreservation-rpc.BookAirTickets.l.svcCtx.TicketsModel.FindOneBySpaceId:spaceID:%d", space.Id)
			}
			return nil, errors.Wrapf(ERRDBERR, "DBERR in flightreservation-rpc.BookAirTickets.l.svcCtx.TicketsModel.FindOneBySpaceId: spaceID: %d, err: %v", space.Id, err)
		}

		// 创建订单
		orderResp, err := l.svcCtx.OrderRpcClient.CreateFlightOrder(l.ctx, &order.CreateFlightOrderReq{
			TicketId: ticket.Id,
			UserId:   in.UserID,
		})
		if err != nil {
			return nil, err
		}
		//更新库存移至创建订单内，用事务处理
		//更新库存,若超时未支付、支付失败或退款， 会在响应代码块将占用的库存释放
		//space.Surplus = space.Surplus - 1
		//err = l.svcCtx.SpacesModel.UpdateWithVersion(session, space)
		//if err != nil {
		//	return err
		//}
		resp = orderResp.Sn
	} else {
		// 库存不足
		return nil, errors.Wrapf(ERRNotEnough, "space.Surplus not enough,spaceid:%d,surplus:%d", space.Id, space.Surplus)
	}

	return &pb.BookAirTicketsResp{OrderSn: resp}, err
}
