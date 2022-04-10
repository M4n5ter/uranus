package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/order/model"
	"uranus/common/tool"
	"uranus/common/uniqueid"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFlightOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFlightOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFlightOrderLogic {
	return &CreateFlightOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateFlightOrder 机票下订单
func (l *CreateFlightOrderLogic) CreateFlightOrder(in *pb.CreateFlightOrderReq) (*pb.CreateFlightOrderResp, error) {
	if in.UserId < 1 || in.TicketId < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input: userID: %d, ticketID: %d", in.UserId, in.TicketId)
	}
	// 获取机票详情
	flightDetailResp, err := l.svcCtx.FlightInquiryClient.FlightDetail(l.ctx, &flightinquiry.FlightDetailReq{TicketID: in.TicketId})
	if err != nil {
		return nil, err
	}
	if flightDetailResp == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("查询不到票所对应的航班信息"), "查询不到票所对应的航班信息, ticketID: %d", in.TicketId)
	}
	// 创建订单 + 锁定库存(本地事务)
	order := new(model.FlightOrder)
	err = l.svcCtx.OrderModel.Trans(func(session sqlx.Session) error {
		// 创建订单
		order.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_FLIGHT_ORDER)
		order.OrderTotalPrice = int64(flightDetailResp.FlightInfo.Price) - int64(float64(flightDetailResp.FlightInfo.Discount)/100*float64(flightDetailResp.FlightInfo.Price))
		order.DepartPosition = flightDetailResp.FlightInfo.DepartPosition
		order.DepartTime = flightDetailResp.FlightInfo.DepartTime.AsTime()
		order.ArrivePosition = flightDetailResp.FlightInfo.ArrivePosition
		order.ArriveTime = flightDetailResp.FlightInfo.ArriveTime.AsTime()
		order.Discount = flightDetailResp.FlightInfo.Discount
		order.TicketId = in.TicketId
		order.UserId = in.UserId
		order.TradeCode = tool.Krand(8, tool.KC_RAND_KIND_ALL)
		order.TicketPrice = int64(flightDetailResp.FlightInfo.Price)
		order.TradeState = model.FlightOrderTradeStateWaitPay
		_, err = l.svcCtx.OrderModel.Insert(session, order)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "下单数据库异常 order : %+v , err: %v", order, err)
		}

		// 锁定库存
		ticket, err := l.svcCtx.TicketsModel.FindOne(in.TicketId)
		if err != nil && err != commonModel.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "查找票失败, ticketID: %d, err: %v", in.TicketId, err)
		}
		if ticket == nil {
			return errors.Wrapf(xerr.NewErrMsg("找不到票信息"), "Not Found ticketID: %d", in.TicketId)
		}

		space, err := l.svcCtx.SpacesModel.FindOne(ticket.SpaceId)
		if err != nil && err != commonModel.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "查找舱位失败, spaceID: %d, err: %v", ticket.SpaceId, err)
		}
		if space == nil {
			return errors.Wrapf(xerr.NewErrMsg("找不到舱位信息"), "Not Found spaceID: %d", ticket.SpaceId)
		}
		if space.Surplus-space.LockedStock-1 < 0 {
			return errors.Wrapf(xerr.NewErrMsg("库存不足"), "stock is not enough,spaceID: %d", space.Id)
		}

		space.LockedStock = space.LockedStock + 1
		err = l.svcCtx.SpacesModel.UpdateWithVersion(session, space)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 延迟关闭订单任务
	_, _ = l.svcCtx.MqueueClient.AqDeferFlightOrderClose(l.ctx, &mqueue.AqDeferFlightOrderCloseReq{Sn: order.Sn})
	return &pb.CreateFlightOrderResp{Sn: order.Sn}, nil
}
