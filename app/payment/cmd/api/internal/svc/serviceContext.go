package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/payment/cmd/api/internal/config"
	"uranus/app/payment/cmd/rpc/payment"
	"uranus/app/userCenter/cmd/rpc/userCenter"
	"uranus/commonModel"
)

type ServiceContext struct {
	Config           config.Config
	UserCenterClient userCenter.Usercenter
	PaymentClient    payment.Payment
	OrderClient      order.Order
	TicketsModel     commonModel.TicketsModel
	SpacesModel      commonModel.SpacesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		UserCenterClient: userCenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpcConf)),
		PaymentClient:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		OrderClient:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		TicketsModel:     commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SpacesModel:      commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
