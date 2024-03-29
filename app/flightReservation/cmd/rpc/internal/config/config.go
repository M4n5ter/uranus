package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		DataSource string
	}
	Cache                cache.CacheConf
	FlightInquiryRpcConf zrpc.RpcClientConf
	OrderRpcConf         zrpc.RpcClientConf
	PaymentRpcConf       zrpc.RpcClientConf
	StockRpcConf         zrpc.RpcClientConf
	UserCenterRpcConf    zrpc.RpcClientConf
	DtmServer            struct {
		Target string
	}
}
