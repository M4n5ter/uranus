package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	WxMiniConf        WxMiniConf
	WxPayConf         WxPayConf
	UserCenterRpcConf zrpc.RpcClientConf
	PaymentRpcConf    zrpc.RpcClientConf
	OrderRpcConf      zrpc.RpcClientConf
	StockRpcConf      zrpc.RpcClientConf

	DB struct {
		DataSource string
	}
	Cache cache.CacheConf

	DtmServer struct {
		Target string
	}
}
