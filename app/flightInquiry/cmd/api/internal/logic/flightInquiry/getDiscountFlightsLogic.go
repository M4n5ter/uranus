package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDiscountFlightsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDiscountFlightsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiscountFlightsLogic {
	return &GetDiscountFlightsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDiscountFlightsLogic) GetDiscountFlights(req *types.GetDiscountFlightsReq) (resp *types.GetDiscountFlightsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
