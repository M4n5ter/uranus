package svc

import (
	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/order/cmd/mq/internal/config"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/usercenter/cmd/rpc/usercenter"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderClient      order.Order
	MqueueClient     mqueue.Mqueue
	UserCenterClient usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderClient:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		MqueueClient:     mqueue.NewMqueue(zrpc.MustNewClient(c.MqueueRpcConf)),
		UserCenterClient: usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpcConf)),
	}
}
