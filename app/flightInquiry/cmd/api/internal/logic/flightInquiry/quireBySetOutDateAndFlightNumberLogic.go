package flightInquiry

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/app/flightInquiry/cmd/rpc/pb"
	"uranus/common/xerr"

	"uranus/app/flightInquiry/cmd/api/internal/svc"
	"uranus/app/flightInquiry/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRIllegalInput = xerr.NewErrCodeMsg(400, "非法输入")
var ERRRpcCall = xerr.NewErrCodeMsg(500, "获取航班信息失败")
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
	if rpcResp.FlightInfos == nil {
		return nil, errors.Wrapf(ERRNotFound, "the rpcResp.FlightInfos is nil: l.svcCtx.FlightInquiryClient.QuireBySetOutDateAndFlightNumber, FlightNumber: %s, SetOutDate: %v\n", req.FlightNumber, req.SetOutDate)
	}
	// 初始化resp，避免空指针
	resp = &types.QuireBySetOutDateAndFlightNumberResp{}
	resp.Flightinfos = make([]*types.Flightinfo, len(rpcResp.FlightInfos))
	copyFlightInfosRpcRespToApiResp(resp.Flightinfos, rpcResp.FlightInfos)
	return
}

func copyFlightInfosRpcRespToApiResp(resp []*types.Flightinfo, rpcResp []*pb.FlightInfo) {
	// 赋值
	for i := range resp {
		//初始化FlightInfos以便下面进行赋值（避免空指针）
		resp[i] = &types.Flightinfo{}
		//退票信息
		for _, timeFee := range rpcResp[i].RefundInfo.TimeFees {
			t := strings.Split(timeFee.Time.AsTime().Local().String(), " +")[0]
			f := timeFee.Fee
			resp[i].RefundInfo.TFs = append(resp[i].RefundInfo.TFs, types.TF{
				Time: t,
				Fee:  f,
			})
		}
		//改票信息
		for _, timeFee := range rpcResp[i].ChangeInfo.TimeFees {
			t := strings.Split(timeFee.Time.AsTime().Local().String(), " +")[0]
			f := timeFee.Fee
			resp[i].ChangeInfo.TFs = append(resp[i].ChangeInfo.TFs, types.TF{
				Time: t,
				Fee:  f,
			})
		}
		resp[i].FlightNumber = rpcResp[i].FlightNumber
		resp[i].FlightType = rpcResp[i].FlightType
		resp[i].SetOutDate = strings.Split(strings.Split(rpcResp[i].SetOutDate.AsTime().Local().String(), " +")[0], " ")[0]
		resp[i].ArriveTime = strings.Split(rpcResp[i].ArriveTime.AsTime().Local().String(), " +")[0]
		resp[i].ArrivePosition = rpcResp[i].ArrivePosition
		resp[i].DepartTime = strings.Split(rpcResp[i].DepartTime.AsTime().Local().String(), " +")[0]
		resp[i].DepartPosition = rpcResp[i].DepartPosition
		resp[i].Price = rpcResp[i].Price
		resp[i].Surplus = rpcResp[i].Surplus
		resp[i].Punctuality = rpcResp[i].Punctuality
		resp[i].Discount = rpcResp[i].Discount
		resp[i].IsFirstClass = rpcResp[i].IsFirstClass
	}
}
