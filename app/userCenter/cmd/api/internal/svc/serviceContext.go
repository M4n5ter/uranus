package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/userCenter/cmd/api/internal/config"
	"uranus/app/userCenter/cmd/rpc/usercenter"
)

type ServiceContext struct {
	Config              config.Config
	UsercenterRpcClient usercenter.Usercenter
	//CasbinCachedEnforcer casbin.CachedEnforcer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		UsercenterRpcClient: usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpc)),
		//CasbinCachedEnforcer: *casbinTools.MustGetEnforcer(c.CasbinConf),
	}
}
