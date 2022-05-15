package logic

import (
	"context"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendReturnJourneyByFlightInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecommendReturnJourneyByFlightInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendReturnJourneyByFlightInfoLogic {
	return &RecommendReturnJourneyByFlightInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecommendReturnJourneyByFlightInfo 根据指定航班信息提供返程机票推荐(仅支持直飞)
func (l *RecommendReturnJourneyByFlightInfoLogic) RecommendReturnJourneyByFlightInfo(in *pb.RecommendReturnJourneyByFlightInfoReq) (*pb.RecommendReturnJourneyByFlightInfoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.RecommendReturnJourneyByFlightInfoResp{}, nil
}
