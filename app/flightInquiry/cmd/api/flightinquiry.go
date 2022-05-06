package main

import (
	"flag"
	"fmt"
	_ "go.uber.org/automaxprocs"
	"uranus/app/flightInquiry/cmd/api/internal/config"
	"uranus/app/flightInquiry/cmd/api/internal/handler"
	"uranus/app/flightInquiry/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/flightinquiry.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
