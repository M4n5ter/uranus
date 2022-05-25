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

func (s *ServiceContext) CopyUniqFlightsRpcRespToApiResp(resp []*types.UniqFlightWithSpaces, rpcResp []*pb.UniqFlightWithSpaces) {
	// 赋值
	for i := range rpcResp {
		//初始化FlightInfos以便下面进行赋值（避免空指针）
		resp[i] = &types.UniqFlightWithSpaces{}
		resp[i].FlightInfoID = rpcResp[i].FlightInfoID
		resp[i].SpacesOfFlightInfo = make([]types.Flightinfo, 2, 2)
		for spaceIndex, info := range rpcResp[i].SpacesOfFlightInfo {
			//退票信息
			for _, timeFee := range info.RefundInfo.TimeFees {
				t := strings.Split(timeFee.Time.AsTime().Local().String(), " +")[0]
				f := timeFee.Fee
				resp[i].SpacesOfFlightInfo[spaceIndex].RefundInfo.TFs = append(resp[i].SpacesOfFlightInfo[spaceIndex].RefundInfo.TFs, types.TF{
					Time: t,
					Fee:  f,
				})
			}

			//改票信息
			for _, timeFee := range info.ChangeInfo.TimeFees {
				t := strings.Split(timeFee.Time.AsTime().Local().String(), " +")[0]
				f := timeFee.Fee
				resp[i].SpacesOfFlightInfo[spaceIndex].ChangeInfo.TFs = append(resp[i].SpacesOfFlightInfo[spaceIndex].ChangeInfo.TFs, types.TF{
					Time: t,
					Fee:  f,
				})
			}

			resp[i].SpacesOfFlightInfo[spaceIndex].FlightInfoID = rpcResp[i].SpacesOfFlightInfo[spaceIndex].FlightInfoID
			resp[i].SpacesOfFlightInfo[spaceIndex].FlightNumber = rpcResp[i].SpacesOfFlightInfo[spaceIndex].FlightNumber
			resp[i].SpacesOfFlightInfo[spaceIndex].FlightType = rpcResp[i].SpacesOfFlightInfo[spaceIndex].FlightType
			resp[i].SpacesOfFlightInfo[spaceIndex].SetOutDate = strings.Split(strings.Split(rpcResp[i].SpacesOfFlightInfo[spaceIndex].SetOutDate.AsTime().Local().String(), " +")[0], " ")[0]
			resp[i].SpacesOfFlightInfo[spaceIndex].ArriveTime = strings.Split(rpcResp[i].SpacesOfFlightInfo[spaceIndex].ArriveTime.AsTime().Local().String(), " +")[0]
			resp[i].SpacesOfFlightInfo[spaceIndex].ArrivePosition = rpcResp[i].SpacesOfFlightInfo[spaceIndex].ArrivePosition
			resp[i].SpacesOfFlightInfo[spaceIndex].DepartTime = strings.Split(rpcResp[i].SpacesOfFlightInfo[spaceIndex].DepartTime.AsTime().Local().String(), " +")[0]
			resp[i].SpacesOfFlightInfo[spaceIndex].DepartPosition = rpcResp[i].SpacesOfFlightInfo[spaceIndex].DepartPosition
			resp[i].SpacesOfFlightInfo[spaceIndex].Price = rpcResp[i].SpacesOfFlightInfo[spaceIndex].Price
			resp[i].SpacesOfFlightInfo[spaceIndex].Surplus = rpcResp[i].SpacesOfFlightInfo[spaceIndex].Surplus
			resp[i].SpacesOfFlightInfo[spaceIndex].Punctuality = rpcResp[i].SpacesOfFlightInfo[spaceIndex].Punctuality
			resp[i].SpacesOfFlightInfo[spaceIndex].Discount = rpcResp[i].SpacesOfFlightInfo[spaceIndex].Discount
			resp[i].SpacesOfFlightInfo[spaceIndex].IsFirstClass = rpcResp[i].SpacesOfFlightInfo[spaceIndex].IsFirstClass
			resp[i].SpacesOfFlightInfo[spaceIndex].SpaceID = rpcResp[i].SpacesOfFlightInfo[spaceIndex].SpaceID
		}
	}
}
