package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlightDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlightDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlightDetailLogic {
	return &FlightDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FlightDetail 通过给定 票id 来获取航班详情
func (l *FlightDetailLogic) FlightDetail(in *pb.FlightDetailReq) (*pb.FlightDetailResp, error) {
	// 检查输入合法性
	if in.TicketID < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input: ticketID: %d", in.TicketID)
	}

	// 找到对应的票信息
	ticket, err := l.svcCtx.TicketsModel.FindOne(in.TicketID)
	if err != nil {
		if err == commonModel.ErrNotFound {
			return nil, errors.Wrapf(ERRGetTickets, "Not Found Ticket: ticketID: %d", in.TicketID)
		}
		return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc: l.svcCtx.TicketsModel.FindOne, ticketID: %d, err: %v", in.TicketID, err)
	}
	// 找到对应的舱位信息
	space, err := l.svcCtx.SpacesModel.FindOne(ticket.SpaceId)
	if err != nil {
		if err == commonModel.ErrNotFound {
			return nil, errors.Wrapf(ERRGetTickets, "Not Found Space: spaceID: %d", ticket.SpaceId)
		}
		return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc: l.svcCtx.SpacesModel.FindOne, spaceID: %d, err: %v", ticket.SpaceId, err)
	}
	// 找到对应的航班信息
	flightInfo, err := l.svcCtx.FlightInfosModel.FindOne(space.FlightInfoId)
	if err != nil {
		if err == commonModel.ErrNotFound {
			return nil, errors.Wrapf(ERRGetTickets, "Not Found FlightInfo: flightInfoID: %d", space.FlightInfoId)
		}
		return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc: l.svcCtx.FlightInfosModel.FindOne, flightInfoID: %d, err: %v", space.FlightInfoId, err)
	}
	// 找到机型
	flight, err := l.svcCtx.Flights.FindOneByNumber(flightInfo.FlightNumber)
	if err != nil {
		if err == commonModel.ErrNotFound {
			err = l.svcCtx.Flights.Trans(func(session sqlx.Session) error {
				_, err = l.svcCtx.Flights.Insert(session, &commonModel.Flights{
					DelState: 0,
					Version:  0,
					Number:   flightInfo.FlightNumber,
					FltType:  "unknown",
				})
				return err
			})
			if err != nil {
				return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc: l.svcCtx.Flights.Insert, flightNumber: %s, err: %v", flightInfo.FlightNumber, err)
			}
		}
		return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc: l.svcCtx.Flights.FindOneByNumber, flightNumber: %s, err: %v", flightInfo.FlightNumber, err)
	}
	// 获取退改票信息
	var ri, ci *commonModel.RefundAndChangeInfos
	err = mr.Finish(func() error {
		ri, err = l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(in.TicketID, 1)
		if err != nil {
			if err == commonModel.ErrNotFound {
			}
			return errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc: l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketId, ticketID: %d, err: %v", in.TicketID, err)
		}
		return nil
	}, func() error {
		ci, err = l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(in.TicketID, 0)
		if err != nil {
			if err == commonModel.ErrNotFound {
			}
			return errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc: l.svcCtx.RefundAndChangeInfosModel.FindOneByTicketId, ticketID: %d, err: %v", in.TicketID, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	resp := new(pb.FlightInfo)
	resp.FlightInfoID = flightInfo.Id
	resp.DepartPosition = flightInfo.DepartPosition
	resp.ArrivePosition = flightInfo.ArrivePosition
	resp.FlightNumber = flightInfo.FlightNumber
	resp.ArriveTime = timestamppb.New(flightInfo.ArriveTime)
	resp.DepartTime = timestamppb.New(flightInfo.DepartTime)
	resp.FlightType = flight.FltType
	resp.Surplus = space.Surplus
	resp.IsFirstClass = space.IsFirstClass > 1
	resp.SetOutDate = timestamppb.New(flightInfo.SetOutDate)
	resp.Discount = ticket.Discount
	resp.Punctuality = uint32(flightInfo.Punctuality)
	resp.Price = uint64(ticket.Price)
	resp.Cba = ticket.Cba
	if ri != nil {
		resp.RefundInfo = &pb.RefundInfo{}
		resp.RefundInfo.TimeFees = append(resp.RefundInfo.TimeFees,
			&pb.TimeFee{Time: timestamppb.New(ri.Time1), Fee: uint64(ri.Fee1)},
			&pb.TimeFee{Time: timestamppb.New(ri.Time2), Fee: uint64(ri.Fee2)},
			&pb.TimeFee{Time: timestamppb.New(ri.Time3), Fee: uint64(ri.Fee3)},
			&pb.TimeFee{Time: timestamppb.New(ri.Time4), Fee: uint64(ri.Fee4)},
		)
		if !ri.Time5.IsZero() && ri.Fee5 > 0 {
			resp.RefundInfo.TimeFees = append(resp.RefundInfo.TimeFees, &pb.TimeFee{Time: timestamppb.New(ri.Time5), Fee: uint64(ri.Fee5)})
		}
	}
	if ci != nil {
		resp.ChangeInfo = &pb.ChangeInfo{}
		resp.ChangeInfo.TimeFees = append(resp.ChangeInfo.TimeFees,
			&pb.TimeFee{Time: timestamppb.New(ci.Time1), Fee: uint64(ci.Fee1)},
			&pb.TimeFee{Time: timestamppb.New(ci.Time2), Fee: uint64(ci.Fee2)},
			&pb.TimeFee{Time: timestamppb.New(ci.Time3), Fee: uint64(ci.Fee3)},
			&pb.TimeFee{Time: timestamppb.New(ci.Time4), Fee: uint64(ci.Fee4)},
		)
		if !ci.Time5.IsZero() && ci.Fee5 > 0 {
			resp.ChangeInfo.TimeFees = append(resp.ChangeInfo.TimeFees, &pb.TimeFee{Time: timestamppb.New(ci.Time5), Fee: uint64(ci.Fee5)})
		}
	}

	return &pb.FlightDetailResp{FlightInfo: resp}, nil
}
