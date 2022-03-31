package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/flightInquiry/cmd/rpc/flightinquiry"
	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/order/cmd/rpc/internal/config"
	"uranus/app/order/model"
)

type ServiceContext struct {
	Config              config.Config
	FlightInquiryClient flightinquiry.FlightInquiry
	MqueueClient        mqueue.Mqueue
	OrderModel          model.FlightOrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		FlightInquiryClient: flightinquiry.NewFlightInquiry(zrpc.MustNewClient(c.FlightInquiryConf)),
		MqueueClient:        mqueue.NewMqueue(zrpc.MustNewClient(c.MqueueRpcConf)),
		OrderModel:          model.NewFlightOrderModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
