package logic

import (
	"context"
	"github.com/pkg/errors"
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

	var flightInfos []*commonModel.FlightInfos

	flightInfos, err := l.svcCtx.FlightInfosModel.FindTransferFlightsByPlace(l.svcCtx.FlightInfosModel.RowBuilder(), in.DepartPosition, in.ArrivePosition, in.SetOutDate.AsTime())
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR: %v", err)
	}

	if flightInfos == nil || len(flightInfos) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("未找到航班信息"), "Err Not Found: departPosition: %s, arrivePosition: %s, err: %v", in.DepartPosition, in.ArrivePosition, err)
	}

	combinedFLIs, err := l.svcCtx.CombineAllInfos(flightInfos)
	if err != nil {
		return nil, err
	}

	l.svcCtx.GetUniqFlightWithSpacesFromCombinedFlightInfos(combinedFLIs)

	return &pb.QuireTransferFlightsResp{}, nil
}
