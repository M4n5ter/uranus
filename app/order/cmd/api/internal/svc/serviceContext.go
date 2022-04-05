package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/order/cmd/api/internal/config"
	"uranus/app/order/cmd/rpc/order"
)

type ServiceContext struct {
	Config         config.Config
	OrderRpcClient order.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		OrderRpcClient: order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
	}
}
