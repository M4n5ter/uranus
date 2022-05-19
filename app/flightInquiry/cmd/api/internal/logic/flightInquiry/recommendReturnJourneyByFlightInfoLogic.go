package flightInquiry

import (
	"context"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"

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
	rpcResp, err := l.svcCtx.FlightInquiryClient.RecommendReturnJourneyByFlightInfo(l.ctx, &flightinquiry.RecommendReturnJourneyByFlightInfoReq{FlightInfoID: req.FlightInfoID})
	if err != nil {
		return nil, err
	}

	resp = &types.RecommendReturnJourneyByFlightInfoResp{}
	resp.UniqFlightWithSpaces = make([]*types.UniqFlightWithSpaces, len(rpcResp.UniqFlightWithSpaces))
	l.svcCtx.CopyUniqFlightsRpcRespToApiResp(resp.UniqFlightWithSpaces, rpcResp.UniqFlightWithSpaces)
	return
}
