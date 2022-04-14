package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	userCenterModel "uranus/app/userCenter/model"

	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/order/cmd/mq/internal/svc"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/order/model"
	paymentModel "uranus/app/payment/model"
	"uranus/app/usercenter/cmd/rpc/usercenter"
	"uranus/common/kqueue"
	"uranus/common/tool"
	"uranus/common/wxminisub"
	"uranus/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

/**
监听支付流水状态变更通知消息队列
*/
type PaymentUpdateStatusMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentUpdateStatusMq(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentUpdateStatusMq {
	return &PaymentUpdateStatusMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentUpdateStatusMq) Consume(_, val string) error {

	fmt.Printf(" PaymentUpdateStatusMq Consume val : %s \n", val)
	//解析数据
	var message kqueue.PaymentUpdatePayStatusNotifyMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Error("PaymentUpdateStatusMq->Consume Unmarshal err : %v , val : %s", err, val)
		return err
	}

	//执行业务..
	if err := l.execService(message); err != nil {
		logx.WithContext(l.ctx).Error("PaymentUpdateStatusMq->execService  err : %v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

//执行业务
func (l *PaymentUpdateStatusMq) execService(message kqueue.PaymentUpdatePayStatusNotifyMessage) error {

	orderTradeState := l.getOrderTradeStateByPaymentTradeState(message.PayStatus)
	if orderTradeState != -99 {
		//更新订单状态
		resp, err := l.svcCtx.OrderClient.UpdateFlightOrderTradeState(l.ctx, &order.UpdateFlightOrderTradeStateReq{
			Sn:         message.OrderSn,
			TradeState: orderTradeState,
		})
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("更新订单状态失败"), " err : %v ,message:%+v", err, message)
		}

		// 判断是否是微信用户，微信用户要推送通知，平台用户不需要
		orderDetail, err := l.svcCtx.OrderClient.FlightOrderDetail(l.ctx, &order.FlightOrderDetailReq{Sn: message.OrderSn})
		if err != nil && err != model.ErrNotFound {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
		}
		if orderDetail != nil {
			_, err = l.svcCtx.UserCenterClient.GetUserAuthByUserId(l.ctx, &usercenter.GetUserAuthByUserIdReq{
				UserId:   orderDetail.FlightOrder.UserId,
				AuthType: "system",
			})
			if err != nil && err == model.ErrNotFound {
				//发送短信、微信等通知用户
				l.notifyUser(resp.Sn, resp.TradeCode, resp.DepartPosition, resp.ArrivePosition, resp.OrderTotalPrice, resp.UserId, resp.DepartTime.AsTime(), resp.ArriveTime.AsTime())
			}
		}

	}

	return nil
}

//根据支付状态获取订单状态.
func (l *PaymentUpdateStatusMq) getOrderTradeStateByPaymentTradeState(paymentPayStatus int64) int64 {

	switch paymentPayStatus {
	case paymentModel.CommonPaySuccess:
		return model.FlightOrderTradeStateWaitUse
	case paymentModel.CommonPayRefund:
		return model.FlightOrderTradeStateRefund
	case paymentModel.CommonPayDiscard:
		return model.FlightOrderTradeStateDiscard
	default:
		return -99
	}

}

//发送小程序模版消息通知用户
func (l *PaymentUpdateStatusMq) notifyUser(sn, code, departPosition, arrivePosition string, orderTotalPrice, userId int64, departTime, arriveTime time.Time) {

	userCenterResp, err := l.svcCtx.UserCenterClient.GetUserAuthByUserId(l.ctx, &usercenter.GetUserAuthByUserIdReq{
		UserId:   userId,
		AuthType: userCenterModel.UserAuthTypeSmallWX,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("PaymentUpdateStatusMq -> notifyUser err:%v , sn :%s , code : %s , userId : %d ", err, sn, code, userId)
	}
	if userCenterResp.UserAuth == nil || len(userCenterResp.UserAuth.AuthKey) == 0 {
		logx.WithContext(l.ctx).Errorf("PaymentUpdateStatusMq -> notifyUser 未保存用户openid信息， 数据异常  sn :%s , code : %s , userId : %d ", sn, code, userId)
	}

	openId := userCenterResp.UserAuth.AuthKey

	//发送小程序订单支付成功订阅消息..
	orderPaySuccessDataParam := wxminisub.OrderPaySuccessDataParam{
		Sn:             sn,
		PayTotal:       fmt.Sprintf("%.2f", tool.Fen2Yuan(orderTotalPrice)),
		DepartPosition: departPosition,
		DepartTime:     departTime.Local().Format("2006-01-02 15:04:05"),
		ArrivePosition: arrivePosition,
		ArriveTime:     arriveTime.Local().Format("2006-01-02 15:04:05"),
	}
	if _, err = l.svcCtx.MqueueClient.SendWxMiniSubMessage(l.ctx, &mqueue.SendWxMiniSubMessageReq{
		TemplateID: wxminisub.OrderPaySuccessTemplateID,
		Openid:     openId,
		Data:       wxminisub.OrderPaySuccessData(orderPaySuccessDataParam),
	}); err != nil {
		logx.WithContext(l.ctx).Errorf("发送小程序订单支付成功订阅消息失败 orderPaySuccessDataParam : %+v \n", orderPaySuccessDataParam)
	}

}
