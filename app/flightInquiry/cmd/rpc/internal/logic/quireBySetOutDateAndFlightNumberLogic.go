package logic

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRGetInfos = xerr.NewErrMsg("暂无直飞航班")
var ERRGetSpaces = xerr.NewErrMsg("暂无舱位信息")
var ERRGetTickets = xerr.NewErrMsg("暂无票信息")
var ERRRefundAndChangeInfos = xerr.NewErrMsg("暂无退改票信息")
var ERRDBERR = xerr.NewErrCode(xerr.DB_ERROR)

//var ERRGetFltType = xerr.NewErrMsg("找不到对应航班机型")

type QuireBySetOutDateAndFlightNumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type FlightQuirer struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type FlightInfoResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	FlightInfos   []*pb.FlightInfo `protobuf:"bytes,1,rep,name=FlightInfos,proto3" json:"FlightInfos,omitempty"`
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
	resp := &pb.QuireBySetOutDateAndFlightNumberResp{}
	// 不走缓存查询 FlightNumber SetOutDate Punctuality DepartPosition DepartTime ArrivePosition ArriveTime
	flightInfos, err := l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightNumber, in.SetOutDate.AsTime())
	if err != nil {
		if err == commonModel.ErrNotFound {
			return nil, errors.Wrapf(ERRGetInfos, "NOT FOUND: can't found flight infos: number->%s setOutTime->%v, ERR: %v\n", in.FlightNumber, in.SetOutDate.AsTime(), err)
		} else {
			return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate : number->%s setOutTime->%v, ERR: %v\n", in.FlightNumber, in.SetOutDate.AsTime(), err)
		}
	}
	v, err := (*FlightQuirer)(l).combineAllInfos(flightInfos, resp, true)
	resp, ok := v.(*pb.QuireBySetOutDateAndFlightNumberResp)
	if !ok {
		return &pb.QuireBySetOutDateAndFlightNumberResp{}, nil
	}
	return resp, err
}

// combineAllInfos respType 为 true 时，将 resp 转化为 *pb.QuireBySetOutDateAndFlightNumberResp,否则转化为 *pb.QuireBySetOutDateStartPositionEndPositionResp
func (l *FlightQuirer) combineAllInfos(flightInfos []*commonModel.FlightInfos, resp any, respType bool) (any, error) {
	// 查询 IsFirstClass Surplus FlightTypes
	for _, info := range flightInfos {
		flt, err := l.svcCtx.Flights.FindOneByNumber(info.FlightNumber)
		if err != nil {
			if err == commonModel.ErrNotFound {
				flt = &commonModel.Flights{FltType: "unknown"}
				logx.WithContext(l.ctx).Infof("NOT FOUND: There is no corresponding flight information for number in this flightInfo.FlightNumber:%s\n", info.FlightNumber)
			} else {
				return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.Flights.FindOneByNumber : FlightNumber:%s, err: %v\n", info.FlightNumber, err)
			}
		}
		// 此处查询不走缓存，直接打到DB上
		spaces, err := l.svcCtx.SpacesModel.FindListByFlightInfoID(l.svcCtx.SpacesModel.RowBuilder(), info.Id)
		if err != nil {
			if err == commonModel.ErrNotFound {
				logx.WithContext(l.ctx).Infof("NOT FOUND: There is no corresponding space information for this flightInfo.FlightInfoID:%d\n", info.Id)
				return nil, errors.Wrapf(ERRGetSpaces, "NOT FOUND: There is no corresponding space information for this flightInfo.FlightInfoID:%d\n", info.Id)
			} else {
				return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.SpacesModel.FindListByFlightInfoID : FlightInfoID:%d, err: %v\n", info.Id, err)
			}
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
			ticket, err := l.svcCtx.TicketsModel.FindOneBySpaceId(space.Id)
			if err != nil {
				if err == commonModel.ErrNotFound {
					logx.WithContext(l.ctx).Infof("NOT FOUND: There is no ticket information for the corresponding space.spaceID:%d\n", space.Id)
					return nil, errors.Wrapf(ERRGetTickets, "NOT FOUND: There is no ticket information for the corresponding space.spaceID:%d\n", space.Id)
				} else {
					return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.TicketsModel.FindListBySpaceID : spaceID:%d\n", space.Id)
				}
			}

			// 退改票信息
			refundInfo := &pb.RefundInfo{}
			changeInfo := &pb.ChangeInfo{}
			// 查退票信息
			ri, err := l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(ticket.Id, 1)
			if err != nil {
				if err == commonModel.ErrNotFound {
					logx.WithContext(l.ctx).Infof("NOT FOUND: There is no refund and change information for the corresponding ticket.ticketID:%d\n", ticket.Id)
					return nil, errors.Wrapf(ERRRefundAndChangeInfos, "NOT FOUND: There is no refund and change information for the corresponding ticket.ticketID:%d\n", ticket.Id)
				} else {
					return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.RefundAndChangeInfosModel.FindListByTicketID : ticketID:%d\n", ticket.Id)
				}
			}
			if ri != nil {
				refundInfo.TimeFees = append(refundInfo.TimeFees,
					&pb.TimeFee{Time: timestamppb.New(ri.Time1), Fee: uint64(ri.Fee1)},
					&pb.TimeFee{Time: timestamppb.New(ri.Time2), Fee: uint64(ri.Fee2)},
					&pb.TimeFee{Time: timestamppb.New(ri.Time3), Fee: uint64(ri.Fee3)},
					&pb.TimeFee{Time: timestamppb.New(ri.Time4), Fee: uint64(ri.Fee4)},
				)
				if !ri.Time5.IsZero() && ri.Fee5 > 0 {
					refundInfo.TimeFees = append(refundInfo.TimeFees, &pb.TimeFee{Time: timestamppb.New(ri.Time5), Fee: uint64(ri.Fee5)})
				}
			}
			// 查改票信息
			ci, err := l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(ticket.Id, 0)
			if err != nil {
				if err == commonModel.ErrNotFound {
					logx.WithContext(l.ctx).Infof("NOT FOUND: There is no refund and change information for the corresponding ticket.ticketID:%d\n", ticket.Id)
					return nil, errors.Wrapf(ERRRefundAndChangeInfos, "NOT FOUND: There is no refund and change information for the corresponding ticket.ticketID:%d\n", ticket.Id)
				} else {
					return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.RefundAndChangeInfosModel.FindListByTicketID : ticketID:%d\n", ticket.Id)
				}
			}
			if ci != nil {
				changeInfo.TimeFees = append(changeInfo.TimeFees,
					&pb.TimeFee{Time: timestamppb.New(ci.Time1), Fee: uint64(ci.Fee1)},
					&pb.TimeFee{Time: timestamppb.New(ci.Time2), Fee: uint64(ci.Fee2)},
					&pb.TimeFee{Time: timestamppb.New(ci.Time3), Fee: uint64(ci.Fee3)},
					&pb.TimeFee{Time: timestamppb.New(ci.Time4), Fee: uint64(ci.Fee4)},
				)
				if !ci.Time5.IsZero() && ci.Fee5 > 0 {
					changeInfo.TimeFees = append(changeInfo.TimeFees, &pb.TimeFee{Time: timestamppb.New(ci.Time5), Fee: uint64(ci.Fee5)})
				}
			}

			// 添加对应信息
			if respType {
				resp.(*pb.QuireBySetOutDateAndFlightNumberResp).FlightInfos = append(resp.(*pb.QuireBySetOutDateAndFlightNumberResp).FlightInfos, &pb.FlightInfo{
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
				return resp.(*pb.QuireBySetOutDateAndFlightNumberResp), nil
			} else {
				resp.(*pb.QuireBySetOutDateStartPositionEndPositionResp).FlightInfos = append(resp.(*pb.QuireBySetOutDateStartPositionEndPositionResp).FlightInfos, &pb.FlightInfo{
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
				return resp.(*pb.QuireBySetOutDateStartPositionEndPositionResp), nil
			}

		}

	}
	return nil, nil
}
