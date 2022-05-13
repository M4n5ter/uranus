package logic

import (
	"context"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFlightsByFlightNumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFlightsByFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightsByFlightNumberLogic {
	return &GetFlightsByFlightNumberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFlightsByFlightNumber 根据航班号获取航班信息
func (l *GetFlightsByFlightNumberLogic) GetFlightsByFlightNumber(in *pb.GetFlightsByFlightNumberReq) (*pb.GetFlightsByFlightNumberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetFlightsByFlightNumberResp{}, nil
}
