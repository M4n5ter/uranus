package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/commonModel"

	"uranus/app/flightInquiry/cmd/rpc/internal/svc"
	"uranus/app/flightInquiry/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuireBySetOutDateStartPositionEndPositionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuireBySetOutDateStartPositionEndPositionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuireBySetOutDateStartPositionEndPositionLogic {
	return &QuireBySetOutDateStartPositionEndPositionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// QuireBySetOutDateStartPositionEndPosition 通过给定日期、出发地、目的地进行航班查询请求
func (l *QuireBySetOutDateStartPositionEndPositionLogic) QuireBySetOutDateStartPositionEndPosition(in *pb.QuireBySetOutDateStartPositionEndPositionReq) (*pb.QuireBySetOutDateStartPositionEndPositionResp, error) {
	//查询 FlightNumber SetOutDate Punctuality DepartPosition DepartTime ArrivePosition ArriveTime
	flightInfos, err := l.svcCtx.FlightInfosModel.FindListBySetOutDateAndPosition(l.svcCtx.FlightInfosModel.RowBuilder(), in.SetOutDate.AsTime(), in.DepartPosition, in.ArrivePosition)
	if err != nil {
		if err == commonModel.ErrNotFound {
			return nil, errors.Wrapf(ERRGetInfos, "NOT FOUND: can't found flight infos: SetOutTime->%v, DepartPosition->%s, ArrivePosition->%s, ERR: %v\n", in.SetOutDate.AsTime(), in.DepartPosition, in.ArrivePosition, err)
		} else {
			return nil, errors.Wrapf(ERRDBERR, "DBERR: when calling flightinquiry-rpc:l.svcCtx.FlightInfosModel.FindListByNumberAndSetOutDate : SetOutTime->%v, DepartPosition->%s, ArrivePosition->%s, ERR: %v\n", in.SetOutDate.AsTime(), in.DepartPosition, in.ArrivePosition, err)
		}
	}

	v, err := l.svcCtx.CombineAllInfos(flightInfos)
	if err != nil {
		return nil, err
	}

	return &pb.QuireBySetOutDateStartPositionEndPositionResp{FlightInfos: v}, err
}
