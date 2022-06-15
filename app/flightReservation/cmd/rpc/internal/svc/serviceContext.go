package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
	"uranus/app/flightReservation/cmd/rpc/internal/config"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/payment/cmd/rpc/payment"
	"uranus/commonModel"
)

type ServiceContext struct {
	Config                    config.Config
	Flights                   commonModel.FlightsModel
	FlightInfosModel          commonModel.FlightInfosModel
	SpacesModel               commonModel.SpacesModel
	TicketsModel              commonModel.TicketsModel
	RefundAndChangeInfosModel commonModel.RefundAndChangeInfosModel
	OrderRpcClient            order.Order
	PaymentRpcClient          payment.Payment
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                    c,
		Flights:                   commonModel.NewFlightsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		FlightInfosModel:          commonModel.NewFlightInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SpacesModel:               commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		TicketsModel:              commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		RefundAndChangeInfosModel: commonModel.NewRefundAndChangeInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		OrderRpcClient:            order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		PaymentRpcClient:          payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
	}
}

// GetFee 获取退改票手续费
func (s *ServiceContext) GetFee(rci *commonModel.RefundAndChangeInfos) (fee int64, ok bool) {
	now := time.Now()
	if now.Before(rci.Time1) {
		return rci.Fee1, true
	} else if now.Before(rci.Time2) {
		return rci.Fee2, true
	} else if now.Before(rci.Time3) {
		return rci.Fee3, true
	} else if now.Before(rci.Time4) {
		return rci.Fee4, true
	} else if now.Before(rci.Time5) {
		return rci.Fee5, true
	}
	return 99999999999, false
}

// ValidateOrderSn 检验订单是否在 orderList 中存在,存在返回 true
func (s *ServiceContext) ValidateOrderSn(orderSn string, orderList []*order.FlightOrder) (ticketID int64, ok bool) {

	for _, flightOrder := range orderList {
		if flightOrder.Sn == orderSn {
			return flightOrder.TicketId, true
		}
	}

	return -1, false
}
