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
	flightsFieldNames          = builder.RawFieldNames(&Flights{})
	flightsRows                = strings.Join(flightsFieldNames, ",")
	flightsRowsExpectAutoSet   = strings.Join(stringx.Remove(flightsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	flightsRowsWithPlaceHolder = strings.Join(stringx.Remove(flightsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheFlightsIdPrefix     = "cache:flights:id:"
	cacheFlightsNumberPrefix = "cache:flights:number:"
)

type (
	FlightsModel interface {
		// Insert 新增数据
		Insert(session sqlx.Session, data *Flights) (sql.Result, error)
		// FindOne 根据主键查询一条数据，走缓存
		FindOne(id int64) (*Flights, error)
		// FindOneByNumber 根据唯一索引查询一条数据，走缓存
		FindOneByNumber(number sql.NullString) (*Flights, error)
		// Delete 删除数据
		Delete(session sqlx.Session, id int64) error
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *Flights) error
		// Update 更新数据
		Update(session sqlx.Session, data *Flights) (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session, data *Flights) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*Flights, error)
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64, error)
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64, error)
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Flights, error)
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Flights, error)
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Flights, error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Flights, error)
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
	}

	defaultFlightsModel struct {
		sqlc.CachedConn
		table string
	}

	Flights struct {
		Id         int64          `db:"id"`
		CreatedAt  sql.NullTime   `db:"created_at"`
		UpdatedAt  sql.NullTime   `db:"updated_at"`
		DeletedAt  sql.NullTime   `db:"deleted_at"`
		DelState   int64          `db:"del_state"`
		Version    int64          `db:"version"`
		Number     sql.NullString `db:"number"`
		FltTypeJmp sql.NullString `db:"flt_type_jmp"`
	}
)

func NewFlightsModel(conn sqlx.SqlConn, c cache.CacheConf) FlightsModel {
	return &defaultFlightsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`flights`",
	}
}

// Insert 新增数据
func (m *defaultFlightsModel) Insert(session sqlx.Session, data *Flights) (sql.Result, error) {

	data.DeletedAt.Time = time.Unix(0, 0)

	flightsIdKey := fmt.Sprintf("%s%v", cacheFlightsIdPrefix, data.Id)
	flightsNumberKey := fmt.Sprintf("%s%v", cacheFlightsNumberPrefix, data.Number)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, flightsRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.DelState, data.Version, data.Number, data.FltTypeJmp)
		}
		return conn.Exec(query, data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.DelState, data.Version, data.Number, data.FltTypeJmp)
	}, flightsIdKey, flightsNumberKey)

}

// FindOne 根据主键查询一条数据，走缓存
func (m *defaultFlightsModel) FindOne(id int64) (*Flights, error) {
	flightsIdKey := fmt.Sprintf("%s%v", cacheFlightsIdPrefix, id)
	var resp Flights
	err := m.QueryRow(&resp, flightsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", flightsRows, m.table)
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

// FindOneBy 根据唯一索引查询一条数据，走缓存
func (m *defaultFlightsModel) FindOneByNumber(number sql.NullString) (*Flights, error) {
	flightsNumberKey := fmt.Sprintf("%s%v", cacheFlightsNumberPrefix, number)
	var resp Flights
	err := m.QueryRowIndex(&resp, flightsNumberKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `number` = ? and del_state = ?  limit 1", flightsRows, m.table)
		if err := conn.QueryRow(&resp, query, number, globalkey.DelStateNo); err != nil {
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
func (m *defaultFlightsModel) Update(session sqlx.Session, data *Flights) (sql.Result, error) {
	flightsIdKey := fmt.Sprintf("%s%v", cacheFlightsIdPrefix, data.Id)
	flightsNumberKey := fmt.Sprintf("%s%v", cacheFlightsNumberPrefix, data.Number)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, flightsRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.DelState, data.Version, data.Number, data.FltTypeJmp, data.Id)
		}
		return conn.Exec(query, data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.DelState, data.Version, data.Number, data.FltTypeJmp, data.Id)
	}, flightsIdKey, flightsNumberKey)
}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *defaultFlightsModel) UpdateWithVersion(session sqlx.Session, data *Flights) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	flightsIdKey := fmt.Sprintf("%s%v", cacheFlightsIdPrefix, data.Id)
	flightsNumberKey := fmt.Sprintf("%s%v", cacheFlightsNumberPrefix, data.Number)
	sqlResult, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, flightsRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.DelState, data.Version, data.Number, data.FltTypeJmp, data.Id, oldVersion)
		}
		return conn.Exec(query, data.CreatedAt, data.UpdatedAt, data.DeletedAt, data.DelState, data.Version, data.Number, data.FltTypeJmp, data.Id, oldVersion)
	}, flightsIdKey, flightsNumberKey)
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
func (m *defaultFlightsModel) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*Flights, error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp Flights
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *defaultFlightsModel) FindSum(sumBuilder squirrel.SelectBuilder) (float64, error) {

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
func (m *defaultFlightsModel) FindCount(countBuilder squirrel.SelectBuilder) (int64, error) {

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
func (m *defaultFlightsModel) FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Flights, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Flights
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 按照页码分页查询数据
func (m *defaultFlightsModel) FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Flights, error) {

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

	var resp []*Flights
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *defaultFlightsModel) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Flights, error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? ", preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Flights
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *defaultFlightsModel) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Flights, error) {

	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? ", preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Flights
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *defaultFlightsModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(flightsRows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *defaultFlightsModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *defaultFlightsModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

// Delete 删除数据
func (m *defaultFlightsModel) Delete(session sqlx.Session, id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	flightsIdKey := fmt.Sprintf("%s%v", cacheFlightsIdPrefix, id)
	flightsNumberKey := fmt.Sprintf("%s%v", cacheFlightsNumberPrefix, data.Number)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, flightsIdKey, flightsNumberKey)
	return err
}

// DeleteSoft 软删除数据
func (m *defaultFlightsModel) DeleteSoft(session sqlx.Session, data *Flights) error {
	data.DelState = globalkey.DelStateYes
	data.DeletedAt.Time = time.Now()
	if err := m.UpdateWithVersion(session, data); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "FlightsModel delete err : %+v", err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *defaultFlightsModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}

// formatPrimary 格式化缓存key
func (m *defaultFlightsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheFlightsIdPrefix, primary)
}

// queryPrimary 根据主键去db查询一条数据
func (m *defaultFlightsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", flightsRows, m.table)
	return conn.QueryRow(v, query, primary, globalkey.DelStateNo)
}

//!!!!! 其他自定义方法，从此处开始写,此处上方不要写自定义方法!!!!!
