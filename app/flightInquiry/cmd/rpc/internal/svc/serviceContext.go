package svc

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/flightInquiry/bizcache"
	"uranus/app/flightInquiry/cmd/rpc/internal/config"
	"uranus/app/flightInquiry/cmd/rpc/pb"
	"uranus/common/xerr"
	"uranus/commonModel"
)

var ERRGetSpaces = xerr.NewErrMsg("暂无舱位信息")
var ERRGetTickets = xerr.NewErrMsg("暂无票信息")
var ERRRefundAndChangeInfos = xerr.NewErrMsg("暂无退改票信息")
var ERRDBERR = xerr.NewErrCode(xerr.DB_ERROR)

type ServiceContext struct {
	Config                    config.Config
	Flights                   commonModel.FlightsModel
	FlightInfosModel          commonModel.FlightInfosModel
	SpacesModel               commonModel.SpacesModel
	TicketsModel              commonModel.TicketsModel
	RefundAndChangeInfosModel commonModel.RefundAndChangeInfosModel
	Redis                     redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                    c,
		Flights:                   commonModel.NewFlightsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		FlightInfosModel:          commonModel.NewFlightInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SpacesModel:               commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		TicketsModel:              commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		RefundAndChangeInfosModel: commonModel.NewRefundAndChangeInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		Redis: *redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
	}
}

func (s *ServiceContext) CombineAllInfos(flightInfos []*commonModel.FlightInfos) ([]*pb.FlightInfo, error) {
	resp := make([]*pb.FlightInfo, 0)
	// 查询 IsFirstClass Surplus FlightTypes
	for _, info := range flightInfos {
		flt, err := s.Flights.FindOneByNumber(info.FlightNumber)
		if err != nil {
			if err == commonModel.ErrNotFound {
				flt = &commonModel.Flights{FltType: "unknown"}
			} else {
				return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.Flights.FindOneByNumber : FlightNumber:%s, err: %v\n", info.FlightNumber, err)
			}
		}

		var spaces []*commonModel.Spaces
		// 查 bizcache
		zset := fmt.Sprintf("fliID_%d", info.Id)
		spaceIdList, err := bizcache.ListAll(s.Redis, zset, bizcache.BizSpaceCachePrefix)
		// 查不到 bizcache 的情况
		if err != nil || len(spaceIdList) == 0 {

			// 此处查询不走缓存，直接打到DB上
			spaces, err = s.SpacesModel.FindListByFlightInfoID(s.SpacesModel.RowBuilder(), info.Id)
			if err != nil {
				if err == commonModel.ErrNotFound {
					return nil, errors.Wrapf(ERRGetSpaces, "NOT FOUND: There is no corresponding space information for this flightInfo.FlightInfoID:%d\n", info.Id)
				} else {
					return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.SpacesModel.FindListByFlightInfoID : FlightInfoID:%d, err: %v\n", info.Id, err)
				}
			}

			// 把 spaceID 列表加进 bizcache
			for _, space := range spaces {
				err := bizcache.AddID(s.Redis, space.Id, zset, bizcache.BizSpaceCachePrefix)
				if err != nil {
					logx.Errorf("ADD bizcache ERR: %v", err)
				}
			}
		} else {
			// 查到 bizcache 的情况
			spaces, err = s.GetSpacesByIdList(spaceIdList)
			if err != nil {
				return nil, err
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
			ticket, err := s.TicketsModel.FindOneBySpaceId(space.Id)
			if err != nil {
				if err == commonModel.ErrNotFound {
					return nil, errors.Wrapf(ERRGetTickets, "NOT FOUND: There is no ticket information for the corresponding space.spaceID:%d\n", space.Id)
				} else {
					return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.TicketsModel.FindListBySpaceID : spaceID:%d\n", space.Id)
				}
			}

			// 退改票信息
			refundInfo := &pb.RefundInfo{}
			changeInfo := &pb.ChangeInfo{}
			// 查退票信息
			ri, err := s.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(ticket.Id, 1)
			if err != nil {
				if err == commonModel.ErrNotFound {
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
			ci, err := s.RefundAndChangeInfosModel.FindOneByTicketIdIsRefund(ticket.Id, 0)
			if err != nil {
				if err == commonModel.ErrNotFound {
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

			resp = append(resp, &pb.FlightInfo{
				FlightInfoID:   info.Id,
				FlightNumber:   info.FlightNumber,
				SetOutDate:     timestamppb.New(info.SetOutDate),
				IsFirstClass:   ifc,
				SpaceID:        space.Id,
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

			//// 添加对应信息
			//if respType {
			//	resp.(*pb.QuireBySetOutDateAndFlightNumberResp).FlightInfos = append(resp.(*pb.QuireBySetOutDateAndFlightNumberResp).FlightInfos, &pb.FlightInfo{
			//		FlightNumber:   info.FlightNumber,
			//		SetOutDate:     timestamppb.New(info.SetOutDate),
			//		IsFirstClass:   ifc,
			//		Price:          uint64(ticket.Price),
			//		Discount:       ticket.Discount,
			//		Surplus:        space.Surplus,
			//		Punctuality:    uint32(info.Punctuality),
			//		DepartPosition: info.DepartPosition,
			//		DepartTime:     timestamppb.New(info.DepartTime),
			//		ArrivePosition: info.ArrivePosition,
			//		ArriveTime:     timestamppb.New(info.ArriveTime),
			//		RefundInfo:     refundInfo,
			//		ChangeInfo:     changeInfo,
			//		Cba:            ticket.Cba,
			//		FlightType:     flt.FltType,
			//	})
			//
			//} else {
			//	resp.(*pb.QuireBySetOutDateStartPositionEndPositionResp).FlightInfos = append(resp.(*pb.QuireBySetOutDateStartPositionEndPositionResp).FlightInfos, &pb.FlightInfo{
			//		FlightNumber:   info.FlightNumber,
			//		SetOutDate:     timestamppb.New(info.SetOutDate),
			//		IsFirstClass:   ifc,
			//		Price:          uint64(ticket.Price),
			//		Discount:       ticket.Discount,
			//		Surplus:        space.Surplus,
			//		Punctuality:    uint32(info.Punctuality),
			//		DepartPosition: info.DepartPosition,
			//		DepartTime:     timestamppb.New(info.DepartTime),
			//		ArrivePosition: info.ArrivePosition,
			//		ArriveTime:     timestamppb.New(info.ArriveTime),
			//		RefundInfo:     refundInfo,
			//		ChangeInfo:     changeInfo,
			//		Cba:            ticket.Cba,
			//		FlightType:     flt.FltType,
			//	})
			//
			//}

		}

	}
	return resp, nil
}

func (s *ServiceContext) GetUniqFlightWithSpacesFromCombinedFlightInfos(combinedFlightInfos []*pb.FlightInfo) []*pb.UniqFlightWithSpaces {
	uniqFlightWithSpaces := make([]*pb.UniqFlightWithSpaces, 0)
	fliIDMap := make(map[int64][]*pb.FlightInfo)
	for _, info := range combinedFlightInfos {
		if _, exist := fliIDMap[info.FlightInfoID]; exist {
			// 已经存在同 flightInfoID 的情况
			fliIDMap[info.FlightInfoID] = append(fliIDMap[info.FlightInfoID], info)
		} else {
			// 首次出现的 flightInfoID
			fliIDMap[info.FlightInfoID] = []*pb.FlightInfo{info}
		}
	}

	// 将 map 填充进 uniqFlightWithSpaces
	for flightInfoID, spacesOfFlightInfo := range fliIDMap {
		uniqFlightWithSpaces = append(uniqFlightWithSpaces, &pb.UniqFlightWithSpaces{
			FlightInfoID:       flightInfoID,
			SpacesOfFlightInfo: spacesOfFlightInfo,
		})
	}

	return uniqFlightWithSpaces
}

func (s *ServiceContext) GetFlightInfosByIdList(idList map[int64]struct{}) ([]*commonModel.FlightInfos, error) {
	var ret []*commonModel.FlightInfos
	for id := range idList {
		fli, err := s.FlightInfosModel.FindOne(id)
		if err != nil && err != commonModel.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
		}

		if fli == nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("未找到对应flightinfo"), "Not Found Err: flightInfo ID: %d", id)
		}

		ret = append(ret, fli)
	}

	return ret, nil
}

func (s *ServiceContext) GetSpacesByIdList(idList map[int64]struct{}) ([]*commonModel.Spaces, error) {
	var ret []*commonModel.Spaces
	for id := range idList {
		space, err := s.SpacesModel.FindOne(id)
		if err != nil && err != commonModel.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
		}

		if space == nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("未找到对应space"), "Not Found Err: flightInfo ID: %d", id)
		}

		ret = append(ret, space)
	}

	return ret, nil
}
