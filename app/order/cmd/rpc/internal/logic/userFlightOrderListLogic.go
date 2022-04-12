package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"uranus/app/order/model"
	"uranus/common/xerr"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserFlightOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserFlightOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserFlightOrderListLogic {
	return &UserFlightOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserFlightOrderList 用户机票订单
func (l *UserFlightOrderListLogic) UserFlightOrderList(in *pb.UserFlightOrderListReq) (*pb.UserFlightOrderListResp, error) {

	whereBuilder := l.svcCtx.OrderModel.RowBuilder().Where(squirrel.Eq{"user_id": in.UserId})
	// 有支持的状态则筛选，没有支持的状态则查所有状态
	if in.TraderState >= model.FlightOrderTradeStateCancel && in.TraderState <= model.FlightOrderTradeStateExpire {
		whereBuilder = whereBuilder.Where(squirrel.Eq{"trade_state": in.TraderState})
	}

	list, err := l.svcCtx.OrderModel.FindPageListByIdDESC(whereBuilder, in.LastId, in.PageSize)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(ERRDBERR, "获取用户机票订单失败, err: %v, in: %+v", err, in)
	}
	if list == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("没有对应状态的订单"), "Not Found err: %v", err)
	}

	var resp []*pb.FlightOrder
	if len(list) > 0 {
		for _, order := range list {
			fltOrder := pb.FlightOrder{}
			_ = copier.Copy(&fltOrder, order)
			fltOrder.DepartTime = timestamppb.New(order.DepartTime)
			fltOrder.ArriveTime = timestamppb.New(order.ArriveTime)
			fltOrder.CreateTime = timestamppb.New(order.CreateTime)
			resp = append(resp, &fltOrder)
		}
	}
	return &pb.UserFlightOrderListResp{List: resp}, nil
}
