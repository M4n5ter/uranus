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

	flightInfos, err := l.svcCtx.FlightInfosModel.FindPageListByNumberAndDays(l.svcCtx.FlightInfosModel.RowBuilder(), in.FlightNumber, in.Days, in.Num)
	if err != nil && err != commonModel.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "DBERR : %+v", err)
	}

	if flightInfos == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("航班信息不存在"), "Err Not Found : flightNumber: %s", in.FlightNumber)
	}

	ret, err := l.svcCtx.CombineAllInfos(flightInfos)
	if err != nil {
		return nil, err
	}
	return &pb.GetFlightsByFlightNumberResp{FlightInfos: ret}, nil
}
