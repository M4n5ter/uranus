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
	flightInfosFieldNames          = builder.RawFieldNames(&FlightInfos{})
	flightInfosRows                = strings.Join(flightInfosFieldNames, ",")
	flightInfosRowsExpectAutoSet   = strings.Join(stringx.Remove(flightInfosFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	flightInfosRowsWithPlaceHolder = strings.Join(stringx.Remove(flightInfosFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheFlightInfosIdPrefix = "cache:flightInfos:id:"
)

type (
	FlightInfosModel interface {
		// Insert 新增数据
		Insert(session sqlx.Session, data *FlightInfos) (sql.Result, error)
		// FindOne 根据主键查询一条数据，走缓存
		FindOne(id int64) (*FlightInfos, error)
		// Delete 删除数据
		Delete(session sqlx.Session, id int64) error
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *FlightInfos) error
		// Update 更新数据
		Update(session sqlx.Session, data *FlightInfos) (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session, data *FlightInfos) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*FlightInfos, error)
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64, error)
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64, error)
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*FlightInfos, error)
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*FlightInfos, error)
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*FlightInfos, error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*FlightInfos, error)
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
		// FindListByNumber 根据航班号查询数据
		FindListByNumber(rowBuilder squirrel.SelectBuilder, number string) ([]*FlightInfos, error)
		// FindListByNumberAndSetOutDate 根据航班号和出发日期查询数据
		FindListByNumberAndSetOutDate(rowBuilder squirrel.SelectBuilder, number string, sot time.Time) ([]*FlightInfos, error)
		// FindListBySetOutDateAndPosition 通过给定日期、出发地、目的地进行航班查询
		FindListBySetOutDateAndPosition(rowBuilder squirrel.SelectBuilder, sod time.Time, depart, arrive string) ([]*FlightInfos, error)
		// FindOneByByNumberAndSetOutDateAndPosition 根据 航班号 出发日期 始末时间地点 锁定一条航班信息
		FindOneByByNumberAndSetOutDateAndPosition(rowBuilder squirrel.SelectBuilder, number string, sod time.Time, depart string, departTime time.Time, arrive string, arriveTime time.Time) (*FlightInfos, error)
	}

	defaultFlightInfosModel struct {
		sqlc.CachedConn
		table string
	}

	FlightInfos struct {
		Id             int64     `db:"id"`
		CreateTime     time.Time `db:"create_time"`
		UpdateTime     time.Time `db:"update_time"`
		DeleteTime     time.Time `db:"delete_time"`
		DelState       int64     `db:"del_state"`       // 是否已经删除
		Version        int64     `db:"version"`         // 版本号
		FlightNumber   string    `db:"flight_number"`   // 对应的航班号
		SetOutDate     time.Time `db:"set_out_date"`    // 出发日期
		Punctuality    int64     `db:"punctuality"`     // 准点率(%)
		DepartPosition string    `db:"depart_position"` // 起飞地点
		DepartTime     time.Time `db:"depart_time"`     // 起飞时间
		ArrivePosition string    `db:"arrive_position"` // 降落地点
		ArriveTime     time.Time `db:"arrive_time"`     // 降落时间
	}
)

func NewFlightInfosModel(conn sqlx.SqlConn, c cache.CacheConf) FlightInfosModel {
	return &defaultFlightInfosModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`flight_infos`",
	}
}

// Insert 新增数据
func (m *defaultFlightInfosModel) Insert(session sqlx.Session, data *FlightInfos) (sql.Result, error) {

	data.DeleteTime = time.Unix(0, 0)

	flightInfosIdKey := fmt.Sprintf("%s%v", cacheFlightInfosIdPrefix, data.Id)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, flightInfosRowsExpectAutoSet)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightNumber, data.SetOutDate, data.Punctuality, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightNumber, data.SetOutDate, data.Punctuality, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime)
	}, flightInfosIdKey)

}

// FindOne 根据主键查询一条数据，走缓存
func (m *defaultFlightInfosModel) FindOne(id int64) (*FlightInfos, error) {
	flightInfosIdKey := fmt.Sprintf("%s%v", cacheFlightInfosIdPrefix, id)
	var resp FlightInfos
	err := m.QueryRow(&resp, flightInfosIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", flightInfosRows, m.table)
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

// Update 修改数据 ,推荐优先使用乐观锁更新
func (m *defaultFlightInfosModel) Update(session sqlx.Session, data *FlightInfos) (sql.Result, error) {
	flightInfosIdKey := fmt.Sprintf("%s%v", cacheFlightInfosIdPrefix, data.Id)
	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, flightInfosRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightNumber, data.SetOutDate, data.Punctuality, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.Id)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightNumber, data.SetOutDate, data.Punctuality, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.Id)
	}, flightInfosIdKey)
}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *defaultFlightInfosModel) UpdateWithVersion(session sqlx.Session, data *FlightInfos) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	flightInfosIdKey := fmt.Sprintf("%s%v", cacheFlightInfosIdPrefix, data.Id)
	sqlResult, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, flightInfosRowsWithPlaceHolder)
		if session != nil {
			return session.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightNumber, data.SetOutDate, data.Punctuality, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.Id, oldVersion)
		}
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.FlightNumber, data.SetOutDate, data.Punctuality, data.DepartPosition, data.DepartTime, data.ArrivePosition, data.ArriveTime, data.Id, oldVersion)
	}, flightInfosIdKey)
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
func (m *defaultFlightInfosModel) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*FlightInfos, error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp FlightInfos
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *defaultFlightInfosModel) FindSum(sumBuilder squirrel.SelectBuilder) (float64, error) {

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
func (m *defaultFlightInfosModel) FindCount(countBuilder squirrel.SelectBuilder) (int64, error) {

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
func (m *defaultFlightInfosModel) FindAll(rowBuilder squirrel.SelectBuilder, orderBy string) ([]*FlightInfos, error) {

	if orderBy == "" {
		rowBuilder = rowBuilder.OrderBy("id DESC")
	} else {
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByPage 按照页码分页查询数据
func (m *defaultFlightInfosModel) FindPageListByPage(rowBuilder squirrel.SelectBuilder, page, pageSize int64, orderBy string) ([]*FlightInfos, error) {

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

	var resp []*FlightInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *defaultFlightInfosModel) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder, preMinId, pageSize int64) ([]*FlightInfos, error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? ", preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *defaultFlightInfosModel) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder, preMaxId, pageSize int64) ([]*FlightInfos, error) {

	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? ", preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *defaultFlightInfosModel) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select(flightInfosRows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *defaultFlightInfosModel) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT(" + field + ")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *defaultFlightInfosModel) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM(" + field + "),0)").From(m.table)
}

// Delete 删除数据
func (m *defaultFlightInfosModel) Delete(session sqlx.Session, id int64) error {

	flightInfosIdKey := fmt.Sprintf("%s%v", cacheFlightInfosIdPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.Exec(query, id)
		}
		return conn.Exec(query, id)
	}, flightInfosIdKey)
	return err
}

// DeleteSoft 软删除数据
func (m *defaultFlightInfosModel) DeleteSoft(session sqlx.Session, data *FlightInfos) error {
	data.DelState = globalkey.DelStateYes
	data.DeleteTime = time.Now()
	if err := m.UpdateWithVersion(session, data); err != nil {
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"), "FlightInfosModel delete err : %+v", err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *defaultFlightInfosModel) Trans(fn func(session sqlx.Session) error) error {

	err := m.Transact(func(session sqlx.Session) error {
		return fn(session)
	})
	return err

}

// formatPrimary 格式化缓存key
func (m *defaultFlightInfosModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheFlightInfosIdPrefix, primary)
}

// queryPrimary 根据主键去db查询一条数据
func (m *defaultFlightInfosModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", flightInfosRows, m.table)
	return conn.QueryRow(v, query, primary, globalkey.DelStateNo)
}

//!!!!! 其他自定义方法，从此处开始写,此处上方不要写自定义方法!!!!!

// FindListByNumber 根据航班号查询数据
func (m *defaultFlightInfosModel) FindListByNumber(rowBuilder squirrel.SelectBuilder, number string) ([]*FlightInfos, error) {

	if len(number) > 0 {
		rowBuilder = rowBuilder.Where(" flight_number = ? ", number)
	} else {
		return nil, ErrNotFound
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindListByNumberAndSetOutDate 根据航班号和出发日期查询数据
func (m *defaultFlightInfosModel) FindListByNumberAndSetOutDate(rowBuilder squirrel.SelectBuilder, number string, sod time.Time) ([]*FlightInfos, error) {

	if len(number) > 0 && !sod.IsZero() {
		sod, _ = time.Parse("2006-01-02", sod.Format("2006-01-02"))
		rowBuilder = rowBuilder.Where(" flight_number = ? AND set_out_date = ? ", number, sod.Format("2006-01-02 15:04:05"))
	} else {
		return nil, ErrNotFound
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindListBySetOutDateAndPosition 通过给定日期、出发地、目的地进行航班查询
func (m *defaultFlightInfosModel) FindListBySetOutDateAndPosition(rowBuilder squirrel.SelectBuilder, sod time.Time, depart, arrive string) ([]*FlightInfos, error) {

	if len(depart) > 0 && len(arrive) > 0 && !sod.IsZero() {
		sod, _ = time.Parse("2006-01-02", sod.Format("2006-01-02"))
		rowBuilder = rowBuilder.Where(" depart_position = ? AND arrive_position = ? AND set_out_date = ? ", depart, arrive, sod.Format("2006-01-02 15:04:05"))
	} else {
		return nil, ErrNotFound
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*FlightInfos
	err = m.QueryRowsNoCache(&resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// FindOneByByNumberAndSetOutDateAndPosition 根据 航班号 出发日期 始末时间地点 锁定一条航班信息
func (m *defaultFlightInfosModel) FindOneByByNumberAndSetOutDateAndPosition(rowBuilder squirrel.SelectBuilder, number string, sod time.Time, depart string, departTime time.Time, arrive string, arriveTime time.Time) (*FlightInfos, error) {

	if len(number) == 0 || sod.IsZero() || len(depart) == 0 || departTime.IsZero() || len(arrive) == 0 || arriveTime.IsZero() {
		return nil, ErrNotFound
	} else {
		sod, _ = time.Parse("2006-01-02", sod.Format("2006-01-02"))
		sodString := sod.Format("2006-01-02 15:04:05")
		departTimeString := strings.Split(departTime.String(), " +")[0]
		arriveTimeString := strings.Split(arriveTime.String(), " +")[0]
		rowBuilder = rowBuilder.Where("flight_number = ? AND depart_position = ? AND arrive_position = ?", number, depart, arrive)
		rowBuilder = rowBuilder.Where(squirrel.Eq{"set_out_date": sodString, "depart_time": departTimeString, "arrive_time": arriveTimeString})
	}
	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp FlightInfos
	err = m.QueryRowNoCache(&resp, query, values...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}
