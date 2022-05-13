package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"uranus/app/flightInquiry/bizcache"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDiscountFlightsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type CanBeSortedFLI []*pb.FlightInfo

func NewGetDiscountFlightsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDiscountFlightsLogic {
	return &GetDiscountFlightsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDiscountFlights 查询折扣航班
func (l *GetDiscountFlightsLogic) GetDiscountFlights(in *pb.GetDiscountFlightsReq) (*pb.GetDiscountFlightsResp, error) {
	// 检查输入
	if len(in.ArrivePosition) == 0 || len(in.DepartPosition) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "")
	}

	var flightInfos []*commonModel.FlightInfos

	// 从 bizcache 查 id 列表
	zset := fmt.Sprintf("%s_%s", in.DepartPosition, in.ArrivePosition)
	idList, err := bizcache.ListAll(l.svcCtx.Redis, zset, bizcache.BizFLICachePrefix)

	// 查不到 bizcache 的情况
	if err != nil || len(idList) == 0 {
		flightInfos, err = l.svcCtx.FlightInfosModel.FindPageListByPositionAndDays(l.svcCtx.FlightInfosModel.RowBuilder(), in.DepartPosition, in.ArrivePosition, in.Days, -1)
		if err != nil {
			if err == commonModel.ErrNotFound {
				return nil, errors.Wrapf(ERRGetInfos, "NOT FOUND: can't found flight infos: departPosition->%s arrivePosition->%v, ERR: %v\n", in.DepartPosition, in.ArrivePosition, err)
			} else {
				return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.FlightInfosModel.FindPageListByPositionAndDays : departPosition->%s arrivePosition->%v, ERR: %v\n", in.DepartPosition, in.ArrivePosition, err)
			}
		}

		// 把id列表加进bizcache
		for _, info := range flightInfos {
			err = bizcache.AddID(l.svcCtx.Redis, info.Id, zset, bizcache.BizFLICachePrefix)
			if err != nil {
				logx.Errorf("ADD bizcache ERR: %v", err)
			}
		}

		// 组合成完整数据组
		combinedInfos, err := l.svcCtx.CombineAllInfos(flightInfos)
		if err != nil {
			return nil, err
		}

		ret := l.getDiscountFlights(combinedInfos, in.Num)
		return &pb.GetDiscountFlightsResp{FlightInfos: ret}, nil
	}
	// 查到 bizcache
	flightInfos, err = l.svcCtx.GetFlightInfosByIdList(idList)
	if err != nil {
		return nil, err
	}

	// 组合成完整数据组
	combinedInfos, err := l.svcCtx.CombineAllInfos(flightInfos)
	if err != nil {
		return nil, err
	}

	ret := l.getDiscountFlights(combinedInfos, in.Num)
	return &pb.GetDiscountFlightsResp{FlightInfos: ret}, nil
}

// 获取折扣最大的 n 条航班信息
func (l *GetDiscountFlightsLogic) getDiscountFlights(flightList []*pb.FlightInfo, n int64) []*pb.FlightInfo {
	// 降序排序
	sort.Sort(CanBeSortedFLI(flightList))
	if n <= 0 {
		return flightList
	}

	if len(flightList) < int(n) {
		return flightList
	}

	return flightList[:n]
}

func (c CanBeSortedFLI) Len() int {
	return len(c)
}

// Less 此处用 > 来降序
func (c CanBeSortedFLI) Less(i, j int) bool {
	return c[i].Discount > c[j].Discount
}

func (c CanBeSortedFLI) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
