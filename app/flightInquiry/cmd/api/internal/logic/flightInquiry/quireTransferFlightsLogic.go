package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/common/timeTools"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireTransferFlightsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuireTransferFlightsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireTransferFlightsLogic {
	return &QuireTransferFlightsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuireTransferFlightsLogic) QuireTransferFlights(req *types.QuireTransferFlightsReq) (resp *types.QuireTransferFlightsResp, err error) {
	var sod time.Time
	if sod, err = timeTools.String2TimeYMD(req.SetOutDate); err == nil {
		if sod.IsZero() {
			return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : SetOutDate is zero time(maybe wrong layout).")
		}
	} else {
		return nil, errors.Wrapf(ERRIllegalInput, "time.Parse(layout,%s) ERR : %v.", req.SetOutDate, err)
	}
	if len(req.DepartPosition) == 0 || len(req.ArrivePosition) == 0 {
		return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : DepartPosition or ArrivePosition is empty.")
	}

	rpcResp, err := l.svcCtx.FlightInquiryClient.QuireTransferFlights(l.ctx, &flightinquiry.QuireTransferFlightsReq{
		DepartPosition: req.DepartPosition,
		ArrivePosition: req.ArrivePosition,
		SetOutDate:     timestamppb.New(sod),
	})
	if err != nil {
		return nil, err
	}

	resp = &types.QuireTransferFlightsResp{}
	resp.Transfers = make([]types.TransferFlights, len(rpcResp.TransferFlights))
	for i, flight := range rpcResp.TransferFlights {
		uniq := make([]*types.UniqFlightWithSpaces, len(flight.UniqFlightWithSpaces))
		l.svcCtx.CopyUniqFlightsRpcRespToApiResp(uniq, flight.UniqFlightWithSpaces)
		resp.Transfers[i].UniqFlightWithSpaceses = uniq
	}
	return
}
