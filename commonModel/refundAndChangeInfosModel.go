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
	refundAndChangeInfosFieldNames          = builder.RawFieldNames(&RefundAndChangeInfos{})
	refundAndChangeInfosRows                = strings.Join(refundAndChangeInfosFieldNames, ",")
	refundAndChangeInfosRowsExpectAutoSet   = strings.Join(stringx.Remove(refundAndChangeInfosFieldNames, "`id`", "`created_at`", "`updated_at`"), ",")
	refundAndChangeInfosRowsWithPlaceHolder = strings.Join(stringx.Remove(refundAndChangeInfosFieldNames, "`id`", "`created_at`", "`updated_at`"), "=?,") + "=?"

	cacheRefundAndChangeInfosIdPrefix               = "cache:refundAndChangeInfos:id:"
	cacheRefundAndChangeInfosTicketIdIsRefundPrefix = "cache:refundAndChangeInfos:ticketId:isRefund:"
)

type (
	RefundAndChangeInfosModel interface {
		// Insert 新增数据
		Insert(session sqlx.Session, data *RefundAndChangeInfos) (sql.Result, error)
		// FindOne 根据主键查询一条数据，走缓存
		FindOne(id int64) (*RefundAndChangeInfos, error)
		// FindOneByTicketIdIsRefund 根据唯一索引查询一条数据，走缓存
		FindOneByTicketIdIsRefund(ticketId int64, isRefund int64) (*RefundAndChangeInfos, error)
		// Delete 删除数据
		Delete(session sqlx.Session, id int64) error
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *RefundAndChangeInfos) error
		// Update 更新数据
		Update(session sqlx.Session, data *RefundAndChangeInfos) (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session, data *RefundAndChangeInfos) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*RefundAndChangeInfos, error)
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64, error)
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64, error)
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*RefundAndChangeInfos, error)
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*RefundAndChangeInfos, error)
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*RefundAndChangeInfos, error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*RefundAndChangeInfos, error)
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
	}

	defaultRefundAndChangeInfosModel struct {
		sqlc.CachedConn
		table string
	}

	RefundAndChangeInfos struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"` // 是否已经删除
		Version    int64     `db:"version"`   // 版本号
		TicketId   int64     `db:"ticket_id"` // 对应的票ID
		IsRefund   int64     `db:"is_refund"` // 1为退订，0为改票
		Time1      time.Time `db:"time1"`     // 时间1
		Fee1       int64     `db:"fee1"`      // 符合时间1时需要的费用(￥/人)
		Time2      time.Time `db:"time2"`     // 时间2
		Fee2       int64     `db:"fee2"`      // 符合时间2时需要的费用(￥/人)
		Time3      time.Time `db:"time3"`     // 时间3
		Fee3       int64     `db:"fee3"`      // 符合时间3时需要的费用(￥/人)
		Time4      time.Time `db:"time4"`     // 时间4
		Fee4       int64     `db:"fee4"`      // 符合时间4时需要的费用(￥/人)
		Time5      time.Time `db:"time5"`     // 时间5
		Fee5       int64     `db:"fee5"`      // 符合时间5时需要的费用(￥/人)
	}
)

func NewRefundAndChangeInfosModel(conn sqlx.SqlConn, c cache.CacheConf) RefundAndChangeInfosModel {
	return &defaultRefundAndChangeInfosModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`refund_and_change_infos`",
	}
}

// Insert 新增数据
func (m *defaultRefundAndChangeInfosModel) Insert(session sqlx.Session, data *RefundAndChangeInfos) (sql.Result, error) {

	data.DeleteTime = time.Unix(0, 0)

	refundAndChangeInfosIdKey := fmt.Sprintf("%s%v", cacheRefundAndChangeInfosIdPrefix, data.Id)
	refundAndChangeInfosTicketIdIsRefundKey := fmt.Sprintf("%s%v:%v", cacheRefundAndChangeInfosTicketIdIsRefundPrefix, data.TicketId, data.IsRefund)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, refundAndChangeInfosRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.TicketId, data.IsRefund, data.Time1, data.Fee1, data.Time2, data.Fee2, data.Time3, data.Fee3, data.Time4, data.Fee4, data.Time5, data.Fee5)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.TicketId, data.IsRefund, data.Time1, data.Fee1, data.Time2, data.Fee2, data.Time3, data.Fee3, data.Time4, data.Fee4, data.Time5, data.Fee5)
	}, refundAndChangeInfosTicketIdIsRefundKey, refundAndChangeInfosIdKey)

}

// FindOne 根据主键查询一条数据，走缓存
func (m *defaultRefundAndChangeInfosModel) FindOne(id int64) (*RefundAndChangeInfos, error) {
	refundAndChangeInfosIdKey := fmt.Sprintf("%s%v", cacheRefundAndChangeInfosIdPrefix, id)
	var resp RefundAndChangeInfos
	err := m.QueryRow(&resp, refundAndChangeInfosIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", refundAndChangeInfosRows, m.table)
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

// FindOneByTicketIdIsRefund 根据唯一索引查询一条数据，走缓存
func (m *defaultRefundAndChangeInfosModel) FindOneByTicketIdIsRefund(ticketId int64, isRefund int64) (*RefundAndChangeInfos, error) {
	refundAndChangeInfosTicketIdIsRefundKey := fmt.Sprintf("%s%v:%v", cacheRefundAndChangeInfosTicketIdIsRefundPrefix, ticketId, isRefund)
	var resp RefundAndChangeInfos
	err := m.QueryRowIndex(&resp, refundAndChangeInfosTicketIdIsRefundKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `ticket_id` = ? and `is_refund` = ? and del_state = ?  limit 1", refundAndChangeInfosRows, m.table)
		if err := conn.QueryRow(&resp, query, ticketId, isRefund, globalkey.DelStateNo); err != nil {
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
func (m *defaultRefundAndChangeInfosModel) Update(session sqlx.Session, data *RefundAndChangeInfos) (sql.Result, error) {
	refundAndChangeInfosIdKey := fmt.Sprintf("%s%v", cacheRefundAndChangeInfosIdPrefix, data.Id)
	refundAndChangeInfosTicketIdIsRefundKey := fmt.Sprintf("%s%v:%v", cacheRefundAndChangeInfosTicketIdIsRefundPrefix, data.TicketId, data.IsRefund)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, refundAndChangeInfosRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.TicketId, data.IsRefund, data.Time1, data.Fee1, data.Time2, data.Fee2, data.Time3, data.Fee3, data.Time4, data.Fee4, data.Time5, data.Fee5, data.Id)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.TicketId, data.IsRefund, data.Time1, data.Fee1, data.Time2, data.Fee2, data.Time3, data.Fee3, data.Time4, data.Fee4, data.Time5, data.Fee5, data.Id)
	}, refundAndChangeInfosIdKey, refundAndChangeInfosTicketIdIsRefundKey)
}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *defaultRefundAndChangeInfosModel) UpdateWithVersion(session sqlx.Session, data *RefundAndChangeInfos) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	refundAndChangeInfosIdKey := fmt.Sprintf("%s%v", cacheRefundAndChangeInfosIdPrefix, data.Id)
	refundAndChangeInfosTicketIdIsRefundKey := fmt.Sprintf("%s%v:%v", cacheRefundAndChangeInfosTicketIdIsRefundPrefix, data.TicketId, data.IsRefund)
	sqlResult, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, refundAndChangeInfosRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.TicketId, data.IsRefund, data.Time1, data.Fee1, data.Time2, data.Fee2, data.Time3, data.Fee3, data.Time4, data.Fee4, data.Time5, data.Fee5, data.Id, oldVersion)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.TicketId, data.IsRefund, data.Time1, data.Fee1, data.Time2, data.Fee2, data.Time3, data.Fee3, data.Time4, data.Fee4, data.Time5, data.Fee5, data.Id, oldVersion)
	}, refundAndChangeInfosIdKey, refundAndChangeInfosTicketIdIsRefundKey)
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
func (m *defaultRefundAndChangeInfosModel) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*RefundAndChangeInfos, error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp RefundAndChangeInfos
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *defaultRefundAndChangeInfosModel) FindSum(sumBuilder squirrel.SelectBuilder) (float64, error) {

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
func (m *defaultRefundAndChangeInfosModel) FindCount(countBuilder squirrel.SelectBuilder) (int64, error) {

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
func (m *defaultRefundAndChangeInfosModel) FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*RefundAndChangeInfos, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*RefundAndChangeInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 按照页码分页查询数据
func (m *defaultRefundAndChangeInfosModel) FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*RefundAndChangeInfos, error) {

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

	var resp []*RefundAndChangeInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *defaultRefundAndChangeInfosModel) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*RefundAndChangeInfos, error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? ", preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*RefundAndChangeInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *defaultRefundAndChangeInfosModel) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*RefundAndChangeInfos, error) {

	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? ", preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*RefundAndChangeInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *defaultRefundAndChangeInfosModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(refundAndChangeInfosRows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *defaultRefundAndChangeInfosModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *defaultRefundAndChangeInfosModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

// Delete 删除数据
func (m *defaultRefundAndChangeInfosModel) Delete(session sqlx.Session, id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	refundAndChangeInfosIdKey := fmt.Sprintf("%s%v", cacheRefundAndChangeInfosIdPrefix, id)
	refundAndChangeInfosTicketIdIsRefundKey := fmt.Sprintf("%s%v:%v", cacheRefundAndChangeInfosTicketIdIsRefundPrefix, data.TicketId, data.IsRefund)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, refundAndChangeInfosIdKey, refundAndChangeInfosTicketIdIsRefundKey)
	return err
}

// DeleteSoft 软删除数据
func (m *defaultRefundAndChangeInfosModel) DeleteSoft(session sqlx.Session, data *RefundAndChangeInfos) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(session, data); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "RefundAndChangeInfosModel delete err : %+v", err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *defaultRefundAndChangeInfosModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}

// formatPrimary 格式化缓存key
func (m *defaultRefundAndChangeInfosModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheRefundAndChangeInfosIdPrefix, primary)
}

// queryPrimary 根据主键去db查询一条数据
func (m *defaultRefundAndChangeInfosModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", refundAndChangeInfosRows, m.table)
	return conn.QueryRow(v, query, primary, globalkey.DelStateNo)
}

//!!!!! 其他自定义方法，从此处开始写,此处上方不要写自定义方法!!!!!
