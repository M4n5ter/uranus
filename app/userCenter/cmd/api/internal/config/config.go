package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserCenterRpc zrpc.RpcClientConf
	JwtAuth       struct {
		AccessSecret string
	}
	WxMiniConf struct {
		AppId  string
		Secret string
	}
	//CasbinConf casbinTools.Conf
}
