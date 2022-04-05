package flightOrder

import (
	"context"

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

func (l *UserFlightOrderListLogic) UserFlightOrderList(req *types.UserFlightOrderListReq) (resp *types.UserFlightOrderListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
