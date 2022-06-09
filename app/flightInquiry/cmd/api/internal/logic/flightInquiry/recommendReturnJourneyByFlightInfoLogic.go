package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/common/timeTools"

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
	var sd time.Time
	if sd, err := timeTools.String2TimeYMD(req.SelectedDate); err == nil {
		if sd.IsZero() {
			return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : SetOutDate is zero time(maybe wrong layout).")
		}
	} else {
		return nil, errors.Wrapf(ERRIllegalInput, "time.Parse(layout,%s) ERR : %v.", req.SelectedDate, err)
	}

	rpcResp, err := l.svcCtx.FlightInquiryClient.RecommendReturnJourneyByFlightInfo(l.ctx, &flightinquiry.RecommendReturnJourneyByFlightInfoReq{FlightInfoID: req.FlightInfoID, SelectedDate: timestamppb.New(sd)})
	if err != nil {
		return nil, err
	}

	resp = &types.RecommendReturnJourneyByFlightInfoResp{}
	resp.UniqFlightWithSpaces = make([]*types.UniqFlightWithSpaces, len(rpcResp.UniqFlightWithSpaces))
	l.svcCtx.CopyUniqFlightsRpcRespToApiResp(resp.UniqFlightWithSpaces, rpcResp.UniqFlightWithSpaces)
	return
}
