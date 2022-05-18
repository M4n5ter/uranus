package logic

import (
	"context"
	"github.com/dtm-labs/dtm/dtmcli/logger"
	"github.com/dtm-labs/dtm/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"uranus/app/order/cmd/rpc/order"
	orderModel "uranus/app/order/model"
	"uranus/app/payment/cmd/rpc/payment"
	paymentModel "uranus/app/payment/model"
	"uranus/app/userCenter/cmd/rpc/usercenter"
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

	var oriSpace *commonModel.Spaces
	//var newSpace *commonModel.Spaces
	var newTicketID int64
	var newPrice int64
	var oriPay int64
	var fee int64
	err = mr.Finish(func() error {
		// 1. 根据票 ID 获取 spaceID, 然后获取相应舱位
		oriTicket, err := l.svcCtx.TicketsModel.FindOne(oriTicketID)
		if err != nil && err != commonModel.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: ticketID: %d, err: %v", oriTicketID, err)
		}

		if oriTicket == nil {
			return errors.Wrapf(xerr.NewErrMsg("未找到对应票信息"), "Err Not Found: ticketID:%d", oriTicketID)
		}

		oriSpace, err = l.svcCtx.SpacesModel.FindOne(oriTicket.SpaceId)
		if err != nil && err != commonModel.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: spaceID: %d, err: %v", oriTicket.SpaceId, err)
		}

		if oriSpace == nil {
			return errors.Wrapf(xerr.NewErrMsg("舱位不存在"), "查不到原舱位,spaceID: %d", oriTicket.SpaceId)
		}
		return nil
	}, func() error {
		// 2. 查询目标舱位库存和价格
		newSpace, err := l.svcCtx.SpacesModel.FindOne(in.SpaceID)
		if err != nil && err != commonModel.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: spaceID: %d, err: %v", in.SpaceID, err)
		}

		if newSpace == nil {
			return errors.Wrapf(xerr.NewErrMsg("舱位不存在"), "查不到新舱位,spaceID: %d", in.SpaceID)
		}

		available := newSpace.Surplus - newSpace.LockedStock
		if available < 1 {
			// 库存不足
			return errors.Wrapf(xerr.NewErrMsg("新舱位库存不足"), "新仓位库存为: %d", available)
		}

		// 查询目标舱位票价
		newTicket, err := l.svcCtx.TicketsModel.FindOneBySpaceId(in.SpaceID)
		if err != nil && err != commonModel.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
		}
		if newTicket == nil {
			return errors.Wrapf(xerr.NewErrMsg("查询不到目标舱位的票信息"), "Err Not Found, Find Ticket By SpaceId: %d", in.SpaceID)
		}
		newTicketID = newTicket.Id
		newPrice = newTicket.Price - int64(float64(newTicket.Discount)/100*float64(newTicket.Price))
		return nil
	}, func() error {
		// 查询旧订单支付时的费用
		paymentResp, err := l.svcCtx.PaymentRpcClient.GetPaymentBySn(l.ctx, &payment.GetPaymentBySnReq{Sn: in.OrderSn})
		if err != nil {
			return err
		}

		if paymentResp.PaymentDetail.PayStatus != paymentModel.PaymentLocalPayStatusSuccess {
			return errors.Wrapf(xerr.NewErrMsg("支付流水异常，请联系管理员"), "支付流水状态与订单状态不符，订单为带使用，流水却不是已支付. paymentSn: %s", paymentResp.PaymentDetail.Sn)
		}

		oriPay = paymentResp.PaymentDetail.PayTotal
		return nil
	}, func() error {
		// 查询改票手续费
		rci, err := l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(oriTicketID, 0)
		if err != nil && err != commonModel.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: ticketID: %d, isRefund:%d, err: %v", oriTicketID, 0, err)
		}

		if rci == nil {
			return errors.Wrapf(xerr.NewErrMsg("未找到退改票信息"), "Err Not Found: ticketID: %d, isRefund:%d", oriTicketID, 0)
		}

		fee, ok = l.svcCtx.GetFee(rci)
		if !ok {
			return errors.Wrapf(xerr.NewErrMsg("获取改票手续费失败"), "Err when getting fee!")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// 3. 创建订单
	newOrderSn := ""
	// 差价
	diff := oriPay - newPrice - fee
	if diff > 0 {
		// 需要退用户钱

		newOrderSn, err = l.discardOldOrderAndCreateNew(in.OrderSn, oriSpace, newTicketID, in.UserID, 0)
		// 退钱
		updateUserWalletReq := &usercenter.AddMoneyReq{
			UserId: in.UserID,
			Money:  diff,
		}
		// 从配置文件生产服务地址
		usercenterServer, err := l.svcCtx.Config.UserCenterRpcConf.BuildTarget()
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("获取用户中心服务地址失败"), "err: %v", err)
		}
		// 执行事务
		dtmServer := l.svcCtx.Config.DtmServer.Target
		gid := dtmgrpc.MustGenGid(dtmServer)
		saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
			Add(usercenterServer+"/pb.usercenter/AddMoney", usercenterServer+"/pb.usercenter/AddMoneyRollback", updateUserWalletReq)
		err = saga.Submit()
		logger.FatalIfError(err)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("未能成功进行退款"), "submit data to  dtm-server err  : %+v \n", err)
		}

	} else if diff < 0 {
		// 需要用户补差价

		newOrderSn, err = l.discardOldOrderAndCreateNew(in.OrderSn, oriSpace, newTicketID, in.UserID, -diff)
	} else {
		// 没有差价

		newOrderSn, err = l.discardOldOrderAndCreateNew(in.OrderSn, oriSpace, newTicketID, in.UserID, 0)
	}

	return &pb.ChangeAirTicketsResp{OrderSn: newOrderSn}, nil
}

// 废弃旧订单，生产新订单
func (l *ChangeAirTicketsLogic) discardOldOrderAndCreateNew(orderSn string, oriSpace *commonModel.Spaces, newTicketID, userID, price int64) (newOrderSn string, err error) {
	_, err = l.svcCtx.OrderRpcClient.UpdateFlightOrderTradeState(l.ctx, &order.UpdateFlightOrderTradeStateReq{
		Sn:         orderSn,
		TradeState: orderModel.FlightOrderTradeStateDiscard,
	})
	if err != nil {
		return "", err
	}

	newOrderSn = ""
	// 释放旧订单占用的库存并且生产新订单
	err = l.svcCtx.SpacesModel.Trans(func(session sqlx.Session) error {
		// 恢复旧舱位库存
		oriSpace.Surplus = oriSpace.Surplus + 1
		err = l.svcCtx.SpacesModel.UpdateWithVersion(session, oriSpace)
		if err != nil {
			return err
		}

		// 创建新订单
		createOrderResp, err := l.svcCtx.OrderRpcClient.CreateFlightOrder(l.ctx, &order.CreateFlightOrderReq{
			TicketId: newTicketID,
			UserId:   userID,
		})
		if err != nil {
			return err
		}

		newOrderSn = createOrderResp.Sn

		// 改价
		_, err = l.svcCtx.OrderRpcClient.ChangeTheOrderPrice(l.ctx, &order.ChangeTheOrderPriceReq{OrderSn: newOrderSn, Price: price})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return newOrderSn, nil
}
