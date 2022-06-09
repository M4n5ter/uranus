package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"uranus/app/flightInquiry/bizcache"
	"uranus/common/xerr"
	"uranus/commonModel"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/zeromicro/go-zero/core/logx"
	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"
)

type GetFlightsByPriceRangeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFlightsByPriceRangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFlightsByPriceRangeLogic {
	return &GetFlightsByPriceRangeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFlightsByPriceRange 查询指定价格区间的航班
func (l *GetFlightsByPriceRangeLogic) GetFlightsByPriceRange(in *pb.GetFlightsByPriceRangeReq) (*pb.GetFlightsByPriceRangeResp, error) {
	// 检查输入
	if len(in.ArrivePosition) == 0 || len(in.DepartPosition) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "")
	}

	var flightInfos []*commonModel.FlightInfos

	// 从 bizcache 查 id 列表
	zset := fmt.Sprintf("GetFlightsByPriceRange-%s_%s", in.DepartPosition, in.ArrivePosition)
	idList, err := bizcache.ListAll(l.svcCtx.Redis, zset, bizcache.BizFLICachePrefix)
	if err != nil || len(idList) == 0 {
		// 查不到 bizcache
		flightInfos, err := l.svcCtx.FlightInfosModel.FindPageListByPositionSODAndDays(l.svcCtx.FlightInfosModel.RowBuilder(), in.DepartPosition, in.ArrivePosition, in.SelectedDate.AsTime(), in.Days, in.Num)
		if err != nil && err != commonModel.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR : %+v", err)
		}

		if flightInfos == nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("未找到航班信息"), "Err Not Found: departPosition: %s, arrivePosition: %s, days: %d, num: %d", in.DepartPosition, in.ArrivePosition, in.Days, in.Num)
		}

		// 把 id 列表加进 bizcache
		for _, info := range flightInfos {
			err = bizcache.AddID(l.svcCtx.Redis, info.Id, zset, bizcache.BizFLICachePrefix)
			if err != nil {
				logx.Errorf("ADD bizcache ERR: %v", err)
			}
		}

		combinedFLIs, err := l.svcCtx.CombineAllInfos(flightInfos)
		filteredCombinedFLIs := slice.Filter(combinedFLIs, func(i int, v *pb.FlightInfo) bool {
			total := int64(v.Price) - int64(float64(v.Discount)/100*float64(v.Price))
			return total >= in.MinPrice && total <= in.MaxPrice
		})
		return &pb.GetFlightsByPriceRangeResp{UniqFlightWithSpaces: l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(filteredCombinedFLIs)}, nil
	}

	// 查到 bizcache
	flightInfos, err = l.svcCtx.GetFlightInfosByIdList(idList)
	if err != nil {
		return nil, err
	}

	combinedFLIs, err := l.svcCtx.CombineAllInfos(flightInfos)
	filteredCombinedFLIs := slice.Filter(combinedFLIs, func(i int, v *pb.FlightInfo) bool {
		total := int64(v.Price) - int64(float64(v.Discount)/100*float64(v.Price))
		return total >= in.MinPrice && total <= in.MaxPrice
	})
	return &pb.GetFlightsByPriceRangeResp{UniqFlightWithSpaces: l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(filteredCombinedFLIs)}, nil
}
