package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/runtime/protoimpl"
	"uranus/app/flightInquiry/bizcache"
	"uranus/common/timeTools"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ERRGetInfos = xerr.NewErrMsg("暂无直飞航班")
var ERRGetTickets = xerr.NewErrMsg("暂无票信息")
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
	var flightInfos []*commonModel.FlightInfos

	// 从 bizcache 查 id 列表
	zset := fmt.Sprintf("%s_%s", in.FlightNumber, timeTools.Timestamppb2TimeStringYMD(in.SetOutDate))
	idList, err := bizcache.ListAll(l.svcCtx.Redis, zset, bizcache.BizFLICachePrefix)
	// 查不到bizcache的情况
	if err != nil || idList == nil {

		// 不走缓存查询 FlightNumber SetOutDate Punctuality DepartPosition DepartTime ArrivePosition ArriveTime
		flightInfos, err = l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightNumber, in.SetOutDate.AsTime())
		if err != nil {
			if err == commonModel.ErrNotFound {
				return nil, errors.Wrapf(ERRGetInfos, "NOT FOUND: can't found flight infos: number->%s setOutTime->%v, ERR: %v\n", in.FlightNumber, in.SetOutDate.AsTime(), err)
			} else {
				return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate : number->%s setOutTime->%v, ERR: %v\n", in.FlightNumber, in.SetOutDate.AsTime(), err)
			}
		}

		// 把id列表加进bizcache
		for _, info := range flightInfos {
			err = bizcache.AddID(l.svcCtx.Redis, info.Id, zset, bizcache.BizFLICachePrefix)
			if err != nil {
				logx.Errorf("ADD bizcache ERR: %v", err)
			}
		}

		// 组合数据并返回
		v, err := l.svcCtx.CombineAllInfos(flightInfos)
		if err != nil {
			return nil, err
		}

		return &pb.QuireBySetOutDateAndFlightNumberResp{FlightInfos: v}, err
	}

	// 查到bizcache的情况
	flightInfos, err = l.svcCtx.GetFlightInfosByIdList(idList)
	if err != nil {
		return nil, err
	}

	// 组合数据并返回
	v, err := l.svcCtx.CombineAllInfos(flightInfos)
	if err != nil {
		return nil, err
	}

	return &pb.QuireBySetOutDateAndFlightNumberResp{FlightInfos: v}, err

	//// 不走缓存查询 FlightNumber SetOutDate Punctuality DepartPosition DepartTime ArrivePosition ArriveTime
	//flightInfos, err = l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightNumber, in.SetOutDate.AsTime())
	//if err != nil {
	//	if err == commonModel.ErrNotFound {
	//		return nil, errors.Wrapf(ERRGetInfos, "NOT FOUND: can't found flight infos: number->%s setOutTime->%v, ERR: %v\n", in.FlightNumber, in.SetOutDate.AsTime(), err)
	//	} else {
	//		return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate : number->%s setOutTime->%v, ERR: %v\n", in.FlightNumber, in.SetOutDate.AsTime(), err)
	//	}
	//}
	//
	//v, err := l.svcCtx.CombineAllInfos(flightInfos)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &pb.QuireBySetOutDateAndFlightNumberResp{FlightInfos: v}, err
}
