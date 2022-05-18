package logic

import (
	"context"
	"fmt"
	"github.com/dtm-labs/dtm/dtmcli/logger"
	"github.com/dtm-labs/dtm/dtmgrpc"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/payment/cmd/rpc/payment"
	"uranus/app/payment/model"
	"uranus/app/stock/cmd/rpc/stock"
	"uranus/app/userCenter/cmd/rpc/usercenter"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightReservation/cmd/rpc/internal/svc"
	"uranus/app/flightReservation/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRRefundFail = xerr.NewErrMsg("退订失败")

type RefundAirTicketsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundAirTicketsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundAirTicketsLogic {
	return &RefundAirTicketsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RefundAirTickets 给定：用户的平台唯一id 用户拥有的对应订单号 来退订机票
func (l *RefundAirTicketsLogic) RefundAirTickets(in *pb.RefundAirTicketsReq) (*pb.RefundAirTicketsResp, error) {

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
	ticketID, ok := l.svcCtx.ValidateOrderSn(in.OrderSn, orderListResp.List)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("只有支付成功的订单才能退款"), "orderSn: %s, userID: %d", in.OrderSn, in.UserID)
	}

	var fee int64
	var pdResp *payment.GetPaymentSuccessRefundByOrderSnResp
	err = mr.Finish(func() error {
		// 查询退票信息
		refundInfo, err := l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(ticketID, 1)
		if err != nil && err != commonModel.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
		}
		if refundInfo == nil {
			return errors.Wrapf(xerr.NewErrMsg("退票失败:该票没有对应的退票信息"), "")
		}

		// 获取退票手续费
		fee, ok = l.svcCtx.GetFee(refundInfo)
		if !ok {
			return errors.Wrapf(xerr.NewErrMsg("无法退款，已经超出最晚期限"), "")
		}
		return nil
	}, func() error {
		// 查询支付流水详情
		pdResp, err = l.svcCtx.PaymentRpcClient.GetPaymentSuccessRefundByOrderSn(l.ctx, &payment.GetPaymentSuccessRefundByOrderSnReq{
			OrderSn: in.OrderSn,
		})
		if err != nil {
			return errors.Wrapf(ERRRefundFail, "查询订单支付流水详情失败, err: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 从配置文件获取服务地址
	paymentServer, err := l.svcCtx.Config.PaymentRpcConf.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取支付服务地址失败"), "err: %v", err)
	}

	usercenterServer, err := l.svcCtx.Config.UserCenterRpcConf.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取用户中心服务地址失败"), "err: %v", err)
	}

	stockServer, err := l.svcCtx.Config.StockRpcConf.BuildTarget()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("获取库存服务地址失败"), "err: %v", err)
	}

	// 更新支付流水
	updateTradeStateReq := &payment.UpdateTradeStateReq{
		Sn:        pdResp.PaymentDetail.Sn,
		PayStatus: model.CommonPayRefund,
		PayTime:   timestamppb.Now(),
		PayMode:   model.PayModeWalletBalance,
	}
	// 退给用户钱
	addMoneyReq := &usercenter.AddMoneyReq{
		UserId: in.UserID,
		Money:  pdResp.PaymentDetail.PayTotal - fee,
	}
	// 更新库存
	addStockReq := &stock.AddStockByTicketIDReq{
		TicketID: ticketID,
		Num:      1,
	}

	// 执行事务
	dtmServer := l.svcCtx.Config.DtmServer.Target
	gid := dtmgrpc.MustGenGid(dtmServer)
	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(paymentServer+"/pb.payment/UpdateTradeState", paymentServer+"/pb.payment/UpdateTradeStateRollBack", updateTradeStateReq).
		Add(stockServer+"/pb.stock/AddStockByTicketID", stockServer+"/pb.stock/AddStockByTicketIDRollBack", addStockReq).
		Add(usercenterServer+"/pb.usercenter/AddMoney", usercenterServer+"/pb.usercenter/AddMoneyRollBack", addMoneyReq).
		EnableConcurrent()
	err = saga.Submit()
	logger.FatalIfError(err)
	if err != nil {
		return nil, fmt.Errorf("submit data to  dtm-server err  : %+v \n", err)
	}

	return &pb.RefundAirTicketsResp{}, nil
}
