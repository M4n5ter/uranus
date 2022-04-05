package flightOrder

import (
	"context"

	"uranus/app/order/cmd/api/internal/svc"
	"uranus/app/order/cmd/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
