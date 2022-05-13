package logic

import (
	"context"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDiscountFlightsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDiscountFlightsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiscountFlightsLogic {
	return &GetDiscountFlightsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDiscountFlights 查询折扣航班
func (l *GetDiscountFlightsLogic) GetDiscountFlights(in *pb.GetDiscountFlightsReq) (*pb.GetDiscountFlightsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetDiscountFlightsResp{}, nil
}
