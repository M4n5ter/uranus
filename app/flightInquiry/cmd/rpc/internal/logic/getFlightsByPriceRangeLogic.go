package logic

import (
	"context"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFlightsByPriceRangeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFlightsByPriceRangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightsByPriceRangeLogic {
	return &GetFlightsByPriceRangeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFlightsByPriceRange 查询指定价格区间的航班
func (l *GetFlightsByPriceRangeLogic) GetFlightsByPriceRange(in *pb.GetFlightsByPriceRangeReq) (*pb.GetFlightsByPriceRangeResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetFlightsByPriceRangeResp{}, nil
}
