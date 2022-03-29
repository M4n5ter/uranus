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
var ERRGetFltType = xerr.NewErrMsg("找不到对应航班机型")
var ERRGetTickets = xerr.NewErrMsg("暂无票信息")
var ERRRefundAndChangeInfos = xerr.NewErrMsg("暂无退改票信息")

type QuireBySetOutDateAndFlightNumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuireBySetOutDateAndFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireBySetOutDateAndFlightNumberLogic {
	return &QuireBySetOutDateAndFlightNumberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuireBySetOutDateAndFlightNumber 通过给定日期、航班号进行航班查询请求
func (l *QuireBySetOutDateAndFlightNumberLogic) QuireBySetOutDateAndFlightNumber(in *pb.QuireBySetOutDateAndFlightNumberReq) (*pb.QuireBySetOutDateAndFlightNumberResp, error) {
	resp := pb.QuireBySetOutDateAndFlightNumberResp{}
	//查询 FlightNumber SetOutDate Punctuality DepartPosition DepartTime ArrivePosition ArriveTime
	flightInfos, err := l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightNumber, in.SetOutDate.AsTime())
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.Wrapf(ERRGetInfos, "NOT FOUND: can't found flight infos: number->%s setOutTime->%v, ERR: %v", in.FlightNumber, in.SetOutDate.AsTime(), err)
		} else {
			return nil, errors.Wrapf(ERRGetInfos, "DBERR: when calling flightinquiry-rpc:l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate : number->%s setOutTime->%v, ERR: %v", in.FlightNumber, in.SetOutDate.AsTime(), err)
		}
	}

	//查询 IsFirstClass Surplus FlightTypes
	for _, info := range flightInfos {
		flt, err := l.svcCtx.Flights.FindOneByNumber(info.FlightNumber)
		if err != nil {
			if err == model.ErrNotFound {
				logx.WithContext(l.ctx).Infof("NOT FOUND: There is no corresponding flight information for number in this flightInfo.FlightNumber:%s", info.FlightNumber)
			}
			return nil, errors.Wrapf(ERRGetFltType, "DBERR: when calling flightinquiry-rpc:l.svcCtx.Flights.FindOneByNumber : FlightNumber:%s, err: %v", info.FlightNumber, err)
		}
		spaces, err := l.svcCtx.SpacesModel.FindListByFlightInfoID(l.svcCtx.SpacesModel.RowBuilder(), info.Id)
		if err != nil {
			if err == model.ErrNotFound {
				logx.WithContext(l.ctx).Infof("NOT FOUND: There is no corresponding space information for this flightInfo.FlightInfoID:%d", info.Id)
			}
			return nil, errors.Wrapf(ERRGetSpaces, "DBERR: when calling flightinquiry-rpc:l.svcCtx.SpacesModel.FindListByFlightInfoID : FlightInfoID:%d, err: %v", info.Id, err)
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
				//退改票信息
				refundInfo := &pb.RefundInfo{}
				changeInfo := &pb.ChangeInfo{}
				rcis, err := l.svcCtx.RefundAndChangeInfosModel.FindListByTicketID(l.svcCtx.RefundAndChangeInfosModel.RowBuilder(), ticket.Id)
				if err != nil {
					if err == model.ErrNotFound {
						logx.WithContext(l.ctx).Infof("NOT FOUND: There is no refund and change information for the corresponding ticket.ticketID:%d", ticket.Id)
					}
					return nil, errors.Wrapf(ERRRefundAndChangeInfos, "DBERR: when calling flightinquiry-rpc:l.svcCtx.RefundAndChangeInfosModel.FindListByTicketID : ticketID:%d", ticket.Id)
				}
				if rcis != nil {
					for _, rci := range rcis {
						switch rci.IsRefund {
						//如果是改票信息
						case 0:
							changeInfo.TimeFees = append(changeInfo.TimeFees,
								&pb.TimeFee{Time: timestamppb.New(rci.Time1), Fee: uint64(rci.Fee1)},
								&pb.TimeFee{Time: timestamppb.New(rci.Time2), Fee: uint64(rci.Fee2)},
								&pb.TimeFee{Time: timestamppb.New(rci.Time3), Fee: uint64(rci.Fee3)},
								&pb.TimeFee{Time: timestamppb.New(rci.Time4), Fee: uint64(rci.Fee4)},
							)
							if rci.Time5.Valid && rci.Fee5.Valid {
								changeInfo.TimeFees = append(changeInfo.TimeFees, &pb.TimeFee{Time: timestamppb.New(rci.Time5.Time), Fee: uint64(rci.Fee5.Int64)})
							}
						//如果是退票信息
						case 1:
							refundInfo.TimeFees = append(refundInfo.TimeFees,
								&pb.TimeFee{Time: timestamppb.New(rci.Time1), Fee: uint64(rci.Fee1)},
								&pb.TimeFee{Time: timestamppb.New(rci.Time2), Fee: uint64(rci.Fee2)},
								&pb.TimeFee{Time: timestamppb.New(rci.Time3), Fee: uint64(rci.Fee3)},
								&pb.TimeFee{Time: timestamppb.New(rci.Time4), Fee: uint64(rci.Fee4)},
							)
							if rci.Time5.Valid && rci.Fee5.Valid {
								refundInfo.TimeFees = append(refundInfo.TimeFees, &pb.TimeFee{Time: timestamppb.New(rci.Time5.Time), Fee: uint64(rci.Fee5.Int64)})
							}
						default:
						}

					}

				}

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
					RefundInfo:     refundInfo,
					ChangeInfo:     changeInfo,
					Cba:            ticket.Cba,
					FlightType:     flt.FltType,
				})
			}

		}

	}
	return &resp, nil
}
