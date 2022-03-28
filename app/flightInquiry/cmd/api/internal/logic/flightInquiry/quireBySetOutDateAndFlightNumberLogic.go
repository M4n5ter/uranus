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

var ERRIllegalInput = xerr.NewErrMsg("非法输入")
var ERRRpcCall = xerr.NewErrMsg("获取航班信息失败")

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
		return nil, errors.Wrapf(ERRIllegalInput, "Illegal input : FlightNumber is empty.")
	}
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
	rpcResp, err := l.svcCtx.FlightInquiryClient.QuireBySetOutDateAndFlightNumber(l.ctx, &flightinquiry.QuireBySetOutDateAndFlightNumberReq{FlightNumber: req.FlightNumber, SetOutDate: timestamppb.New(sod)})
	if err != nil {
		return nil, errors.Wrapf(ERRRpcCall, "Rpc err in callling flightInquiry-rpc.QuireBySetOutDateAndFlightNumber, FlightNumber: %s, SetOutDate: %v, err: %v", req.FlightNumber, req.SetOutDate, err)
	}
	// 赋值
	// 初始化resp，避免空指针
	resp = &types.QuireBySetOutDateAndFlightNumberResp{}
	resp.Flightinfos = make([]*types.Flightinfo, len(rpcResp.FlightInfos))
	for i, _ := range resp.Flightinfos {
		//初始化FlightInfos以便下面进行赋值（避免空指针）
		resp.Flightinfos[i] = &types.Flightinfo{}
		//退票信息
		for _, timeFee := range rpcResp.FlightInfos[i].RefundInfo.TimeFees {
			t := timeFee.Time.AsTime().String()
			f := timeFee.Fee
			resp.Flightinfos[i].RefundInfo.TFs = append(resp.Flightinfos[i].RefundInfo.TFs, types.TF{
				Time: t,
				Fee:  f,
			})
		}
		//改票信息
		for _, timeFee := range rpcResp.FlightInfos[i].ChangeInfo.TimeFees {
			t := timeFee.Time.AsTime().String()
			f := timeFee.Fee
			resp.Flightinfos[i].ChangeInfo.TFs = append(resp.Flightinfos[i].ChangeInfo.TFs, types.TF{
				Time: t,
				Fee:  f,
			})
		}
		resp.Flightinfos[i].FlightNumber = rpcResp.FlightInfos[i].FlightNumber
		resp.Flightinfos[i].SetOutDate = rpcResp.FlightInfos[i].SetOutDate.AsTime().String()
		resp.Flightinfos[i].ArriveTime = rpcResp.FlightInfos[i].ArriveTime.AsTime().String()
		resp.Flightinfos[i].ArrivePosition = rpcResp.FlightInfos[i].ArrivePosition
		resp.Flightinfos[i].DepartTime = rpcResp.FlightInfos[i].DepartTime.AsTime().String()
		resp.Flightinfos[i].DepartPosition = rpcResp.FlightInfos[i].DepartPosition
		resp.Flightinfos[i].Price = rpcResp.FlightInfos[i].Price
		resp.Flightinfos[i].Surplus = rpcResp.FlightInfos[i].Surplus
		resp.Flightinfos[i].Punctuality = rpcResp.FlightInfos[i].Punctuality
		resp.Flightinfos[i].Discount = rpcResp.FlightInfos[i].Discount
		resp.Flightinfos[i].IsFirstClass = rpcResp.FlightInfos[i].IsFirstClass
	}
	return
}
