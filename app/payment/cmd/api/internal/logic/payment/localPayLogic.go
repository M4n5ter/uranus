package payment

import (
	"context"
	"github.com/pkg/errors"
	"uranus/common/xerr"

	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"

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

	return
}
