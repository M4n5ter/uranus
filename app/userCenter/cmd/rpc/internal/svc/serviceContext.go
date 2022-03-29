package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"uranus/app/uranusAuth/cmd/rpc/auth"
	"uranus/app/usercenter/cmd/rpc/internal/config"
	"uranus/app/usercenter/model"
)

type ServiceContext struct {
	Config        config.Config
	AuthRpcClient auth.Auth
	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		AuthRpcClient: auth.NewAuth(zrpc.MustNewClient(c.AuthRpc)),
		UserModel:     model.NewUserModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UserAuthModel: model.NewUserAuthModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
