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

	combinedResp, err := l.svcCtx.CombineAllInfos(flightInfos)
	if err != nil {
		return nil, err
	}

	uniqFlightWithSpaces := make([]*pb.UniqFlightWithSpaces, 0)
	fliIDMap := make(map[int64][]*pb.FlightInfo)
	for _, info := range combinedResp {
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

	return &pb.GetFlightsByFlightNumberResp{UniqFlightWithSpaces: uniqFlightWithSpaces}, nil
}
