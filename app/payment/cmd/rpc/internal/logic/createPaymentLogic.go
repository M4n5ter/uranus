package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/payment/model"
	"uranus/common/uniqueid"
	"uranus/common/xerr"

	"uranus/app/payment/cmd/rpc/internal/svc"
	"uranus/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreatePayment 创建微信支付预处理订单
func (l *CreatePaymentLogic) CreatePayment(in *pb.CreatePaymentReq) (*pb.CreatePaymentResp, error) {
	// 检查输入合法性
	if len(in.AuthKey) == 0 || len(in.OrderSn) == 0 || len(in.PayModel) == 0 || in.PayTotal < 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input %+v", in)
	}

	// 创建支付
	data := &model.Payment{}
	data.Sn			= uniqueid.GenSn(uniqueid.SN_PREFIX_THIRD_PAYMENT)
	data.OrderSn	= in.OrderSn
	data.AuthKey	= in.AuthKey
	data.PayMode	= in.PayModel
	data.PayTotal	= in.PayTotal

	_, err := l.svcCtx.

	return &pb.CreatePaymentResp{}, nil
}
