package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/payment/model"
	"uranus/common/xerr"

	"uranus/app/payment/cmd/rpc/internal/svc"
	"uranus/app/payment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateRollBackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateRollBackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateRollBackLogic {
	return &UpdateTradeStateRollBackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateTradeStateRollBack 回滚更新交易状态
func (l *UpdateTradeStateRollBackLogic) UpdateTradeStateRollBack(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {

	// 检查输入合法性
	if len(in.TradeStateDesc) == 0 || len(in.TradeType) == 0 || len(in.Sn) == 0 || len(in.TransactionId) == 0 ||
		in.PayStatus < -1 || in.PayStatus > 2 || !in.PayTime.IsValid() {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input : %+v", in)
	}

	//1、流水记录确认.
	payment, err := l.svcCtx.PaymentModel.FindOneBySn(in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(ERRDBERR, "DBERR: 确认流水记录失败, err: %v, Sn: %s", err, in.Sn)
	}
	if payment == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("流水记录不存在"), "Not Found: Sn: %s", in.Sn)
	}

	return &pb.UpdateTradeStateResp{}, nil
}
