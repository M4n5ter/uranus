package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"uranus/app/userCenter/cmd/rpc/internal/config"
	"uranus/app/userCenter/model"
)

type ServiceContext struct {
	Config        config.Config
	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
	WalletModel   model.WalletModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UserModel:     model.NewUserModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UserAuthModel: model.NewUserAuthModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		WalletModel:   model.NewWalletModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
