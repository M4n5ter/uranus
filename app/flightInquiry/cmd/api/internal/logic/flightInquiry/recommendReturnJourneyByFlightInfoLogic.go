package flightInquiry

import (
	"context"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendReturnJourneyByFlightInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecommendReturnJourneyByFlightInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendReturnJourneyByFlightInfoLogic {
	return &RecommendReturnJourneyByFlightInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecommendReturnJourneyByFlightInfoLogic) RecommendReturnJourneyByFlightInfo(req *types.RecommendReturnJourneyByFlightInfoReq) (resp *types.RecommendReturnJourneyByFlightInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
