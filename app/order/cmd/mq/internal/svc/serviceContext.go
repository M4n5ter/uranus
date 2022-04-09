package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/order/cmd/mq/internal/config"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/usercenter/cmd/rpc/usercenter"
	"uranus/commonModel"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	SpacesModel      commonModel.SpacesModel
	TicketsModel     commonModel.TicketsModel
	OrderClient      order.Order
	MqueueClient     mqueue.Mqueue
	UserCenterClient usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		SpacesModel:      commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		TicketsModel:     commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		OrderClient:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		MqueueClient:     mqueue.NewMqueue(zrpc.MustNewClient(c.MqueueRpcConf)),
		UserCenterClient: usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpcConf)),
	}
}
