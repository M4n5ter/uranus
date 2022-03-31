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

	OrderRpc      order.Order
	MqueueRpc     mqueue.Mqueue
	UserCenterRpc usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		MqueueRpc:     mqueue.NewMqueue(zrpc.MustNewClient(c.MqueueRpcConf)),
		UserCenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpcConf)),
	}
}
