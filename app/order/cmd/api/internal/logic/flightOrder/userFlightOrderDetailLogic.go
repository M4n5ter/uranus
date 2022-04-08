package flightOrder

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"uranus/app/order/cmd/api/internal/svc"
	"uranus/app/order/cmd/api/internal/types"
	"uranus/app/order/cmd/rpc/order"
	"uranus/common/ctxdata"
	"uranus/common/timeTools"
	"uranus/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFlightOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFlightOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserFlightOrderDetailLogic {
	return UserFlightOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserFlightOrderDetailLogic) UserFlightOrderDetail(req *types.UserFlightOrderDetailReq) (resp *types.UserFlightOrderDetailResp, err error) {
	// 检查输入是否合法
	if len(req.Sn) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "")
	}

	orderDetailResp, err := l.svcCtx.OrderRpcClient.FlightOrderDetail(l.ctx, &order.FlightOrderDetailReq{Sn: req.Sn})
	if err != nil {
		return nil, err
	}

	if orderDetailResp.FlightOrder.UserId != ctxdata.GetUidFromCtx(l.ctx) {
		return nil, errors.Wrapf(xerr.NewErrCodeMsg(400, "只能查看自己的订单"), "")
	}

	fltOrder := types.FlightOrder{}
	_ = copier.Copy(&fltOrder, &orderDetailResp.FlightOrder)
	fltOrder.DepartTime = timeTools.Timestamppb2TimeStringYMDhms(orderDetailResp.FlightOrder.DepartTime)
	fltOrder.ArriveTime = timeTools.Timestamppb2TimeStringYMDhms(orderDetailResp.FlightOrder.ArriveTime)
	fltOrder.CreateTime = timeTools.Timestamppb2TimeStringYMDhms(orderDetailResp.FlightOrder.CreateTime)

	resp = &types.UserFlightOrderDetailResp{FlightOrderDetail: types.FlightOrder{}}
	resp.FlightOrderDetail = fltOrder
	return
}
