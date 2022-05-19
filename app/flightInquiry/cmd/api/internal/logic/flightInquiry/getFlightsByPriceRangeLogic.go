package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFlightsByPriceRangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFlightsByPriceRangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightsByPriceRangeLogic {
	return &GetFlightsByPriceRangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFlightsByPriceRangeLogic) GetFlightsByPriceRange(req *types.GetFlightsByPriceRangeReq) (resp *types.GetFlightsByPriceRangeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
