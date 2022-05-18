package payment

import (
	"context"
	"fmt"
	"github.com/dtm-labs/dtm/dtmcli/logger"
	"github.com/dtm-labs/dtm/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"
	"uranus/app/payment/cmd/rpc/payment"
	"uranus/app/payment/model"
	"uranus/app/stock/cmd/rpc/stock"
	"uranus/app/userCenter/cmd/rpc/usercenter"
	"uranus/common/ctxdata"
	"uranus/common/uniqueid"
	"uranus/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type LocalPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

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

	if !strings.HasPrefix(req.OrderSn, fmt.Sprintf("%s", uniqueid.SN_PREFIX_FLIGHT_ORDER)) {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法订单号"), "invalid orderSn: %s", req.OrderSn)
	}

	// 获取订单详情
	orderDetail, err := l.svcCtx.OrderClient.FlightOrderDetail(l.ctx, &order.FlightOrderDetailReq{Sn: req.OrderSn})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("无该订单号对应的订单详情"), "无该订单号对应的订单详情 err: %v, orderSn: %s", err, req.OrderSn)
	}

	// 执行本地支付
	resp, err = l.execLocalPay(orderDetail)
	//if err != nil {
	//
	//	// 支付失败，恢复库存
	//	if e := l.recoverSurplus(orderDetail.FlightOrder.TicketId); e != nil {
	//		return nil, errors.Wrapf(err, "恢复库存失败, err: %v", e)
	//	}
	//	return nil, errors.Wrapf(err, "支付失败, err: %v", err)
	//}
	if err != nil {
		return nil, err
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

	var paymentSn string
	var userMoney int64
	err = mr.Finish(func() error {
		// 创建支付流水
		paymentResp, err := l.svcCtx.PaymentClient.CreatePayment(l.ctx, &payment.CreatePaymentReq{
			UserID:   userID,
			PayMode:  model.PayModeWalletBalance,
			PayTotal: orderDetail.FlightOrder.OrderTotalPrice,
			OrderSn:  orderDetail.FlightOrder.Sn,
		})
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("创建支付流水失败"), "创建支付流水失败 paymentSn: %s, userId: %d, orderSn: %s, payTotal: %d, err: %v",
				paymentResp.Sn, userID, orderDetail.FlightOrder.Sn, orderDetail.FlightOrder.OrderTotalPrice, err)
		}

		paymentSn = paymentResp.Sn
		return nil
	}, func() error {
		// 获取用户钱包余额
		getMoneyResp, err := l.svcCtx.UserCenterClient.GetUserMoney(l.ctx, &usercenter.GetUserMoneyReq{UserId: userID})
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("获取用户钱包余额失败"), "获取用户钱包余额失败 err: %v, userID: %d", err, userID)
		}

		userMoney = getMoneyResp.Money
		return nil
	})
	if err != nil {
		return nil, err
	}

	if userMoney-orderDetail.FlightOrder.OrderTotalPrice < 0 {
		// 余额不足，更新支付状态为支付失败
		updatePaymentReq := &payment.UpdateTradeStateReq{
			Sn:        paymentSn,
			PayStatus: model.PaymentLocalPayStatusFAIL,
			PayTime:   timestamppb.Now(),
			PayMode:   model.PayModeWalletBalance,
		}
		gid := dtmgrpc.MustGenGid(l.svcCtx.Config.DtmServer.Target)
		saga := dtmgrpc.NewSagaGrpc(l.svcCtx.Config.DtmServer.Target, gid).
			Add(l.svcCtx.Config.PaymentRpcConf.Target+"/pb.payment/UpdateTradeState", l.svcCtx.Config.PaymentRpcConf.Target+"/pb.payment/UpdateTradeStateRollBack", updatePaymentReq)
		err = saga.Submit()
		logger.FatalIfError(err)
		if err != nil {
			return nil, fmt.Errorf("submit data to  dtm-server err  : %+v \n", err)
		}
		return nil, errors.Wrapf(xerr.NewErrMsg("用户钱包余额不足"), "余额不足, orderTotalPrice: %d, userId: %d, user's money: %d", orderDetail.FlightOrder.OrderTotalPrice, userID, userMoney)
	}

	// 用分布式事务处理
	// 更新支付状态
	updatePaymentReq := &payment.UpdateTradeStateReq{
		Sn:        paymentSn,
		PayStatus: model.PaymentLocalPayStatusSuccess,
		PayTime:   timestamppb.Now(),
		PayMode:   model.PayModeWalletBalance,
	}
	// 扣除用户钱包余额
	updateUserWalletReq := &usercenter.DeductMoneyReq{
		UserId: userID,
		Money:  orderDetail.FlightOrder.OrderTotalPrice,
	}
	// 扣除库存
	deductReq := &stock.DeductStockByTicketIDReq{
		TicketID: orderDetail.FlightOrder.TicketId,
		Num:      1,
	}
	// 释放锁定的库存
	releaseReq := &stock.ReleaseStockByTicketIDReq{
		TicketID: orderDetail.FlightOrder.TicketId,
		Num:      1,
	}

	// 从配置文件生产服务地址
	stockServer, err := l.svcCtx.Config.StockRpcConf.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取库存服务地址失败"), "err: %v", err)
	}
	usercenterServer, err := l.svcCtx.Config.UserCenterRpcConf.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取用户中心服务地址失败"), "err: %v", err)
	}
	paymentServer, err := l.svcCtx.Config.PaymentRpcConf.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取支付服务地址失败"), "err: %v", err)
	}

	// 执行事务
	dtmServer := l.svcCtx.Config.DtmServer.Target
	gid := dtmgrpc.MustGenGid(dtmServer)
	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(usercenterServer+"/pb.usercenter/DeductMoney", usercenterServer+"/pb.usercenter/DeductMontyRollBack", updateUserWalletReq).
		Add(paymentServer+"/pb.payment/UpdateTradeState", paymentServer+"/pb.payment/UpdateTradeStateRollBack", updatePaymentReq).
		Add(stockServer+"/pb.stock/DeductStockByTicketID", stockServer+"/pb.stock/DeductStockByTicketIDRollBack", deductReq).
		Add(stockServer+"/pb.stock/ReleaseStockByTicketID", stockServer+"/pb.stock/ReleaseStockByTicketIDRollBack", releaseReq).
		EnableConcurrent().
		AddBranchOrder(3, []int{2})
	err = saga.Submit()
	logger.FatalIfError(err)
	if err != nil {
		return nil, fmt.Errorf("submit data to  dtm-server err  : %+v \n", err)
	}
	return
}

//// 恢复占用的库存
//func (l *LocalPayLogic) recoverSurplus(ticketID int64) error {
//	// 恢复该订单占用的库存
//	ticket, err := l.svcCtx.TicketsModel.FindOne(ticketID)
//	if err != nil && err != commonModel.ErrNotFound {
//		return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
//	}
//	if ticket == nil {
//		return errors.Wrapf(xerr.NewErrMsg("票不存在"), "票不存在, ticketID: %d", ticketID)
//	}
//	space, err := l.svcCtx.SpacesModel.FindOne(ticket.SpaceId)
//	if err != nil && err != commonModel.ErrNotFound {
//		return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
//	}
//	if ticket == nil {
//		return errors.Wrapf(xerr.NewErrMsg("舱位不存在"), "舱位不存在, spaceID: %d", ticket.SpaceId)
//	}
//	space.Surplus = space.Surplus + 1
//	err = l.svcCtx.SpacesModel.UpdateWithVersion(nil, space)
//	if err != nil {
//		return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "更新库存失败 space: %+v", space)
//	}
//	return nil
//}
