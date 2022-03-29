package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"uranus/common/globalkey"
	"uranus/common/xerr"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userTicketsFieldNames          = builder.RawFieldNames(&UserTickets{})
	userTicketsRows                = strings.Join(userTicketsFieldNames, ",")
	userTicketsRowsExpectAutoSet   = strings.Join(stringx.Remove(userTicketsFieldNames, "`id`", "`created_at`", "`updated_at`"), ",")
	userTicketsRowsWithPlaceHolder = strings.Join(stringx.Remove(userTicketsFieldNames, "`id`", "`created_at`", "`updated_at`"), "=?,") + "=?"

	cacheUserTicketsIdPrefix              = "cache:userTickets:id:"
	cacheUserTicketsAuthKeyPrefix         = "cache:userTickets:authKey:"
	cacheUserTicketsAuthKeyTicketIdPrefix = "cache:userTickets:authKey:ticketId:"
)

type (
	UserTicketsModel interface {
		// Insert 新增数据
		Insert(session sqlx.Session, data *UserTickets) (sql.Result, error)
		// FindOne 根据主键查询一条数据，走缓存
		FindOne(id int64) (*UserTickets, error)
		// FindOneByAuthKey 根据唯一索引查询一条数据，走缓存
		FindOneByAuthKey(authKey string) (*UserTickets, error)
		// FindOneByAuthKeyTicketId 根据唯一索引查询一条数据，走缓存
		FindOneByAuthKeyTicketId(authKey string, ticketId int64) (*UserTickets, error)
		// Delete 删除数据
		Delete(session sqlx.Session, id int64) error
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *UserTickets) error
		// Update 更新数据
		Update(session sqlx.Session, data *UserTickets) (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session, data *UserTickets) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*UserTickets, error)
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64, error)
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64, error)
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*UserTickets, error)
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*UserTickets, error)
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*UserTickets, error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*UserTickets, error)
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
	}

	defaultUserTicketsModel struct {
		sqlc.CachedConn
		table string
	}

	UserTickets struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"` // 是否已经删除
		Version    int64     `db:"version"`   // 版本号
		AuthKey    string    `db:"auth_key"`  // 用户平台唯一id
		TicketId   int64     `db:"ticket_id"` // 票id
	}
)

func NewUserTicketsModel(conn sqlx.SqlConn, c cache.CacheConf) UserTicketsModel {
	return &defaultUserTicketsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_tickets`",
	}
}

// Insert 新增数据
func (m *defaultUserTicketsModel) Insert(session sqlx.Session, data *UserTickets) (sql.Result, error) {

	data.DeleteTime = time.Unix(0, 0)

	userTicketsIdKey := fmt.Sprintf("%s%v", cacheUserTicketsIdPrefix, data.Id)
	userTicketsAuthKeyKey := fmt.Sprintf("%s%v", cacheUserTicketsAuthKeyPrefix, data.AuthKey)
	userTicketsAuthKeyTicketIdKey := fmt.Sprintf("%s%v:%v", cacheUserTicketsAuthKeyTicketIdPrefix, data.AuthKey, data.TicketId)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userTicketsRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.AuthKey, data.TicketId)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.AuthKey, data.TicketId)
	}, userTicketsIdKey, userTicketsAuthKeyKey, userTicketsAuthKeyTicketIdKey)

}

// FindOne 根据主键查询一条数据，走缓存
func (m *defaultUserTicketsModel) FindOne(id int64) (*UserTickets, error) {
	userTicketsIdKey := fmt.Sprintf("%s%v", cacheUserTicketsIdPrefix, id)
	var resp UserTickets
	err := m.QueryRow(&resp, userTicketsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", userTicketsRows, m.table)
		return conn.QueryRow(v, query, id, globalkey.DelStateNo)
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

// FindOneByAuthKey 根据唯一索引查询一条数据，走缓存
func (m *defaultUserTicketsModel) FindOneByAuthKey(authKey string) (*UserTickets, error) {
	userTicketsAuthKeyKey := fmt.Sprintf("%s%v", cacheUserTicketsAuthKeyPrefix, authKey)
	var resp UserTickets
	err := m.QueryRowIndex(&resp, userTicketsAuthKeyKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `auth_key` = ? and del_state = ?  limit 1", userTicketsRows, m.table)
		if err := conn.QueryRow(&resp, query, authKey, globalkey.DelStateNo); err != nil {
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

// FindOneByAuthKeyTicketId 根据唯一索引查询一条数据，走缓存
func (m *defaultUserTicketsModel) FindOneByAuthKeyTicketId(authKey string, ticketId int64) (*UserTickets, error) {
	userTicketsAuthKeyTicketIdKey := fmt.Sprintf("%s%v:%v", cacheUserTicketsAuthKeyTicketIdPrefix, authKey, ticketId)
	var resp UserTickets
	err := m.QueryRowIndex(&resp, userTicketsAuthKeyTicketIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `auth_key` = ? and `ticket_id` = ? and del_state = ?  limit 1", userTicketsRows, m.table)
		if err := conn.QueryRow(&resp, query, authKey, ticketId, globalkey.DelStateNo); err != nil {
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

// Update 修改数据 ,推荐优先使用乐观锁更新
func (m *defaultUserTicketsModel) Update(session sqlx.Session, data *UserTickets) (sql.Result, error) {
	userTicketsIdKey := fmt.Sprintf("%s%v", cacheUserTicketsIdPrefix, data.Id)
	userTicketsAuthKeyKey := fmt.Sprintf("%s%v", cacheUserTicketsAuthKeyPrefix, data.AuthKey)
	userTicketsAuthKeyTicketIdKey := fmt.Sprintf("%s%v:%v", cacheUserTicketsAuthKeyTicketIdPrefix, data.AuthKey, data.TicketId)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userTicketsRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.AuthKey, data.TicketId, data.Id)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.AuthKey, data.TicketId, data.Id)
	}, userTicketsIdKey, userTicketsAuthKeyKey, userTicketsAuthKeyTicketIdKey)
}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *defaultUserTicketsModel) UpdateWithVersion(session sqlx.Session, data *UserTickets) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	userTicketsIdKey := fmt.Sprintf("%s%v", cacheUserTicketsIdPrefix, data.Id)
	userTicketsAuthKeyKey := fmt.Sprintf("%s%v", cacheUserTicketsAuthKeyPrefix, data.AuthKey)
	userTicketsAuthKeyTicketIdKey := fmt.Sprintf("%s%v:%v", cacheUserTicketsAuthKeyTicketIdPrefix, data.AuthKey, data.TicketId)
	sqlResult, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, userTicketsRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.AuthKey, data.TicketId, data.Id, oldVersion)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.AuthKey, data.TicketId, data.Id, oldVersion)
	}, userTicketsIdKey, userTicketsAuthKeyKey, userTicketsAuthKeyTicketIdKey)
	if err != nil {
		return err
	}

	updateCount, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}

	if updateCount == 0 {
		return xerr.NewErrCode(xerr.DB_UPDATE_AFFECTED_ZERO_ERROR)
	}

	return nil

}

// FindOneByQuery 根据条件查询一条数据
func (m *defaultUserTicketsModel) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*UserTickets, error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp UserTickets
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *defaultUserTicketsModel) FindSum(sumBuilder squirrel.SelectBuilder) (float64, error) {

	query, values, err := sumBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindCount 根据某个字段查询数据数量
func (m *defaultUserTicketsModel) FindCount(countBuilder squirrel.SelectBuilder) (int64, error) {

	query, values, err := countBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindAll 查询所有数据
func (m *defaultUserTicketsModel) FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*UserTickets, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*UserTickets
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 按照页码分页查询数据
func (m *defaultUserTicketsModel) FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*UserTickets, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*UserTickets
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *defaultUserTicketsModel) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*UserTickets, error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? ", preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*UserTickets
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *defaultUserTicketsModel) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*UserTickets, error) {

	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? ", preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*UserTickets
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *defaultUserTicketsModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(userTicketsRows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *defaultUserTicketsModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *defaultUserTicketsModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

// Delete 删除数据
func (m *defaultUserTicketsModel) Delete(session sqlx.Session, id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	userTicketsAuthKeyTicketIdKey := fmt.Sprintf("%s%v:%v", cacheUserTicketsAuthKeyTicketIdPrefix, data.AuthKey, data.TicketId)
	userTicketsIdKey := fmt.Sprintf("%s%v", cacheUserTicketsIdPrefix, id)
	userTicketsAuthKeyKey := fmt.Sprintf("%s%v", cacheUserTicketsAuthKeyPrefix, data.AuthKey)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, userTicketsIdKey, userTicketsAuthKeyKey, userTicketsAuthKeyTicketIdKey)
	return err
}

// DeleteSoft 软删除数据
func (m *defaultUserTicketsModel) DeleteSoft(session sqlx.Session, data *UserTickets) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(session, data); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "UserTicketsModel delete err : %+v", err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *defaultUserTicketsModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}

// formatPrimary 格式化缓存key
func (m *defaultUserTicketsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserTicketsIdPrefix, primary)
}

// queryPrimary 根据主键去db查询一条数据
func (m *defaultUserTicketsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", userTicketsRows, m.table)
	return conn.QueryRow(v, query, primary, globalkey.DelStateNo)
}

//!!!!! 其他自定义方法，从此处开始写,此处上方不要写自定义方法!!!!!
