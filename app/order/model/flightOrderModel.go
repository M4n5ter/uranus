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
	flightOrderFieldNames          = builder.RawFieldNames(&FlightOrder{})
	flightOrderRows                = strings.Join(flightOrderFieldNames, ",")
	flightOrderRowsExpectAutoSet   = strings.Join(stringx.Remove(flightOrderFieldNames, "`id`", "`created_at`", "`updated_at`"), ",")
	flightOrderRowsWithPlaceHolder = strings.Join(stringx.Remove(flightOrderFieldNames, "`id`", "`created_at`", "`updated_at`"), "=?,") + "=?"
)

type (
	FlightOrderModel interface {
		// Insert 新增数据
		Insert(session sqlx.Session, data *FlightOrder) (sql.Result, error)
		// FindOne 根据主键查询一条数据，走缓存
		FindOne(id int64) (*FlightOrder, error)
		// FindOneBySn 根据唯一索引查询一条数据，走缓存
		FindOneBySn(sn string) (*FlightOrder, error)
		// Delete 删除数据
		Delete(session sqlx.Session, id int64) error
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *FlightOrder) error
		// Update 更新数据
		Update(session sqlx.Session, data *FlightOrder) (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session, data *FlightOrder) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*FlightOrder, error)
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64, error)
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64, error)
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*FlightOrder, error)
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*FlightOrder, error)
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*FlightOrder, error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*FlightOrder, error)
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
	}

	defaultFlightOrderModel struct {
		conn  sqlx.SqlConn
		table string
	}

	FlightOrder struct {
		Id              int64     `db:"id"`
		CreateTime      time.Time `db:"create_time"`
		UpdateTime      time.Time `db:"update_time"`
		DeleteTime      time.Time `db:"delete_time"`
		DelState        int64     `db:"del_state"`
		Version         int64     `db:"version"`           // 版本号
		Sn              string    `db:"sn"`                // 订单号
		UserId          int64     `db:"user_id"`           // 下单用户id
		TicketId        int64     `db:"ticket_id"`         // 票id
		DepartPosition  string    `db:"depart_position"`   // 起飞地点
		DepartTime      time.Time `db:"depart_time"`       // 起飞时间
		ArrivePosition  string    `db:"arrive_position"`   // 降落地点
		ArriveTime      time.Time `db:"arrive_time"`       // 降落时间
		TicketPrice     int64     `db:"ticket_price"`      // 票价(分)
		Discount        int64     `db:"discount"`          // 折扣(-n%)
		TradeState      int64     `db:"trade_state"`       // -1: 已取消 0:待支付 1:未使用 2:已使用  3:已退款 4:已过期
		TradeCode       string    `db:"trade_code"`        // 确认码
		OrderTotalPrice int64     `db:"order_total_price"` // 订单总价格(分)
	}
)

func NewFlightOrderModel(conn sqlx.SqlConn) FlightOrderModel {
	return &defaultFlightOrderModel{
		conn:  conn,
		table: "`flight_order`",
	}
}

// Insert 新增数据
func (m *defaultFlightOrderModel) Insert(session sqlx.Session, data *FlightOrder) (sql.Result, error) {

	data.DeleteTime = time.Unix(0, 0)

	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, flightOrderRowsExpectAutoSet)
	if session != nil {
		return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.Sn, data.UserId, data.TicketId, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.TicketPrice, data.Discount, data.TradeState, data.TradeCode, data.OrderTotalPrice)
	}
	return m.conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.Sn, data.UserId, data.TicketId, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.TicketPrice, data.Discount, data.TradeState, data.TradeCode, data.OrderTotalPrice)

}

// FindOne 根据主键查询一条数据，走缓存
func (m *defaultFlightOrderModel) FindOne(id int64) (*FlightOrder, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ?  and del_state = ? limit 1", flightOrderRows, m.table)
	var resp FlightOrder
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

// FindOneBySn 根据唯一索引查询一条数据，走缓存
func (m *defaultFlightOrderModel) FindOneBySn(sn string) (*FlightOrder, error) {
	var resp FlightOrder
	query := fmt.Sprintf("select %s from %s where `sn` = ? and del_state = ? limit 1", flightOrderRows, m.table)
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
func (m *defaultFlightOrderModel) Update(session sqlx.Session, data *FlightOrder) (sql.Result, error) {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, flightOrderRowsWithPlaceHolder)
	if session != nil {
		return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.Sn, data.UserId, data.TicketId, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.TicketPrice, data.Discount, data.TradeState, data.TradeCode, data.OrderTotalPrice, data.Id)
	} else {
		return m.conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.Sn, data.UserId, data.TicketId, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.TicketPrice, data.Discount, data.TradeState, data.TradeCode, data.OrderTotalPrice, data.Id)
	}

}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *defaultFlightOrderModel) UpdateWithVersion(session sqlx.Session, data *FlightOrder) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, flightOrderRowsWithPlaceHolder)
	if session != nil {
		sqlResult, err = session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.Sn, data.UserId, data.TicketId, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.TicketPrice, data.Discount, data.TradeState, data.TradeCode, data.OrderTotalPrice, data.Id, oldVersion)
	} else {
		sqlResult, err = m.conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.Sn, data.UserId, data.TicketId, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.TicketPrice, data.Discount, data.TradeState, data.TradeCode, data.OrderTotalPrice, data.Id, oldVersion)
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
func (m *defaultFlightOrderModel) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*FlightOrder, error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp FlightOrder

	err = m.conn.QueryRow(&resp, query, values...)

	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *defaultFlightOrderModel) FindSum(sumBuilder squirrel.SelectBuilder) (float64, error) {

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
func (m *defaultFlightOrderModel) FindCount(countBuilder squirrel.SelectBuilder) (int64, error) {

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
func (m *defaultFlightOrderModel) FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*FlightOrder, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightOrder

	err = m.conn.QueryRows(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 按照页码分页查询数据
func (m *defaultFlightOrderModel) FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*FlightOrder, error) {

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

	var resp []*FlightOrder

	err = m.conn.QueryRows(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *defaultFlightOrderModel) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*FlightOrder, error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? ", preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightOrder

	err = m.conn.QueryRows(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *defaultFlightOrderModel) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*FlightOrder, error) {

	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? ", preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightOrder

	err = m.conn.QueryRows(&resp, query, values...)

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *defaultFlightOrderModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(flightOrderRows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *defaultFlightOrderModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *defaultFlightOrderModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

// Delete 删除数据
func (m *defaultFlightOrderModel) Delete(session sqlx.Session, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	if session != nil {
		_, err := session.Exec(query, id)
		return err
	}
	_, err := m.conn.Exec(query, id)
	return err
}

// DeleteSoft 软删除数据
func (m *defaultFlightOrderModel) DeleteSoft(session sqlx.Session, data *FlightOrder) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(session, data); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "FlightOrderModel delete err : %+v", err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *defaultFlightOrderModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.conn.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}
