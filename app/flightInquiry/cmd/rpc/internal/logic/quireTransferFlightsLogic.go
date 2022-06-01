package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"uranus/app/flightInquiry/bizcache"
	"uranus/common/timeTools"
	"uranus/common/xerr"
	"uranus/commonModel"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireTransferFlightsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuireTransferFlightsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireTransferFlightsLogic {
	return &QuireTransferFlightsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuireTransferFlights 中转航班查询
func (l *QuireTransferFlightsLogic) QuireTransferFlights(in *pb.QuireTransferFlightsReq) (*pb.QuireTransferFlightsResp, error) {
	if len(in.DepartPosition) == 0 || len(in.ArrivePosition) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "")
	}

	var transfers []*commonModel.Transfer
	var combinedTransfers []*pb.TransferFlightInfo

	// 从 bizcache 查 id 列表
	zset := fmt.Sprintf("QuireTransferFlights-%s_%s_%s", in.DepartPosition, in.ArrivePosition, timeTools.Timestamppb2TimeStringYMD(in.SetOutDate))
	transfers, err := bizcache.ListAllTransfers(l.svcCtx.Redis, zset, bizcache.BizFLICachePrefix)

	// 查不到 bizcache
	if len(transfers) == 0 || err != nil {
		transfers, err = l.svcCtx.FlightInfosModel.FindTransferFlightsByPlace(l.svcCtx.FlightInfosModel.RowBuilder(), in.DepartPosition, in.ArrivePosition, in.SetOutDate.AsTime())
		if err != nil && err != commonModel.ErrNotFound {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
		}

		if transfers == nil || len(transfers) == 0 {
			return nil, errors.Wrapf(xerr.NewErrMsg("未找到航班信息"), "Err Not Found: departPosition: %s, arrivePosition: %s, err: %v", in.DepartPosition, in.ArrivePosition, err)
		}

		for _, transfer := range transfers {
			err := bizcache.AddTransfer(l.svcCtx.Redis, *transfer, zset, bizcache.BizTransferCachePrefix)
			if err != nil {
				logx.Errorf("ADD bizcache ERR: %v", err)
			}
		}

		for _, transfer := range transfers {
			combinedFLIs, err := l.svcCtx.CombineAllInfos(transfer.F)
			if err != nil {
				return nil, err
			}

			combinedTransfers = append(combinedTransfers, &pb.TransferFlightInfo{UniqFlightWithSpaces: l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(combinedFLIs)})
		}
		return &pb.QuireTransferFlightsResp{TransferFlights: combinedTransfers}, nil
	}

	// 查到 bizcache
	for _, transfer := range transfers {
		combinedFLIs, err := l.svcCtx.CombineAllInfos(transfer.F)
		if err != nil {
			return nil, err
		}

		combinedTransfers = append(combinedTransfers, &pb.TransferFlightInfo{UniqFlightWithSpaces: l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(combinedFLIs)})
	}

	return &pb.QuireTransferFlightsResp{TransferFlights: combinedTransfers}, nil
}
