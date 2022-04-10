package commonModel

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
	spacesFieldNames          = builder.RawFieldNames(&Spaces{})
	spacesRows                = strings.Join(spacesFieldNames, ",")
	spacesRowsExpectAutoSet   = strings.Join(stringx.Remove(spacesFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	spacesRowsWithPlaceHolder = strings.Join(stringx.Remove(spacesFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheSpacesIdPrefix                       = "cache:spaces:id:"
	cacheSpacesFlightInfoIdIsFirstClassPrefix = "cache:spaces:flightInfoId:isFirstClass:"
)

type (
	SpacesModel interface {
		// Insert 新增数据
		Insert(session sqlx.Session, data *Spaces) (sql.Result, error)
		// FindOne 根据主键查询一条数据，走缓存
		FindOne(id int64) (*Spaces, error)
		// FindOneBy 根据唯一索引查询一条数据，走缓存
		FindOneByFlightInfoIdIsFirstClass(flightInfoId int64, isFirstClass int64) (*Spaces, error)
		// Delete 删除数据
		Delete(session sqlx.Session, id int64) error
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *Spaces) error
		// Update 更新数据
		Update(session sqlx.Session, data *Spaces) (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session, data *Spaces) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*Spaces, error)
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64, error)
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64, error)
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Spaces, error)
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Spaces, error)
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Spaces, error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Spaces, error)
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
		// FindListByFlightInfoID 按照flightInfoID查询数据，不支持排序
		FindListByFlightInfoID(rowBuilder squirrel.SelectBuilder, flightInfoID int64) ([]*Spaces, error)
	}

	defaultSpacesModel struct {
		sqlc.CachedConn
		table string
	}

	Spaces struct {
		Id           int64     `db:"id"`
		CreateTime   time.Time `db:"create_time"`
		UpdateTime   time.Time `db:"update_time"`
		DeleteTime   time.Time `db:"delete_time"`
		DelState     int64     `db:"del_state"`      // 是否已经删除
		Version      int64     `db:"version"`        // 版本号
		FlightInfoId int64     `db:"flight_info_id"` // 对应的航班信息id
		IsFirstClass int64     `db:"is_first_class"` // 是否是头等舱/商务舱
		Total        int64     `db:"total"`          // 总量
		Surplus      int64     `db:"surplus"`        // 剩余量
		LockedStock  int64     `db:"locked_stock"`   // 已经被锁定的库存
	}
)

func NewSpacesModel(conn sqlx.SqlConn, c cache.CacheConf) SpacesModel {
	return &defaultSpacesModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`spaces`",
	}
}

// Insert 新增数据
func (m *defaultSpacesModel) Insert(session sqlx.Session, data *Spaces) (sql.Result, error) {

	data.DeleteTime = time.Unix(0, 0)

	spacesIdKey := fmt.Sprintf("%s%v", cacheSpacesIdPrefix, data.Id)
	spacesFlightInfoIdIsFirstClassKey := fmt.Sprintf("%s%v:%v", cacheSpacesFlightInfoIdIsFirstClassPrefix, data.FlightInfoId, data.IsFirstClass)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, spacesRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightInfoId, data.IsFirstClass, data.Total, data.Surplus, data.LockedStock)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightInfoId, data.IsFirstClass, data.Total, data.Surplus, data.LockedStock)
	}, spacesFlightInfoIdIsFirstClassKey, spacesIdKey)

}

// FindOne 根据主键查询一条数据，走缓存
func (m *defaultSpacesModel) FindOne(id int64) (*Spaces, error) {
	spacesIdKey := fmt.Sprintf("%s%v", cacheSpacesIdPrefix, id)
	var resp Spaces
	err := m.QueryRow(&resp, spacesIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", spacesRows, m.table)
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
func (m *defaultSpacesModel) FindOneByFlightInfoIdIsFirstClass(flightInfoId int64, isFirstClass int64) (*Spaces, error) {
	spacesFlightInfoIdIsFirstClassKey := fmt.Sprintf("%s%v:%v", cacheSpacesFlightInfoIdIsFirstClassPrefix, flightInfoId, isFirstClass)
	var resp Spaces
	err := m.QueryRowIndex(&resp, spacesFlightInfoIdIsFirstClassKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `flight_info_id` = ? and `is_first_class` = ? and del_state = ?  limit 1", spacesRows, m.table)
		if err := conn.QueryRow(&resp, query, flightInfoId, isFirstClass, globalkey.DelStateNo); err != nil {
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
func (m *defaultSpacesModel) Update(session sqlx.Session, data *Spaces) (sql.Result, error) {
	spacesFlightInfoIdIsFirstClassKey := fmt.Sprintf("%s%v:%v", cacheSpacesFlightInfoIdIsFirstClassPrefix, data.FlightInfoId, data.IsFirstClass)
	spacesIdKey := fmt.Sprintf("%s%v", cacheSpacesIdPrefix, data.Id)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, spacesRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightInfoId, data.IsFirstClass, data.Total, data.Surplus, data.LockedStock, data.Id)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightInfoId, data.IsFirstClass, data.Total, data.Surplus, data.LockedStock, data.Id)
	}, spacesFlightInfoIdIsFirstClassKey, spacesIdKey)
}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *defaultSpacesModel) UpdateWithVersion(session sqlx.Session, data *Spaces) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	spacesFlightInfoIdIsFirstClassKey := fmt.Sprintf("%s%v:%v", cacheSpacesFlightInfoIdIsFirstClassPrefix, data.FlightInfoId, data.IsFirstClass)
	spacesIdKey := fmt.Sprintf("%s%v", cacheSpacesIdPrefix, data.Id)
	sqlResult, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, spacesRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightInfoId, data.IsFirstClass, data.Total, data.Surplus, data.LockedStock, data.Id, oldVersion)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightInfoId, data.IsFirstClass, data.Total, data.Surplus, data.LockedStock, data.Id, oldVersion)
	}, spacesFlightInfoIdIsFirstClassKey, spacesIdKey)
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
func (m *defaultSpacesModel) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*Spaces, error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp Spaces
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *defaultSpacesModel) FindSum(sumBuilder squirrel.SelectBuilder) (float64, error) {

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
func (m *defaultSpacesModel) FindCount(countBuilder squirrel.SelectBuilder) (int64, error) {

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
func (m *defaultSpacesModel) FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Spaces, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Spaces
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 按照页码分页查询数据
func (m *defaultSpacesModel) FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Spaces, error) {

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

	var resp []*Spaces
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *defaultSpacesModel) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Spaces, error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? ", preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Spaces
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *defaultSpacesModel) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Spaces, error) {

	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? ", preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Spaces
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *defaultSpacesModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(spacesRows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *defaultSpacesModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *defaultSpacesModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

// Delete 删除数据
func (m *defaultSpacesModel) Delete(session sqlx.Session, id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	spacesIdKey := fmt.Sprintf("%s%v", cacheSpacesIdPrefix, id)
	spacesFlightInfoIdIsFirstClassKey := fmt.Sprintf("%s%v:%v", cacheSpacesFlightInfoIdIsFirstClassPrefix, data.FlightInfoId, data.IsFirstClass)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, spacesIdKey, spacesFlightInfoIdIsFirstClassKey)
	return err
}

// DeleteSoft 软删除数据
func (m *defaultSpacesModel) DeleteSoft(session sqlx.Session, data *Spaces) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(session, data); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "SpacesModel delete err : %+v", err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *defaultSpacesModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}

// formatPrimary 格式化缓存key
func (m *defaultSpacesModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheSpacesIdPrefix, primary)
}

// queryPrimary 根据主键去db查询一条数据
func (m *defaultSpacesModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", spacesRows, m.table)
	return conn.QueryRow(v, query, primary, globalkey.DelStateNo)
}

//!!!!! 其他自定义方法，从此处开始写,此处上方不要写自定义方法!!!!!

// FindListByFlightInfoID 按照flightInfoID查询数据，不支持排序
func (m *defaultSpacesModel) FindListByFlightInfoID(rowBuilder squirrel.SelectBuilder, flightInfoID int64) ([]*Spaces, error) {

	if flightInfoID > 0 {
		rowBuilder = rowBuilder.Where(" flight_info_id = ? ", flightInfoID)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Spaces
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
