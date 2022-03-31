package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/payment/model"
	"uranus/common/xerr"

	"uranus/app/payment/cmd/rpc/internal/svc"
	"uranus/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentSuccessRefundByOrderSnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentSuccessRefundByOrderSnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentSuccessRefundByOrderSnLogic {
	return &GetPaymentSuccessRefundByOrderSnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetPaymentSuccessRefundByOrderSn 根据订单sn查询流水记录
func (l *GetPaymentSuccessRefundByOrderSnLogic) GetPaymentSuccessRefundByOrderSn(in *pb.GetPaymentSuccessRefundByOrderSnReq) (*pb.GetPaymentSuccessRefundByOrderSnResp, error) {
	// 检查输入合法性
	if len(in.OrderSn) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input : orderSn: %s", in.OrderSn)
	}

	whereBuilder := l.svcCtx.PaymentModel.RowBuilder().Where(
		"order_sn = ? AND (trade_state = ? OR trade_state = ?)",
		in.OrderSn, model.PaymentPayTradeStateSuccess, model.PaymentPayTradeStateRefund)

	payment, err := l.svcCtx.PaymentModel.FindOneByQuery(whereBuilder)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(ERRDBERR, "DBERR: 根据订单sn查询支付流水失败, err: %v, in: %+v", err, in)
	}

	var resp pb.PaymentDetail
	if payment != nil {
		_ = copier.Copy(&resp, payment)
		resp.CreateTime = timestamppb.New(payment.CreateTime)
		resp.UpdateTime = timestamppb.New(payment.UpdateTime)
		resp.PayTime = timestamppb.New(payment.PayTime)
	}

	return &pb.GetPaymentSuccessRefundByOrderSnResp{PaymentDetail: &resp}, nil
}
