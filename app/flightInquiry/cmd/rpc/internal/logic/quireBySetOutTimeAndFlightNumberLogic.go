package logic

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/common/xerr"
	"uranus/model"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRGetInfos = xerr.NewErrMsg("暂无直飞航班")
var ERRGetSpaces = xerr.NewErrMsg("暂无舱位信息")
var ERRGetTickets = xerr.NewErrMsg("暂无票信息")

type QuireBySetOutTimeAndFlightNumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuireBySetOutTimeAndFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireBySetOutTimeAndFlightNumberLogic {
	return &QuireBySetOutTimeAndFlightNumberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuireBySetOutTimeAndFlightNumber 通过给定日期、航班号进行航班查询请求
func (l *QuireBySetOutTimeAndFlightNumberLogic) QuireBySetOutTimeAndFlightNumber(in *pb.QuireBySetOutTimeAndFlightNumberReq) (*pb.QuireBySetOutTimeAndFlightNumberResp, error) {
	resp := pb.QuireBySetOutTimeAndFlightNumberResp{}
	//查询 FlightNumber SetOutDate Punctuality DepartPosition DepartTime ArrivePosition ArriveTime
	flightInfos, err := l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutTime(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightNumber, in.SetOutTime.AsTime())
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Wrapf(ERRGetInfos, "NOT FOUND: can't found flight infos: number->%s setOutTime->%v, ERR: %v", in.FlightNumber, in.SetOutTime.AsTime(), err)
		} else {
			return nil, errors.Wrapf(ERRGetInfos, "DBERR: when calling flightinquiry-rpc:l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutTime : number->%s setOutTime->%v, ERR: %v", in.FlightNumber, in.SetOutTime.AsTime(), err)
		}
	}

	//查询 IsFirstClass Surplus
	for _, info := range flightInfos {
		spaces, err := l.svcCtx.SpacesModel.FindListByFlightInfoID(l.svcCtx.SpacesModel.RowBuilder(), info.Id)
		if err != nil {
			if err == model.ErrNotFound {
				logx.WithContext(l.ctx).Infof("NOT FOUND: There is no corresponding space information for this flightInfo.FlightInfoID:%d", info.Id)
			}
			return nil, errors.Wrapf(ERRGetSpaces, "DBERR: when calling flightinquiry-rpc:l.svcCtx.SpacesModel.FindListByFlightInfoID : FlightInfoID:%d", info.Id)
		}
		for _, space := range spaces {
			// 是否是头等舱/商务舱
			var ifc bool
			if space.IsFirstClass == 0 {
				ifc = false
			} else {
				ifc = true
			}
			// 查询 Price RefundInfo ChangeInfo
			tickets, err := l.svcCtx.TicketsModel.FindListBySpaceID(l.svcCtx.TicketsModel.RowBuilder(), space.Id)
			if err != nil {
				if err == model.ErrNotFound {
					logx.WithContext(l.ctx).Infof("NOT FOUND: There is no ticket information for the corresponding space.spaceID:%d", space.Id)
				}
				return nil, errors.Wrapf(ERRGetTickets, "DBERR: when calling flightinquiry-rpc:l.svcCtx.TicketsModel.FindListBySpaceID : spaceID:%d", space.Id)
			}
			for _, ticket := range tickets {
				// 添加对应信息
				resp.FlightInfos = append(resp.FlightInfos, &pb.FlightInfo{
					FlightNumber:   info.FlightNumber,
					SetOutDate:     timestamppb.New(info.SetOutDate),
					IsFirstClass:   ifc,
					Price:          uint64(ticket.Price),
					Discount:       ticket.Discount,
					Surplus:        space.Surplus,
					Punctuality:    uint32(info.Punctuality),
					DepartPosition: info.DepartPosition,
					DepartTime:     timestamppb.New(info.DepartTime),
					ArrivePosition: info.ArrivePosition,
					ArriveTime:     timestamppb.New(info.ArriveTime),
					RefundInfo:     nil,
					ChangeInfo:     nil,
					Cba:            ticket.Cba,
				})
			}

		}

	}
	return &resp, nil
}
