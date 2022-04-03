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
	walletFieldNames          = builder.RawFieldNames(&Wallet{})
	walletRows                = strings.Join(walletFieldNames, ",")
	walletRowsExpectAutoSet   = strings.Join(stringx.Remove(walletFieldNames, "`created_at`", "`updated_at`"), ",")
	walletRowsWithPlaceHolder = strings.Join(stringx.Remove(walletFieldNames, "`id`", "`created_at`", "`updated_at`"), "=?,") + "=?"

	cacheWalletIdPrefix     = "cache:wallet:id:"
	cacheWalletUserIdPrefix = "cache:wallet:userId:"
)

type (
	WalletModel interface {
		// Insert 新增数据
		Insert(session sqlx.Session, data *Wallet) (sql.Result, error)
		// FindOne 根据主键查询一条数据，走缓存
		FindOne(id int64) (*Wallet, error)
		// FindOneByUserId 根据唯一索引查询一条数据，走缓存
		FindOneByUserId(userId int64) (*Wallet, error)
		// Delete 删除数据
		Delete(session sqlx.Session, id int64) error
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *Wallet) error
		// Update 更新数据
		Update(session sqlx.Session, data *Wallet) (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session, data *Wallet) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*Wallet, error)
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64, error)
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64, error)
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Wallet, error)
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Wallet, error)
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Wallet, error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Wallet, error)
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
	}

	defaultWalletModel struct {
		sqlc.CachedConn
		table string
	}

	Wallet struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		Version    int64     `db:"version"`
		UserId     int64     `db:"user_id"`
		Money      int64     `db:"money"`
	}
)

func NewWalletModel(conn sqlx.SqlConn, c cache.CacheConf) WalletModel {
	return &defaultWalletModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`wallet`",
	}
}

// Insert 新增数据
func (m *defaultWalletModel) Insert(session sqlx.Session, data *Wallet) (sql.Result, error) {

	data.DeleteTime = time.Unix(0, 0)

	walletUserIdKey := fmt.Sprintf("%s%v", cacheWalletUserIdPrefix, data.UserId)
	walletIdKey := fmt.Sprintf("%s%v", cacheWalletIdPrefix, data.Id)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, walletRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.Id, data.DeleteTime, data.DelState, data.Version, data.UserId, data.Money)
		}
		return conn.Exec(query, data.Id, data.DeleteTime, data.DelState, data.Version, data.UserId, data.Money)
	}, walletIdKey, walletUserIdKey)

}

// FindOne 根据主键查询一条数据，走缓存
func (m *defaultWalletModel) FindOne(id int64) (*Wallet, error) {
	walletIdKey := fmt.Sprintf("%s%v", cacheWalletIdPrefix, id)
	var resp Wallet
	err := m.QueryRow(&resp, walletIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", walletRows, m.table)
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

// FindOneByUserId 根据唯一索引查询一条数据，走缓存
func (m *defaultWalletModel) FindOneByUserId(userId int64) (*Wallet, error) {
	walletUserIdKey := fmt.Sprintf("%s%v", cacheWalletUserIdPrefix, userId)
	var resp Wallet
	err := m.QueryRowIndex(&resp, walletUserIdKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and del_state = ?  limit 1", walletRows, m.table)
		if err := conn.QueryRow(&resp, query, userId, globalkey.DelStateNo); err != nil {
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
func (m *defaultWalletModel) Update(session sqlx.Session, data *Wallet) (sql.Result, error) {
	walletIdKey := fmt.Sprintf("%s%v", cacheWalletIdPrefix, data.Id)
	walletUserIdKey := fmt.Sprintf("%s%v", cacheWalletUserIdPrefix, data.UserId)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, walletRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.Money, data.Id)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.Money, data.Id)
	}, walletIdKey, walletUserIdKey)
}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *defaultWalletModel) UpdateWithVersion(session sqlx.Session, data *Wallet) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	walletIdKey := fmt.Sprintf("%s%v", cacheWalletIdPrefix, data.Id)
	walletUserIdKey := fmt.Sprintf("%s%v", cacheWalletUserIdPrefix, data.UserId)
	sqlResult, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, walletRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.Money, data.Id, oldVersion)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.Money, data.Id, oldVersion)
	}, walletIdKey, walletUserIdKey)
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
func (m *defaultWalletModel) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*Wallet, error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp Wallet
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *defaultWalletModel) FindSum(sumBuilder squirrel.SelectBuilder) (float64, error) {

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
func (m *defaultWalletModel) FindCount(countBuilder squirrel.SelectBuilder) (int64, error) {

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
func (m *defaultWalletModel) FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Wallet, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Wallet
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 按照页码分页查询数据
func (m *defaultWalletModel) FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Wallet, error) {

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

	var resp []*Wallet
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *defaultWalletModel) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Wallet, error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? ", preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Wallet
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *defaultWalletModel) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Wallet, error) {

	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? ", preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Wallet
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *defaultWalletModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(walletRows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *defaultWalletModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *defaultWalletModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

// Delete 删除数据
func (m *defaultWalletModel) Delete(session sqlx.Session, id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	walletIdKey := fmt.Sprintf("%s%v", cacheWalletIdPrefix, id)
	walletUserIdKey := fmt.Sprintf("%s%v", cacheWalletUserIdPrefix, data.UserId)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, walletIdKey, walletUserIdKey)
	return err
}

// DeleteSoft 软删除数据
func (m *defaultWalletModel) DeleteSoft(session sqlx.Session, data *Wallet) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(session, data); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "WalletModel delete err : %+v", err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *defaultWalletModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}

// formatPrimary 格式化缓存key
func (m *defaultWalletModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheWalletIdPrefix, primary)
}

// queryPrimary 根据主键去db查询一条数据
func (m *defaultWalletModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", walletRows, m.table)
	return conn.QueryRow(v, query, primary, globalkey.DelStateNo)
}

//!!!!! 其他自定义方法，从此处开始写,此处上方不要写自定义方法!!!!!
