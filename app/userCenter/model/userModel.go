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
	userFieldNames          = builder.RawFieldNames(&User{})
	userRows                = strings.Join(userFieldNames, ",")
	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_time`", "`update_time`", "`version`"), "=?,") + "=?"

	cacheUserIdPrefix     = "cache:user:id:"
	cacheUserMobilePrefix = "cache:user:mobile:"
)

type (
	UserModel interface {
		FindOne(id int64) (*User, error)
		FindOneByMobile(mobile string) (*User, error)
		Insert(session sqlx.Session, data *User) (sql.Result, error)
		Update(session sqlx.Session, data *User) error
		Delete(session sqlx.Session, data *User) error
		Trans(fn func(session sqlx.Session) error) error
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		Mobile     string    `db:"mobile"`
		Password   string    `db:"password"`
		Nickname   string    `db:"nickname"`
		Sex        int64     `db:"sex"` // 性别 0:男 1:女
		Avatar     string    `db:"avatar"`
		Info       string    `db:"info"`
	}
)

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf) UserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user`",
	}
}

func (m *defaultUserModel) Insert(session sqlx.Session, data *User) (sql.Result, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {

		//@todo self edit  value , because change table field is trouble in here , so self fix field is easy
		query := fmt.Sprintf("insert into .... (%s) values ...", m.table)
		if session != nil {
			//@todo self edit  value , because change table field is trouble in here , so self fix field is easy
			return session.Exec(query, data.DeleteTime, data.DelState, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info)
		}
		//@todo self edit  value , because change table field is trouble in here , so self fix field is easy
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info)
	}, userIdKey, userMobileKey)

}

func (m *defaultUserModel) FindOne(id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRow(&resp, userIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
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

func (m *defaultUserModel) FindOneByMobile(mobile string) (*User, error) {
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, mobile)
	var resp User
	err := m.QueryRowIndex(&resp, userMobileKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `mobile` = ? limit 1", userRows, m.table)
		if err := conn.QueryRow(&resp, query, mobile); err != nil {
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

func (m *defaultUserModel) Update(session sqlx.Session, data *User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info, data.Id)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info, data.Id)
	}, userIdKey, userMobileKey)
	return err
}

func (m *defaultUserModel) Delete(session sqlx.Session, data *User) error {
	data.DelState = globalkey.DelStateYes
	return m.Update(session, data)
}

func (m *defaultUserModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}

func (m *defaultUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}

func (m *defaultUserModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
	return conn.QueryRow(v, query, primary)
}
