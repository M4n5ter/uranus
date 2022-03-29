package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"uranus/common/globalkey"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userAuthFieldNames          = builder.RawFieldNames(&UserAuth{})
	userAuthRows                = strings.Join(userAuthFieldNames, ",")
	userAuthRowsWithPlaceHolder = strings.Join(stringx.Remove(userAuthFieldNames, "`id`", "`create_time`", "`update_time`", "`version`"), "=?,") + "=?"

	cacheUserAuthIdPrefix              = "cache:userAuth:id:"
	cacheUserAuthAuthTypeAuthKeyPrefix = "cache:userAuth:authType:authKey:"
	cacheUserAuthUserIdAuthTypePrefix  = "cache:userAuth:userId:authType:"
)

type (
	UserAuthModel interface {
		FindOne(id int64) (*UserAuth, error)
		FindOneByAuthTypeAuthKey(authType string, authKey string) (*UserAuth, error)
		FindOneByUserIdAuthType(userId int64, authType string) (*UserAuth, error)
		Insert(session sqlx.Session, data *UserAuth) (sql.Result, error)
		Update(session sqlx.Session, data *UserAuth) error
		Delete(session sqlx.Session, data *UserAuth) error
		Trans(fn func(session sqlx.Session) error) error
	}

	defaultUserAuthModel struct {
		sqlc.CachedConn
		table string
	}

	UserAuth struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		UserId     int64     `db:"user_id"`
		AuthKey    string    `db:"auth_key"`  // 平台唯一id
		AuthType   string    `db:"auth_type"` // 平台类型
	}
)

func NewUserAuthModel(conn sqlx.SqlConn, c cache.CacheConf) UserAuthModel {
	return &defaultUserAuthModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_auth`",
	}
}

func (m *defaultUserAuthModel) Insert(session sqlx.Session, data *UserAuth) (sql.Result, error) {
	userAuthIdKey := fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, data.Id)
	userAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	userAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {

		//@todo self edit  value , because change table field is trouble in here , so self fix field is easy
		query := fmt.Sprintf("insert into .... (%s) values ...", m.table)
		if session != nil {
			//@todo self edit  value , because change table field is trouble in here , so self fix field is easy
			return session.Exec(query, data.DeleteTime, data.DelState, data.UserId, data.AuthKey, data.AuthType)
		}
		//@todo self edit  value , because change table field is trouble in here , so self fix field is easy
		return conn.Exec(query, data.DeleteTime, data.DelState, data.UserId, data.AuthKey, data.AuthType)
	}, userAuthUserIdAuthTypeKey, userAuthIdKey, userAuthAuthTypeAuthKeyKey)

}

func (m *defaultUserAuthModel) FindOne(id int64) (*UserAuth, error) {
	userAuthIdKey := fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, id)
	var resp UserAuth
	err := m.QueryRow(&resp, userAuthIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		if resp.DelState == globalkey.DelStateYes {
			return nil, ErrNotFound
		}
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindOneByAuthTypeAuthKey(authType string, authKey string) (*UserAuth, error) {
	userAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheUserAuthAuthTypeAuthKeyPrefix, authType, authKey)
	var resp UserAuth
	err := m.QueryRowIndex(&resp, userAuthAuthTypeAuthKeyKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `auth_type` = ? and `auth_key` = ? limit 1", userAuthRows, m.table)
		if err := conn.QueryRow(&resp, query, authType, authKey); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		if resp.DelState == globalkey.DelStateYes {
			return nil, ErrNotFound
		}
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindOneByUserIdAuthType(userId int64, authType string) (*UserAuth, error) {
	userAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheUserAuthUserIdAuthTypePrefix, userId, authType)
	var resp UserAuth
	err := m.QueryRowIndex(&resp, userAuthUserIdAuthTypeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and `auth_type` = ? limit 1", userAuthRows, m.table)
		if err := conn.QueryRow(&resp, query, userId, authType); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		if resp.DelState == globalkey.DelStateYes {
			return nil, ErrNotFound
		}
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) Update(session sqlx.Session, data *UserAuth) error {
	userAuthIdKey := fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, data.Id)
	userAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	userAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userAuthRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.UserId, data.AuthKey, data.AuthType, data.Id)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.UserId, data.AuthKey, data.AuthType, data.Id)
	}, userAuthAuthTypeAuthKeyKey, userAuthUserIdAuthTypeKey, userAuthIdKey)
	return err
}

func (m *defaultUserAuthModel) Delete(session sqlx.Session, data *UserAuth) error {
	data.DelState = globalkey.DelStateYes
	return m.Update(session, data)
}

func (m *defaultUserAuthModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}

func (m *defaultUserAuthModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserAuthIdPrefix, primary)
}

func (m *defaultUserAuthModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthRows, m.table)
	return conn.QueryRow(v, query, primary)
}
