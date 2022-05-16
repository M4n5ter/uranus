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

type GetFlightsByFlightNumberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFlightsByFlightNumberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightsByFlightNumberLogic {
	return &GetFlightsByFlightNumberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFlightsByFlightNumber 根据航班号获取航班信息
func (l *GetFlightsByFlightNumberLogic) GetFlightsByFlightNumber(in *pb.GetFlightsByFlightNumberReq) (*pb.GetFlightsByFlightNumberResp, error) {

	if len(in.FlightNumber) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "")
	}

	var flightInfos []*commonModel.FlightInfos

	// 从 bizcache 查 id 列表
	zset := fmt.Sprintf("GetFlightsByFlightNumber-%s", in.FlightNumber)
	idList, err := bizcache.ListAll(l.svcCtx.Redis, zset, bizcache.BizFLICachePrefix)
	if err != nil || len(idList) == 0 {
		// 查不到 bizcache
		flightInfos, err := l.svcCtx.FlightInfosModel.FindPageListByNumberAndDays(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightNumber, in.Days, in.Num)
		if err != nil && err != commonModel.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR : %+v", err)
		}

		if flightInfos == nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("航班信息不存在"), "Err Not Found : flightNumber: %s", in.FlightNumber)
		}

		// 把 id 列表加进 bizcache
		for _, info := range flightInfos {
			err = bizcache.AddID(l.svcCtx.Redis, info.Id, zset, bizcache.BizFLICachePrefix)
			if err != nil {
				logx.Errorf("ADD bizcache ERR: %v", err)
			}
		}

		combinedResp, err := l.svcCtx.CombineAllInfos(flightInfos)
		if err != nil {
			return nil, err
		}

		uniqFlightWithSpaces := l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(combinedResp)

		return &pb.GetFlightsByFlightNumberResp{UniqFlightWithSpaces: uniqFlightWithSpaces}, nil
	}

	// 查到 bizcache
	flightInfos, err = l.svcCtx.GetFlightInfosByIdList(idList)
	if err != nil {
		return nil, err
	}

	combinedResp, err := l.svcCtx.CombineAllInfos(flightInfos)
	if err != nil {
		return nil, err
	}

	uniqFlightWithSpaces := l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(combinedResp)

	return &pb.GetFlightsByFlightNumberResp{UniqFlightWithSpaces: uniqFlightWithSpaces}, nil
}
