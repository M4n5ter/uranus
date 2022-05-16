package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"uranus/app/flightInquiry/bizcache"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecommendReturnJourneyByFlightInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecommendReturnJourneyByFlightInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecommendReturnJourneyByFlightInfoLogic {
	return &RecommendReturnJourneyByFlightInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecommendReturnJourneyByFlightInfo 根据指定航班信息提供返程机票推荐(仅支持直飞)
func (l *RecommendReturnJourneyByFlightInfoLogic) RecommendReturnJourneyByFlightInfo(in *pb.RecommendReturnJourneyByFlightInfoReq) (*pb.RecommendReturnJourneyByFlightInfoResp, error) {
	if in.FlightInfo == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "")
	}

	var flightInfos []*commonModel.FlightInfos

	// 从 bizcache 查 id 列表(跟 GetDiscountFlights 共用 bizcache)
	zset := fmt.Sprintf("GetDiscountFlights-%s_%s", in.FlightInfo.ArrivePosition, in.FlightInfo.DepartPosition)
	idList, err := bizcache.ListAll(l.svcCtx.Redis, zset, bizcache.BizFLICachePrefix)

	// 查不到 bizcache 的情况
	if err != nil || len(idList) == 0 {
		flightInfos, err = l.svcCtx.FlightInfosModel.FindPageListByPositionAndDays(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightInfo.ArrivePosition, in.FlightInfo.DepartPosition, 7, 5)
		if err != nil && err != commonModel.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %+v", err)
		}

		if flightInfos == nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("未找到航班信息"), "Err Not Found: departPosition: %s, arrivePosiiton: %s, days: %d", in.FlightInfo.ArrivePosition, in.FlightInfo.DepartTime, 7)
		}

		// 把id列表加进bizcache
		for _, info := range flightInfos {
			err = bizcache.AddID(l.svcCtx.Redis, info.Id, zset, bizcache.BizFLICachePrefix)
			if err != nil {
				logx.Errorf("ADD bizcache ERR: %v", err)
			}
		}

		combinedFLIs, err := l.svcCtx.CombineAllInfos(flightInfos)
		if err != nil {
			return nil, err
		}

		ret := l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(combinedFLIs)
		return &pb.RecommendReturnJourneyByFlightInfoResp{UniqFlightWithSpaces: ret}, nil
	}

	// 找到 bizcache
	flightInfos, err = l.svcCtx.GetFlightInfosByIdList(idList)
	if err != nil {
		return nil, err
	}

	combinedFLIs, err := l.svcCtx.CombineAllInfos(flightInfos)
	if err != nil {
		return nil, err
	}

	ret := l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(combinedFLIs)
	return &pb.RecommendReturnJourneyByFlightInfoResp{UniqFlightWithSpaces: ret}, nil
}
