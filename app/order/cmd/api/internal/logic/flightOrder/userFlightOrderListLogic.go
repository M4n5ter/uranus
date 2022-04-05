package flightOrder

import (
	"context"
	"uranus/app/order/cmd/rpc/order"
	"uranus/common/ctxdata"

	"uranus/app/order/cmd/api/internal/svc"
	"uranus/app/order/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFlightOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserFlightOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) UserFlightOrderListLogic {
	return UserFlightOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserFlightOrderList 分页查询订单列表
func (l *UserFlightOrderListLogic) UserFlightOrderList(req *types.UserFlightOrderListReq) (resp *types.UserFlightOrderListResp, err error) {
	listResp, err := l.svcCtx.OrderRpcClient.UserFlightOrderList(l.ctx, &order.UserFlightOrderListReq{
		LastId:      req.LastId,
		PageSize:    req.PageSize,
		UserId:      ctxdata.GetUidFromCtx(l.ctx),
		TraderState: req.TraderState,
	})
	if err != nil {
		return nil, err
	}

	list := make([]string, len(listResp.List))
	for i, _ := range list {
		list[i] = listResp.List[i].Sn
	}
	resp.FlightOrderList = list
	return
}
