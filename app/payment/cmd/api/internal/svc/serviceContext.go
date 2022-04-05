package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/order/cmd/rpc/order"
	"uranus/app/payment/cmd/api/internal/config"
	"uranus/app/payment/cmd/rpc/payment"
	"uranus/app/userCenter/cmd/rpc/usercenter"
	"uranus/commonModel"
)

type ServiceContext struct {
	Config           config.Config
	UserCenterClient usercenter.Usercenter
	PaymentClient    payment.Payment
	OrderClient      order.Order
	TicketsModel     commonModel.TicketsModel
	SpacesModel      commonModel.SpacesModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		UserCenterClient: usercenter.NewUsercenter(zrpc.MustNewClient(c.UserCenterRpcConf)),
		PaymentClient:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		OrderClient:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		TicketsModel:     commonModel.NewTicketsModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SpacesModel:      commonModel.NewSpacesModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}