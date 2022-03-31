package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/order/model"
	"uranus/common/xerr"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlightOrderDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlightOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlightOrderDetailLogic {
	return &FlightOrderDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FlightOrderDetail 机票订单详情
func (l *FlightOrderDetailLogic) FlightOrderDetail(in *pb.FlightOrderDetailReq) (*pb.FlightOrderDetailResp, error) {
	if len(in.Sn) == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("订单号不能为空"), "ERR: empty Sn!")
	}

	flightOrder, err := l.svcCtx.OrderModel.FindOneBySn(in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "flightOrderModel  FindOne err : %v , sn : %s", err, in.Sn)
	}

	var resp pb.FlightOrder
	if flightOrder != nil {
		_ = copier.Copy(&resp, flightOrder)
		resp.CreateTime = timestamppb.New(flightOrder.CreateTime)
		resp.DepartTime = timestamppb.New(flightOrder.DepartTime)
		resp.ArriveTime = timestamppb.New(flightOrder.ArriveTime)
	}

	return &pb.FlightOrderDetailResp{FlightOrder: &resp}, nil
}
