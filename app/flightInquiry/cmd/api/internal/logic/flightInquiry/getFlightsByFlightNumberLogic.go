package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFlightsByFlightNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFlightsByFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightsByFlightNumberLogic {
	return &GetFlightsByFlightNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFlightsByFlightNumberLogic) GetFlightsByFlightNumber(req *types.GetFlightsByFlightNumberReq) (resp *types.GetFlightsByFlightNumberResp, err error) {
	if len(req.FlightNumber) == 0 {
		return nil, errors.Wrapf(ERRIllegalInput, "")
	}

	rpcResp, err := l.svcCtx.FlightInquiryClient.GetFlightsByFlightNumber(l.ctx, &flightinquiry.GetFlightsByFlightNumberReq{
		FlightNumber: req.FlightNumber,
		Days:         req.Days,
		Num:          req.Num,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetFlightsByFlightNumberResp{}
	resp.UniqFlightWithSpaces = make([]*types.UniqFlightWithSpaces, len(rpcResp.UniqFlightWithSpaces))
	l.svcCtx.CopyUniqFlightsRpcRespToApiResp(resp.UniqFlightWithSpaces, rpcResp.UniqFlightWithSpaces)
	return
}
