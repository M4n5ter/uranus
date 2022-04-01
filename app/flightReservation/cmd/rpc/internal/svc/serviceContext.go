package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/app/flightReservation/cmd/rpc/internal/config"
	"uranus/app/order/cmd/rpc/order"
	"uranus/commonModel"
)

type ServiceContext struct {
	Config                    config.Config
	Flights                   commonModel.FlightsModel
	FlightInfosModel          commonModel.FlightInfosModel
	SpacesModel               commonModel.SpacesModel
	TicketsModel              commonModel.TicketsModel
	RefundAndChangeInfosModel commonModel.RefundAndChangeInfosModel
	FlightInquiryRpcClient    flightinquiry.FlightInquiry
	OrderRpcClient            order.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:                    c,
		Flights:                   commonModel.NewFlightsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		FlightInfosModel:          commonModel.NewFlightInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SpacesModel:               commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		TicketsModel:              commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		RefundAndChangeInfosModel: commonModel.NewRefundAndChangeInfosModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		FlightInquiryRpcClient:    flightinquiry.NewFlightInquiry(zrpc.MustNewClient(c.FlightInquiryRpcConf)),
		OrderRpcClient:            order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
	}
}
