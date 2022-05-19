package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/common/timeTools"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFlightDetailByTicketIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFlightDetailByTicketIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightDetailByTicketIDLogic {
	return &GetFlightDetailByTicketIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFlightDetailByTicketIDLogic) GetFlightDetailByTicketID(req *types.GetFlightDetailByTicketIDReq) (resp *types.GetFlightDetailByTicketIDResp, err error) {
	if req.TicketID < 1 {
		return nil, errors.Wrapf(ERRIllegalInput, "")
	}

	rpcResp, err := l.svcCtx.FlightInquiryClient.FlightDetail(l.ctx, &flightinquiry.FlightDetailReq{TicketID: req.TicketID})
	if err != nil {
		return nil, err
	}

	resp = &types.GetFlightDetailByTicketIDResp{}
	resp.Flightinfo = types.Flightinfo{
		FlightInfoID:   rpcResp.FlightInfo.FlightInfoID,
		FlightNumber:   rpcResp.FlightInfo.FlightNumber,
		FlightType:     rpcResp.FlightInfo.FlightType,
		SetOutDate:     timeTools.Timestamppb2TimeStringYMD(rpcResp.FlightInfo.SetOutDate),
		IsFirstClass:   rpcResp.FlightInfo.IsFirstClass,
		SpaceID:        rpcResp.FlightInfo.SpaceID,
		Price:          rpcResp.FlightInfo.Price,
		Discount:       rpcResp.FlightInfo.Discount,
		Surplus:        rpcResp.FlightInfo.Surplus,
		Punctuality:    rpcResp.FlightInfo.Punctuality,
		DepartPosition: rpcResp.FlightInfo.DepartPosition,
		DepartTime:     timeTools.Timestamppb2TimeStringYMDhms(rpcResp.FlightInfo.DepartTime),
		ArrivePosition: rpcResp.FlightInfo.ArrivePosition,
		ArriveTime:     timeTools.Timestamppb2TimeStringYMDhms(rpcResp.FlightInfo.ArriveTime),
	}

	ri := make([]types.TF, len(rpcResp.FlightInfo.RefundInfo.TimeFees))
	for i, tf := range rpcResp.FlightInfo.RefundInfo.TimeFees {
		ri[i].Fee = tf.Fee
		ri[i].Time = timeTools.Timestamppb2TimeStringYMDhms(tf.Time)
	}
	ci := make([]types.TF, len(rpcResp.FlightInfo.ChangeInfo.TimeFees))
	for i, tf := range rpcResp.FlightInfo.ChangeInfo.TimeFees {
		ci[i].Fee = tf.Fee
		ci[i].Time = timeTools.Timestamppb2TimeStringYMDhms(tf.Time)
	}

	resp.Flightinfo.RefundInfo.TFs, resp.Flightinfo.ChangeInfo.TFs = ri, ci
	return resp, nil
}
