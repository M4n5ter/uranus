package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/order/cmd/rpc/internal/config"
	"uranus/app/order/model"
	"uranus/commonModel"
)

type ServiceContext struct {
	Config              config.Config
	FlightInquiryClient flightinquiry.FlightInquiry
	MqueueClient        mqueue.Mqueue
	OrderModel          model.FlightOrderModel
	SpacesModel         commonModel.SpacesModel
	TicketsModel        commonModel.TicketsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		FlightInquiryClient: flightinquiry.NewFlightInquiry(zrpc.MustNewClient(c.FlightInquiryConf)),
		MqueueClient:        mqueue.NewMqueue(zrpc.MustNewClient(c.MqueueRpcConf)),
		OrderModel:          model.NewFlightOrderModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SpacesModel:         commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		TicketsModel:        commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
