package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"

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
	if len(req.DepartPosition) == 0 || len(req.ArrivePosition) == 0 {
		return nil, errors.Wrapf(ERRIllegalInput, "")
	}

	rpcResp, err := l.svcCtx.FlightInquiryClient.GetFlightsByPriceRange(l.ctx, &flightinquiry.GetFlightsByPriceRangeReq{
		DepartPosition: req.DepartPosition,
		ArrivePosition: req.ArrivePosition,
		MinPrice:       req.MinPrice,
		MaxPrice:       req.MaxPrice,
		Days:           req.Days,
		Num:            req.Num,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetFlightsByPriceRangeResp{}
	resp.UniqFlightWithSpaces = make([]*types.UniqFlightWithSpaces, len(rpcResp.UniqFlightWithSpaces))
	l.svcCtx.CopyUniqFlightsRpcRespToApiResp(resp.UniqFlightWithSpaces, rpcResp.UniqFlightWithSpaces)
	return
}
