package logic

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/order/model"
	"uranus/common/xerr"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRDBERR = xerr.NewErrCode(xerr.DB_ERROR)

type UpdateFlightOrderTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFlightOrderTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFlightOrderTradeStateLogic {
	return &UpdateFlightOrderTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateFlightOrderTradeState 更新机票订单状态
func (l *UpdateFlightOrderTradeStateLogic) UpdateFlightOrderTradeState(in *pb.UpdateFlightOrderTradeStateReq) (*pb.UpdateFlightOrderTradeStateResp, error) {
	// 查询当前订单信息
	order, err := l.svcCtx.OrderModel.FindOneBySn(in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(ERRDBERR, "DBERR: err: %v, Sn: %s", err, in.Sn)
	}
	if order == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("订单不存在"), "Not Found order : Sn: %s", in.Sn)
	}

	// 如果已经一致，则不需要更新，直接返回
	if order.TradeState == in.TradeState {
		return &pb.UpdateFlightOrderTradeStateResp{}, nil
	}

	// 校验订单状态
	if err := l.verifyOrderTradeState(in.TradeState, order.TradeState); err != nil {
		return nil, errors.WithMessagef(err, " , in : %+v", in)
	}

	// 更新订单状态
	order.TradeState = in.TradeState
	if err := l.svcCtx.OrderModel.UpdateWithVersion(nil, order); err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("更新机票订单状态失败"), "更新机票订单状态失败 err:%v , in : %+v", err, in)

	}

	return &pb.UpdateFlightOrderTradeStateResp{
		Id:              order.Id,
		UserId:          order.UserId,
		Sn:              order.Sn,
		TradeCode:       order.TradeCode,
		OrderTotalPrice: order.OrderTotalPrice,
		DepartPosition:  order.DepartPosition,
		DepartTime:      timestamppb.New(order.DepartTime),
		ArrivePosition:  order.ArrivePosition,
		ArriveTime:      timestamppb.New(order.ArriveTime),
	}, nil
}

// 校验订单状态，校验通过则返回 nil，否则返回 error
func (l *UpdateFlightOrderTradeStateLogic) verifyOrderTradeState(newTradeState, oldTradeState int64) error {

	switch newTradeState {
	case model.FlightOrderTradeStateWaitPay:
		return errors.Wrapf(xerr.NewErrMsg("不支持更改此状态"),
			"不支持更改为待支付状态 newTradeState: %d, oldTradeState: %d", newTradeState, oldTradeState)
	case model.FlightOrderTradeStateCancel:
		if oldTradeState != model.FlightOrderTradeStateWaitPay {
			return errors.Wrapf(xerr.NewErrMsg("只有待支付的订单才能被取消"),
				"只有待支付的订单才能被取消 newTradeState: %d, oldTradeState: %d", newTradeState, oldTradeState)
		}
	case model.FlightOrderTradeStateWaitUse:
		if oldTradeState != model.FlightOrderTradeStateWaitPay {
			return errors.Wrapf(xerr.NewErrMsg("只有待支付的订单才能更改为此状态"),
				"只有待支付的订单才能更改为未使用状态 newTradeState: %d, oldTradeState: %d", newTradeState, oldTradeState)
		}
	case model.FlightOrderTradeStateUsed:
		if oldTradeState != model.FlightOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("只有未使用的订单才能更改为此状态"),
				"只有未使用的订单才能更改为已使用状态 newTradeState: %d, oldTradeState: %d", newTradeState, oldTradeState)
		}
	case model.FlightOrderTradeStateRefund:
		if oldTradeState != model.FlightOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("只有未使用的订单才能更改为此状态"),
				"只有未使用的订单才能更改为已退款状态 newTradeState: %d, oldTradeState: %d", newTradeState, oldTradeState)
		}
	case model.FlightOrderTradeStateExpire:
		if oldTradeState != model.FlightOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("只有未使用的订单才能更改为此状态"),
				"只有未使用的订单才能更改为已过期状态 newTradeState: %d, oldTradeState: %d", newTradeState, oldTradeState)
		}
	}

	return nil
}
