package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	QiniuOSS struct {
		AccessKey string
		SecretKey string
		Bucket    string
		Domain    string
	}

	DB struct {
		DataSource string
	}

	Cache cache.CacheConf

	AuthRpc zrpc.RpcClientConf
}
