package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/order/cmd/rpc/order"
	orderModel "uranus/app/order/model"
	"uranus/app/payment/cmd/rpc/payment"
	paymentModel "uranus/app/payment/model"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightReservation/cmd/rpc/internal/svc"
	"uranus/app/flightReservation/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAirTicketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeAirTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAirTicketsLogic {
	return &ChangeAirTicketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChangeAirTickets 给定：用户的平台唯一id 用户拥有的对应订单号 目标舱位id
func (l *ChangeAirTicketsLogic) ChangeAirTickets(in *pb.ChangeAirTicketsReq) (*pb.ChangeAirTicketsResp, error) {
	// 查询用户支付成功的订单列表
	orderListResp, err := l.svcCtx.OrderRpcClient.UserFlightOrderList(l.ctx, &order.UserFlightOrderListReq{
		LastId:      0,
		PageSize:    1024,
		UserId:      in.UserID,
		TraderState: 1,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "查询用户支付成功的订单列表失败, userID: %d, err: %v", in.UserID, err)
	}
	// 检验该订单是否在支付成功的订单列表中
	oriTicketID, ok := l.svcCtx.ValidateOrderSn(in.OrderSn, orderListResp.List)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("只有支付成功的订单才能退款"), "orderSn: %s, userID: %d", in.OrderSn, in.UserID)
	}

	// 1. 根据票 ID 获取 spaceID, 然后获取相应舱位
	oriTicket, err := l.svcCtx.TicketsModel.FindOne(oriTicketID)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: ticketID: %d, err: %v", oriTicketID, err)
	}

	if oriTicket == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("未找到对应票信息"), "Err Not Found: ticketID:%d", oriTicketID)
	}

	oriSpace, err := l.svcCtx.SpacesModel.FindOne(oriTicket.SpaceId)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: spaceID: %d, err: %v", oriTicket.SpaceId, err)
	}

	if oriSpace == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("舱位不存在"), "查不到原舱位,spaceID: %d", oriTicket.SpaceId)
	}

	// 2. 查询目标舱位库存和价格
	newSpace, err := l.svcCtx.SpacesModel.FindOne(in.SpaceID)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: spaceID: %d, err: %v", in.SpaceID, err)
	}

	if newSpace == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("舱位不存在"), "查不到新舱位,spaceID: %d", in.SpaceID)
	}

	available := newSpace.Surplus - newSpace.LockedStock
	if available < 1 {
		// 库存不足
		return nil, errors.Wrapf(xerr.NewErrMsg("新舱位库存不足"), "新仓位库存为: %d", available)
	}

	// 查询目标舱位票价
	newTicket, err := l.svcCtx.TicketsModel.FindOneBySpaceId(in.SpaceID)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}
	if newTicket == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("查询不到目标舱位的票信息"), "Err Not Found, Find Ticket By SpaceId: %d", in.SpaceID)
	}

	newPrice := newTicket.Price - int64(float64(newTicket.Discount)/100*float64(newTicket.Price))

	// 3. 创建订单
	// 查询旧订单支付时的费用
	paymentResp, err := l.svcCtx.PaymentRpcClient.GetPaymentBySn(l.ctx, &payment.GetPaymentBySnReq{Sn: in.OrderSn})
	if err != nil {
		return nil, err
	}

	if paymentResp.PaymentDetail.PayStatus != paymentModel.PaymentLocalPayStatusSuccess {
		return nil, errors.Wrapf(xerr.NewErrMsg("支付流水异常，请联系管理员"), "支付流水状态与订单状态不符，订单为带使用，流水却不是已支付. paymentSn: %s", paymentResp.PaymentDetail.Sn)
	}

	oriPay := paymentResp.PaymentDetail.PayTotal

	// 查询改票手续费
	rci, err := l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(oriTicketID, 0)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: ticketID: %d, isRefund:%d, err: %v", oriTicketID, 0, err)
	}

	if rci == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("未找到退改票信息"), "Err Not Found: ticketID: %d, isRefund:%d", oriTicketID, 0)
	}

	fee, ok := l.svcCtx.GetFee(rci)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取改票手续费失败"), "Err when getting fee!")
	}

	// 废弃旧订单
	err = l.discardOldOrder(in.OrderSn)
	if err != nil {
		return nil, err
	}

	// 差价
	diff := oriPay - newPrice - fee
	if diff > 0 {
		// 需要退用户钱

	} else if diff < 0 {
		// 需要用户补差价
	} else {
		// 没有差价
	}
	return &pb.ChangeAirTicketsResp{}, nil
}

// 废弃订单
func (l *ChangeAirTicketsLogic) discardOldOrder(orderSn string) error {
	_, err := l.svcCtx.OrderRpcClient.UpdateFlightOrderTradeState(l.ctx, &order.UpdateFlightOrderTradeStateReq{
		Sn:         orderSn,
		TradeState: orderModel.FlightOrderTradeStateDiscard,
	})
	if err != nil {
		return err
	}

	return nil
}
