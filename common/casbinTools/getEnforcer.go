package casbinTools

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
	"log"
	"uranus/common/xerr"
)

var ErrCasbinErr = xerr.NewErrMsg("连接 Casbin 失败")

func GetEnforcer(conf Conf) (enforcer *casbin.CachedEnforcer, err error) {
	a, err := gormadapter.NewAdapter("mysql", conf.DB.DataSourceWithoutDBName, conf.DB.DBName)
	if err != nil {
		return nil, errors.Wrapf(ErrCasbinErr, "create adapter err: %v", err)
	}
	m, err := model.NewModelFromString(conf.Model)
	if err != nil {
		return nil, errors.Wrapf(ErrCasbinErr, "create casbinTools commonModel err: %v", err)
	}
	enforcer, err = casbin.NewCachedEnforcer(m, a)
	if err != nil {
		return nil, errors.Wrapf(ErrCasbinErr, "create cachedEnforcer err: %v", err)
	}
	return enforcer, nil
}

func MustGetEnforcer(conf Conf) *casbin.CachedEnforcer {
	e, err := GetEnforcer(conf)
	if err != nil {
		log.Fatalln(err)
	}
	return e
}
