package payment

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/payment/cmd/rpc/payment"
	"uranus/app/payment/model"
	"uranus/app/userCenter/cmd/rpc/userCenter"
	"uranus/common/ctxdata"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LocalPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

var ERRFailToPay = xerr.NewErrCodeMsg(400, "支付失败")
var ERRRpcErr = xerr.NewErrCodeMsg(500, "支付失败")

func NewLocalPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) LocalPayLogic {
	return LocalPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LocalPayLogic) LocalPay(req *types.LocalPaymentReq) (resp *types.LocalPaymentResp, err error) {
	// 检查输入是否合法
	if len(req.OrderSn) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("订单号不能为空"), "订单号为空")
	}

	// 获取订单详情
	orderDetail, err := l.svcCtx.OrderClient.FlightOrderDetail(l.ctx, &order.FlightOrderDetailReq{Sn: req.OrderSn})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("无该订单号对应的订单详情"), "无该订单号对应的订单详情 err: %v, orderSn: %s", err, req.OrderSn)
	}

	// 执行本地支付
	resp, err = l.execLocalPay(orderDetail)
	if err != nil {
		// 支付失败，恢复库存
		if e := l.recoverSurplus(orderDetail.FlightOrder.TicketId); e != nil {
			return nil, errors.Wrapf(ERRRpcErr, "恢复库存失败, err: %v", e)
		}
		return nil, errors.Wrapf(ERRFailToPay, "支付失败, err: %v", err)
	}

	return
}

// 本地支付
func (l *LocalPayLogic) execLocalPay(orderDetail *order.FlightOrderDetailResp) (resp *types.LocalPaymentResp, err error) {
	resp = &types.LocalPaymentResp{}

	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID != orderDetail.FlightOrder.UserId {
		return nil, errors.Wrapf(xerr.NewErrMsg("该订单号不属于该用户"), "该订单号不属于该用户 userID: %d, order's userID: %d", userID, orderDetail.FlightOrder.UserId)
	}

	// 创建支付流水
	paymentSn, err := l.svcCtx.PaymentClient.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		UserID:   userID,
		PayMode:  model.PayModeWalletBalance,
		PayTotal: orderDetail.FlightOrder.OrderTotalPrice,
		OrderSn:  orderDetail.FlightOrder.Sn,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("创建支付流水失败"), "创建支付流水失败 paymentSn: %s, userId: %d, orderSn: %s, payTotal: %d, err: %v",
			paymentSn, userID, orderDetail.FlightOrder.Sn, orderDetail.FlightOrder.OrderTotalPrice, err)
	}

	// 获取用户钱包余额
	getMoneyResp, err := l.svcCtx.UserCenterClient.GetUserMoney(l.ctx, &userCenter.GetUserMoneyReq{UserId: userID})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取用户钱包余额失败"), "获取用户钱包余额失败 err: %v, userID: %d", err, userID)
	}
	if getMoneyResp.Money-orderDetail.FlightOrder.OrderTotalPrice < 0 {
		// 余额不足，更新支付状态为支付失败
		_, err = l.svcCtx.PaymentClient.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
			Sn:        paymentSn.Sn,
			PayStatus: model.PaymentLocalPayStatusFAIL,
			PayTime:   timestamppb.Now(),
			PayMode:   model.PayModeWalletBalance,
		})
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("用户钱包余额不足,应更新支付状态为失败，但操作失败"), "用户钱包余额不足,应更新支付状态为失败，但操作失败 err: %v", err)
		}
		return nil, errors.Wrapf(xerr.NewErrMsg("用户钱包余额不足"), "余额不足, orderTotalPrice: %d, userId: %d", orderDetail.FlightOrder.OrderTotalPrice, userID)
	}

	// 更新支付状态和扣除钱包余额应使用分布式事务 todo
	//gid := dtmgrpc.MustGenGid(l.svcCtx.Config.DtmServer.Target)
	//saga := dtmgrpc.NewSagaGrpc(l.svcCtx.Config.DtmServer.Target, gid).
	//	Add().
	//	Add()

	// 更新支付状态
	_, err = l.svcCtx.PaymentClient.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
		Sn:        paymentSn.Sn,
		PayStatus: model.PaymentLocalPayStatusSuccess,
		PayTime:   timestamppb.Now(),
		PayMode:   model.PayModeWalletBalance,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("更新支付状态失败"), "err %v", err)
	}

	// 扣除用户钱包余额
	_, err = l.svcCtx.UserCenterClient.UpdateUserWallet(l.ctx, &userCenter.UpdateUserWalletReq{
		UserId: userID,
		Money:  getMoneyResp.Money - orderDetail.FlightOrder.OrderTotalPrice,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("从用户钱包余额扣除金额失败"), "从用户钱包余额扣除金额失败 userId: %d, err: %v", userID, err)
	}

	return
}

// 恢复占用的库存
func (l *LocalPayLogic) recoverSurplus(ticketID int64) error {
	// 恢复该订单占用的库存
	ticket, err := l.svcCtx.TicketsModel.FindOne(ticketID)
	if err != nil && err != commonModel.ErrNotFound {
		return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}
	if ticket == nil {
		return errors.Wrapf(xerr.NewErrMsg("票不存在"), "票不存在, ticketID: %d", ticketID)
	}
	space, err := l.svcCtx.SpacesModel.FindOne(ticket.SpaceId)
	if err != nil && err != commonModel.ErrNotFound {
		return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}
	if ticket == nil {
		return errors.Wrapf(xerr.NewErrMsg("舱位不存在"), "舱位不存在, spaceID: %d", ticket.SpaceId)
	}
	space.Surplus = space.Surplus + 1
	err = l.svcCtx.SpacesModel.UpdateWithVersion(nil, space)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "更新库存失败 space: %+v", space)
	}
	return nil
}
