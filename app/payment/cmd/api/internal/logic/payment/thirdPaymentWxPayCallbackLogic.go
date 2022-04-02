package payment

import (
	"context"

	"uranus/app/payment/cmd/api/internal/svc"
	"uranus/app/payment/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdPaymentWxPayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThirdPaymentWxPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) ThirdPaymentWxPayCallbackLogic {
	return ThirdPaymentWxPayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentWxPayCallbackLogic) ThirdPaymentWxPayCallback(req *types.ThirdPaymentWxPayCallbackReq) (resp *types.ThirdPaymentWxPayCallbackResp, err error) {
	// todo: add your logic here and delete this line

	return
}
