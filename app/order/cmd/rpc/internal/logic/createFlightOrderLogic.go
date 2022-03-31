package logic

import (
	"context"
	"github.com/pkg/errors"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/order/model"
	"uranus/common/tool"
	"uranus/common/uniqueid"
	"uranus/common/xerr"

	"uranus/app/order/cmd/rpc/internal/svc"
	"uranus/app/order/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFlightOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFlightOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFlightOrderLogic {
	return &CreateFlightOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 机票下订单
func (l *CreateFlightOrderLogic) CreateFlightOrder(in *pb.CreateFlightOrderReq) (*pb.CreateFlightOrderResp, error) {
	if in.UserId < 1 || in.TicketId < 1 {
		return nil, errors.Wrapf(xerr.NewErrMsg("非法输入"), "invalid input: userID: %d, ticketID: %d", in.UserId, in.TicketId)
	}
	// 获取机票详情
	flightDetailResp, err := l.svcCtx.FlightInquiryClient.FlightDetail(l.ctx, &flightinquiry.FlightDetailReq{TicketID: in.TicketId})
	if err != nil {
		return nil, err
	}
	if flightDetailResp == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("查询不到票所对应的航班信息"), "查询不到票所对应的航班信息, ticketID: %d", in.TicketId)
	}
	// 创建订单
	order := new(model.FlightOrder)
	order.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_FLIGHT_ORDER)
	order.OrderTotalPrice = int64(flightDetailResp.FlightInfo.Price) - int64(float64(flightDetailResp.FlightInfo.Discount)/100*float64(flightDetailResp.FlightInfo.Price))
	order.DepartPosition = flightDetailResp.FlightInfo.DepartPosition
	order.DepartTime = flightDetailResp.FlightInfo.DepartTime.AsTime()
	order.ArrivePosition = flightDetailResp.FlightInfo.ArrivePosition
	order.ArriveTime = flightDetailResp.FlightInfo.ArriveTime.AsTime()
	order.Discount = flightDetailResp.FlightInfo.Discount
	order.TicketId = in.TicketId
	order.UserId = in.UserId
	order.TradeCode = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	order.TicketPrice = int64(flightDetailResp.FlightInfo.Price)
	order.TradeState = model.FlightOrderTradeStateWaitPay
	_, err = l.svcCtx.OrderModel.Insert(nil, order)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "下单数据库异常 order : %+v , err: %v", order, err)
	}

	// 延迟关闭订单任务
	_, _ = l.svcCtx.MqueueClient.AqDeferFlightOrderClose(l.ctx, &mqueue.AqDeferFlightOrderCloseReq{Sn: order.Sn})
	return &pb.CreateFlightOrderResp{Sn: order.Sn}, nil
}
