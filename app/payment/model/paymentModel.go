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
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	paymentFieldNames          = builder.RawFieldNames(&Payment{})
	paymentRows                = strings.Join(paymentFieldNames, ",")
	paymentRowsExpectAutoSet   = strings.Join(stringx.Remove(paymentFieldNames, "`id`", "`created_at`", "`updated_at`"), ",")
	paymentRowsWithPlaceHolder = strings.Join(stringx.Remove(paymentFieldNames, "`id`", "`created_at`", "`updated_at`"), "=?,") + "=?"
)

type (
	PaymentModel interface {
		// Insert 新增数据
		Insert(session sqlx.Session, data *Payment) (sql.Result, error)
		// FindOne 根据主键查询一条数据，走缓存
		FindOne(id int64) (*Payment, error)
		// FindOneBy 根据唯一索引查询一条数据，走缓存
		FindOneBySn(sn string) (*Payment, error)
		// Delete 删除数据
		Delete(session sqlx.Session, id int64) error
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *Payment) error
		// Update 更新数据
		Update(session sqlx.Session, data *Payment) (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session, data *Payment) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*Payment, error)
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64, error)
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64, error)
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Payment, error)
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Payment, error)
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Payment, error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Payment, error)
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
	}

	defaultPaymentModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Payment struct {
		Id             int64     `db:"id"`
		Sn             string    `db:"sn"` // 流水单号
		CreateTime     time.Time `db:"create_time"`
		UpdateTime     time.Time `db:"update_time"`
		DeleteTime     time.Time `db:"delete_time"`
		DelState       int64     `db:"del_state"`
		Version        int64     `db:"version"`          // 乐观锁版本号
		UserId         int64     `db:"user_id"`          // 用户id
		PayMode        string    `db:"pay_mode"`         // 支付方式 1:微信支付 2:支付宝 3:钱包余额
		TradeType      string    `db:"trade_type"`       // 第三方支付类型
		TradeState     string    `db:"trade_state"`      // 第三方交易状态
		PayTotal       int64     `db:"pay_total"`        // 支付总金额(分)
		TransactionId  string    `db:"transaction_id"`   // 第三方支付单号
		TradeStateDesc string    `db:"trade_state_desc"` // 支付状态描述
		OrderSn        string    `db:"order_sn"`         // 业务单号
		PayStatus      int64     `db:"pay_status"`       // 平台内交易状态   -1:支付失败 0:未支付 1:支付成功 2:已退款
		PayTime        time.Time `db:"pay_time"`         // 支付成功时间
	}
)

func NewPaymentModel(conn sqlx.SqlConn) PaymentModel {
	return &defaultPaymentModel{
		conn:  conn,
		table: "`payment`",
	}
}

// Insert 新增数据
func (m *defaultPaymentModel) Insert(session sqlx.Session, data *Payment) (sql.Result, error) {

	data.DeleteTime = time.Unix(0, 0)

	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, paymentRowsExpectAutoSet)
	if session != nil {
		return session.Exec(query, data.Sn, data.DeleteTime, data.DelState, data.Version, data.UserId, data.PayMode, data.TradeType, data.TradeState, data.PayTotal, data.TransactionId, data.TradeStateDesc, data.OrderSn, data.PayStatus, data.PayTime)
	}
	return m.conn.Exec(query, data.Sn, data.DeleteTime, data.DelState, data.Version, data.UserId, data.PayMode, data.TradeType, data.TradeState, data.PayTotal, data.TransactionId, data.TradeStateDesc, data.OrderSn, data.PayStatus, data.PayTime)

}

// FindOne 根据主键查询一条数据，走缓存
func (m *defaultPaymentModel) FindOne(id int64) (*Payment, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ?  and del_state = ? limit 1", paymentRows, m.table)
	var resp Payment
	err := m.conn.QueryRow(&resp, query, id, globalkey.DelStateNo)
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
func (m *defaultPaymentModel) FindOneBySn(sn string) (*Payment, error) {
	var resp Payment
	query := fmt.Sprintf("select %s from %s where `sn` = ? and del_state = ? limit 1", paymentRows, m.table)
	err := m.conn.QueryRow(&resp, query, sn, globalkey.DelStateNo)
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
func (m *defaultPaymentModel) Update(session sqlx.Session, data *Payment) (sql.Result, error) {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, paymentRowsWithPlaceHolder)
	if session != nil {
		return session.Exec(query, data.Sn, data.DeleteTime, data.DelState, data.Version, data.UserId, data.PayMode, data.TradeType, data.TradeState, data.PayTotal, data.TransactionId, data.TradeStateDesc, data.OrderSn, data.PayStatus, data.PayTime, data.Id)
	} else {
		return m.conn.Exec(query, data.Sn, data.DeleteTime, data.DelState, data.Version, data.UserId, data.PayMode, data.TradeType, data.TradeState, data.PayTotal, data.TransactionId, data.TradeStateDesc, data.OrderSn, data.PayStatus, data.PayTime, data.Id)
	}

}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *defaultPaymentModel) UpdateWithVersion(session sqlx.Session, data *Payment) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, paymentRowsWithPlaceHolder)
	if session != nil {
		sqlResult, err = session.Exec(query, data.Sn, data.DeleteTime, data.DelState, data.Version, data.UserId, data.PayMode, data.TradeType, data.TradeState, data.PayTotal, data.TransactionId, data.TradeStateDesc, data.OrderSn, data.PayStatus, data.PayTime, data.Id, oldVersion)
	} else {
		sqlResult, err = m.conn.Exec(query, data.Sn, data.DeleteTime, data.DelState, data.Version, data.UserId, data.PayMode, data.TradeType, data.TradeState, data.PayTotal, data.TransactionId, data.TradeStateDesc, data.OrderSn, data.PayStatus, data.PayTime, data.Id, oldVersion)
	}

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
func (m *defaultPaymentModel) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*Payment, error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp Payment

	err = m.conn.QueryRow(&resp, query, values...)

	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *defaultPaymentModel) FindSum(sumBuilder squirrel.SelectBuilder) (float64, error) {

	query, values, err := sumBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64

	err = m.conn.QueryRow(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindCount 根据某个字段查询数据数量
func (m *defaultPaymentModel) FindCount(countBuilder squirrel.SelectBuilder) (int64, error) {

	query, values, err := countBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64

	err = m.conn.QueryRow(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindAll 查询所有数据
func (m *defaultPaymentModel) FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*Payment, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Payment

	err = m.conn.QueryRows(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 按照页码分页查询数据
func (m *defaultPaymentModel) FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*Payment, error) {

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

	var resp []*Payment

	err = m.conn.QueryRows(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *defaultPaymentModel) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*Payment, error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? ", preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Payment

	err = m.conn.QueryRows(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *defaultPaymentModel) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*Payment, error) {

	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? ", preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*Payment

	err = m.conn.QueryRows(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *defaultPaymentModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(paymentRows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *defaultPaymentModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *defaultPaymentModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

// Delete 删除数据
func (m *defaultPaymentModel) Delete(session sqlx.Session, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	if session != nil {
		_, err := session.Exec(query, id)
		return err
	}
	_, err := m.conn.Exec(query, id)
	return err
}

// DeleteSoft 软删除数据
func (m *defaultPaymentModel) DeleteSoft(session sqlx.Session, data *Payment) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(session, data); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "PaymentModel delete err : %+v", err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *defaultPaymentModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.conn.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}
