package svc

import (
	"github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/uranusAuth/cmd/rpc/auth"
	"uranus/app/userCenter/cmd/api/internal/config"
	"uranus/app/userCenter/cmd/rpc/usercenter"
	"uranus/common/casbinTools"
)

type ServiceContext struct {
	Config               config.Config
	UsercenterRpcClient  usercenter.Usercenter
	AuthRpcClient        auth.Auth
	CasbinCachedEnforcer casbin.CachedEnforcer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:               c,
		UsercenterRpcClient:  usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpc)),
		AuthRpcClient:        auth.NewAuth(zrpc.MustNewClient(c.AuthRpc)),
		CasbinCachedEnforcer: *casbinTools.MustGetEnforcer(c.CasbinConf),
	}
}
