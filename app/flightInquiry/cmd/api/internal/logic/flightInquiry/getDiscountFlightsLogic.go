package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"

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
	if len(req.DepartPosition) == 0 || len(req.ArrivePosition) == 0 {
		return nil, errors.Wrapf(ERRIllegalInput, "")
	}

	rpcResp, err := l.svcCtx.FlightInquiryClient.GetDiscountFlights(l.ctx, &flightinquiry.GetDiscountFlightsReq{
		DepartPosition: req.DepartPosition,
		ArrivePosition: req.ArrivePosition,
		Days:           req.Days,
		Num:            req.Num,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetDiscountFlightsResp{}
	resp.UniqFlightWithSpaces = make([]*types.UniqFlightWithSpaces, len(rpcResp.UniqFlightWithSpaces))
	l.svcCtx.CopyUniqFlightsRpcRespToApiResp(resp.UniqFlightWithSpaces, rpcResp.UniqFlightWithSpaces)
	return
}
