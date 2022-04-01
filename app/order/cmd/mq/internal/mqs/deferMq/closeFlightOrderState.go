package deferMq

import (
	"context"
	"encoding/json"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/order/model"
	"uranus/common/asynqmq"
	"uranus/common/xerr"
	"uranus/commonModel"

	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
)

func (l *AsynqTask) closeFlightOrderStateMqHandler(ctx context.Context, t *asynq.Task) error {

	var p asynqmq.FlightOrderCloseTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("解析asynq task payload err"), "closeFlightOrderStateMqHandler payload err:%v, payLoad:%+v", err, t.Payload())
	}

	resp, err := l.svcCtx.OrderClient.FlightOrderDetail(ctx, &order.FlightOrderDetailReq{
		Sn: p.Sn,
	})
	if err != nil || resp.FlightOrder == nil {
		return errors.Wrapf(xerr.NewErrMsg("获取订单失败"), "closeFlightOrderStateMqHandler 获取订单失败 or 订单不存在 err:%v, sn:%s ,FlightOrder : %+v", err, p.Sn, resp.FlightOrder)
	}

	if resp.FlightOrder.TradeState == model.FlightOrderTradeStateWaitPay {
		_, err := l.svcCtx.OrderClient.UpdateFlightOrderTradeState(ctx, &order.UpdateFlightOrderTradeStateReq{
			Sn:         p.Sn,
			TradeState: model.FlightOrderTradeStateCancel,
		})
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("关闭订单失败"), "closeFlightOrderStateMqHandler 关闭订单失败  err:%v, sn:%s ", err, p.Sn)
		}
	}
	// 恢复该订单占用的库存
	ticket, err := l.svcCtx.TicketsModel.FindOne(resp.FlightOrder.TicketId)
	if err != nil && err != commonModel.ErrNotFound {
		return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}
	if ticket == nil {
		return errors.Wrapf(xerr.NewErrMsg("票不存在"), "票不存在, ticketID: %d", resp.FlightOrder.TicketId)
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
