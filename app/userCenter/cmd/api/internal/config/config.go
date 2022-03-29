package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/common/casbinTools"
)

type Config struct {
	rest.RestConf
	UserCenterRpc zrpc.RpcClientConf
	AuthRpc       zrpc.RpcClientConf
	WxMiniConf    struct {
		AppId  string
		Secret string
	}
	CasbinConf casbinTools.Conf
}
