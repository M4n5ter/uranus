package deferMq

import (
	"context"
	"fmt"
	"log"
	"uranus/app/order/cmd/mq/internal/svc"
	"uranus/common/asynqmq"

	"github.com/hibiken/asynq"
)

/**
监听关闭订单
*/
type AsynqTask struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAsynqTask(ctx context.Context, svcCtx *svc.ServiceContext) *AsynqTask {
	return &AsynqTask{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AsynqTask) Start() {

	fmt.Println("AsynqTask start ")

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: l.svcCtx.Config.Redis.Host, Password: l.svcCtx.Config.Redis.Pass},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	mux := asynq.NewServeMux()

	//关闭航班订单任务
	mux.HandleFunc(asynqmq.TypeFlightOrderCloseDelivery, l.closeFlightOrderStateMqHandler)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func (l *AsynqTask) Stop() {
	fmt.Println("AsynqTask stop")
}
