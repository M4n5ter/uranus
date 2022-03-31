package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	FlightInquiryConf zrpc.RpcClientConf
	MqueueRpcConf     zrpc.RpcClientConf
	DB                struct {
		DataSource string
	}
	Cache cache.CacheConf
}
