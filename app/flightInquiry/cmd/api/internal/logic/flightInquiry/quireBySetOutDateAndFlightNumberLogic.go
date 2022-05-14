package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"

	"time"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"

	"uranus/common/xerr"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRIllegalInput = xerr.NewErrCodeMsg(400, "非法输入")
var ERRRpcCall = xerr.NewErrCodeMsg(xerr.DB_ERROR, "获取航班信息失败")
var ERRNotFound = xerr.NewErrCodeMsg(400, "暂无对应航班信息")

type QuireBySetOutDateAndFlightNumberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuireBySetOutDateAndFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) QuireBySetOutDateAndFlightNumberLogic {
	return QuireBySetOutDateAndFlightNumberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuireBySetOutDateAndFlightNumberLogic) QuireBySetOutDateAndFlightNumber(req *types.QuireBySetOutDateAndFlightNumberReq) (resp *types.QuireBySetOutDateAndFlightNumberResp, err error) {
	//检查输入的合法性
	if len(req.FlightNumber) == 0 {
		return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : FlightNumber is empty.\n")
	}
	if len(req.SetOutDate) == 0 {
		return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : SetOutDate is empty.\n")
	}
	var sod time.Time
	if sod, err = time.Parse("2006-01-02", req.SetOutDate); err == nil {
		if sod.IsZero() {
			return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : SetOutDate is zero time(maybe wrong layout).\n")
		}
	} else {
		return nil, errors.Wrapf(ERRIllegalInput, "time.Parse(layout,%s) ERR : %v.\n", req.SetOutDate, err)
	}
	rpcResp, err := l.svcCtx.FlightInquiryClient.QuireBySetOutDateAndFlightNumber(l.ctx, &flightinquiry.QuireBySetOutDateAndFlightNumberReq{FlightNumber: req.FlightNumber, SetOutDate: timestamppb.New(sod)})
	if err != nil {
		return nil, errors.Wrapf(ERRRpcCall, "Rpc err in callling flightInquiry-rpc.QuireBySetOutDateAndFlightNumber, FlightNumber: %s, SetOutDate: %v, err: %v\n", req.FlightNumber, req.SetOutDate, err)
	}
	if rpcResp.UniqFlightWithSpaces == nil {
		return nil, errors.Wrapf(ERRNotFound, "the rpcResp.FlightInfos is nil: l.svcCtx.FlightInquiryClient.QuireBySetOutDateAndFlightNumber, FlightNumber: %s, SetOutDate: %v\n", req.FlightNumber, req.SetOutDate)
	}
	// 初始化resp，避免空指针
	resp = &types.QuireBySetOutDateAndFlightNumberResp{}
	resp.UniqFlightWithSpaces = make([]*types.UniqFlightWithSpaces, len(rpcResp.UniqFlightWithSpaces))
	l.svcCtx.CopyFlightInfosRpcRespToApiResp(resp.UniqFlightWithSpaces, rpcResp.UniqFlightWithSpaces)
	return
}
