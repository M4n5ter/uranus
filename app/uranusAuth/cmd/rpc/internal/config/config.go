package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Jwt struct {
		AccessSecret string
		AccessExpire int64
	}
	logx.LogConf
}
