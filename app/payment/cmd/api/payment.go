package main

import (
	"flag"
	"fmt"
	"uranus/common/middleware"

	"uranus/app/payment/cmd/api/internal/config"
	"uranus/app/payment/cmd/api/internal/handler"
	"uranus/app/payment/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/payment.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	// 全局中间件
	//将nginx网关验证后的userId设置到ctx中.
	server.Use(middleware.NewSetUidToCtxMiddleware().Handle)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
