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

var ERRDBERR = xerr.NewErrCode(xerr.DB_ERROR)

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

// CreatePayment 创建支付预处理订单
func (l *CreatePaymentLogic) CreatePayment(in *pb.CreatePaymentReq) (*pb.CreatePaymentResp, error) {
	// 检查输入合法性
	if in.UserID < 1 || len(in.OrderSn) == 0 || len(in.PayModel) == 0 || in.PayTotal < 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input %+v", in)
	}

	// 创建支付
	data := &model.Payment{}
	data.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_THIRD_PAYMENT)
	data.OrderSn = in.OrderSn
	data.UserId = in.UserID
	data.PayMode = in.PayModel
	data.PayTotal = in.PayTotal

	_, err := l.svcCtx.PaymentModel.Insert(nil, data)
	if err != nil {
		return nil, errors.Wrapf(ERRDBERR, "DBERR: 创建支付流水记录失败: err: %v, data: %+v", err, data)
	}

	return &pb.CreatePaymentResp{Sn: data.Sn}, nil
}
