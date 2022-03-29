package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/uranusAuth/cmd/api/internal/config"
	"uranus/app/uranusAuth/cmd/rpc/auth"
)

type ServiceContext struct {
	Config        config.Config
	AuthRpcClient auth.Auth
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		AuthRpcClient: auth.NewAuth(zrpc.MustNewClient(c.RpcConf)),
	}
}
