package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/order/model"
	"uranus/common/xerr"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeTheOrderPriceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeTheOrderPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeTheOrderPriceLogic {
	return &ChangeTheOrderPriceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 机票改价
func (l *ChangeTheOrderPriceLogic) ChangeTheOrderPrice(in *pb.ChangeTheOrderPriceReq) (*pb.ChangeTheOrderPriceResp, error) {
	if len(in.OrderSn) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "")
	}

	order, err := l.svcCtx.OrderModel.FindOneBySn(in.OrderSn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: orderSn: %s, err: %v", in.OrderSn, err)
	}

	if order == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("订单不存在"), "Err Not Found: orderSn:%s", in.OrderSn)
	}

	order.OrderTotalPrice = in.Price
	err = l.svcCtx.OrderModel.UpdateWithVersion(nil, order)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("改价失败"), "orderSn: %s, price: %d, err: %v", in.OrderSn, in.Price, err)
	}

	return &pb.ChangeTheOrderPriceResp{}, nil
}
