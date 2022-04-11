package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"strings"
	"uranus/app/flightInquiry/cmd/api/internal/config"
	"uranus/app/flightInquiry/cmd/api/internal/types"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/app/flightInquiry/cmd/rpc/pb"
)

type ServiceContext struct {
	Config              config.Config
	FlightInquiryClient flightinquiry.FlightInquiry
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		FlightInquiryClient: flightinquiry.NewFlightInquiry(zrpc.MustNewClient(c.FlightInquiryConf)),
	}
}

func (s *ServiceContext) CopyFlightInfosRpcRespToApiResp(resp []*types.Flightinfo, rpcResp []*pb.FlightInfo) {
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
