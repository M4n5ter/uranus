package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/mqueue/cmd/rpc/mqueue"
	"uranus/app/payment/cmd/rpc/internal/config"
	"uranus/app/payment/model"
)

type ServiceContext struct {
	Config       config.Config
	PaymentModel model.PaymentModel
	MqueueClient mqueue.Mqueue
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		PaymentModel: model.NewPaymentModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		MqueueClient: mqueue.NewMqueue(zrpc.MustNewClient(c.MqueueRpcConf)),
	}
}
