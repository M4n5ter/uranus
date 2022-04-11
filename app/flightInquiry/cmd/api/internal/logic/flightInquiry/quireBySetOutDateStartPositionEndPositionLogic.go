package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutDateStartPositionEndPositionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuireBySetOutDateStartPositionEndPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) QuireBySetOutDateStartPositionEndPositionLogic {
	return QuireBySetOutDateStartPositionEndPositionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuireBySetOutDateStartPositionEndPositionLogic) QuireBySetOutDateStartPositionEndPosition(req *types.QuireBySetOutDateStartPositionEndPositionReq) (resp *types.QuireBySetOutDateStartPositionEndPositionResp, err error) {
	//检查输入合法性
	if len(req.SetOutDate) == 0 {
		return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : SetOutDate is empty.")
	}
	var sod time.Time
	if sod, err = time.Parse("2006-01-02", req.SetOutDate); err == nil {
		if sod.IsZero() {
			return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : SetOutDate is zero time(maybe wrong layout).")
		}
	} else {
		return nil, errors.Wrapf(ERRIllegalInput, "time.Parse(layout,%s) ERR : %v.", req.SetOutDate, err)
	}
	if len(req.DepartPosition) == 0 || len(req.ArrivePosition) == 0 {
		return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : DepartPosition or ArrivePosition is empty.")
	}

	rpcResp, err := l.svcCtx.FlightInquiryClient.QuireBySetOutDateStartPositionEndPosition(l.ctx, &flightinquiry.QuireBySetOutDateStartPositionEndPositionReq{
		SetOutDate:     timestamppb.New(sod),
		ArrivePosition: req.ArrivePosition,
		DepartPosition: req.DepartPosition,
	})
	if err != nil {
		return nil, errors.Wrapf(ERRRpcCall, "Rpc err in callling flightInquiry-rpc.QuireBySetOutDateStartPositionEndPosition, ArrivePosition: %s, DepartPosition: %s, SetOutDate: %v, err: %v", req.ArrivePosition, req.DepartPosition, req.SetOutDate, err)
	}
	if rpcResp.FlightInfos == nil {
		return nil, errors.Wrapf(ERRRpcCall, "the rpcResp.FlightInfos is nil: l.svcCtx.FlightInquiryClient.QuireBySetOutDateStartPositionEndPosition, ArrivePosition: %s, DepartPosition: %s, SetOutDate: %v", req.ArrivePosition, req.DepartPosition, req.SetOutDate)
	}
	//初始化resp，避免空指针错误
	resp = &types.QuireBySetOutDateStartPositionEndPositionResp{}
	resp.Flightinfos = make([]*types.Flightinfo, len(rpcResp.FlightInfos))
	l.svcCtx.CopyFlightInfosRpcRespToApiResp(resp.Flightinfos, rpcResp.FlightInfos)
	return
}
