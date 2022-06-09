package commonModel

import (
	"database/sql"
	"fmt"
	"github.com/duke-git/lancet/v2/datetime"
	"github.com/zeromicro/go-zero/core/mr"
	"strings"
	"time"
	"uranus/common/timeTools"

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
		// FindPageListByNumberAndDays 根据航班号查询数据
		FindPageListByNumberAndDays(rowBuilder squirrel.SelectBuilder, number string, days, limit int64) ([]*FlightInfos, error)
		// FindListByNumberAndSetOutDate 根据航班号和出发日期查询数据
		FindListByNumberAndSetOutDate(rowBuilder squirrel.SelectBuilder, number string, sot time.Time) ([]*FlightInfos, error)
		// FindListBySetOutDateAndPosition 通过给定日期、出发地、目的地进行航班查询
		FindListBySetOutDateAndPosition(rowBuilder squirrel.SelectBuilder, sod time.Time, depart, arrive string) ([]*FlightInfos, error)
		// FindOneByByNumberAndSetOutDateAndPosition 根据 航班号 出发日期 始末时间地点 锁定一条航班信息
		FindOneByByNumberAndSetOutDateAndPosition(rowBuilder squirrel.SelectBuilder, number string, sod time.Time, depart string, departTime time.Time, arrive string, arriveTime time.Time) (*FlightInfos, error)
		// FindPageListByPositionAndDays 根据起始地点和距今日的日期差分页查询, limit <= 0 表示不分页
		FindPageListByPositionAndDays(rowBuilder squirrel.SelectBuilder, departPosition, arrivePosition string, days int64, limit int64) ([]*FlightInfos, error)
		// FindTransferFlightsByPlace 根据出发地、目的地、出发日期查询中转航班(中转航班)
		FindTransferFlightsByPlace(rowBuilder squirrel.SelectBuilder, departPosition, arrivePosition string, sod time.Time) ([]*Transfer, error)
		// FindPageListByPositionSODAndDays 根据起始地点和距今日的日期差分页查询, limit <= 0 表示不分页
		FindPageListByPositionSODAndDays(rowBuilder squirrel.SelectBuilder, departPosition, arrivePosition string, selectedDate time.Time, days int64, limit int64) ([]*FlightInfos, error)
	}

	defaultFlightInfosModel struct {
		sqlc.CachedConn
		table string
	}

	FlightInfos struct {
		Id             int64     `db:"id" json:"id"`
		CreateTime     time.Time `db:"create_time" json:"createTime"`
		UpdateTime     time.Time `db:"update_time" json:"updateTime"`
		DeleteTime     time.Time `db:"delete_time" json:"deleteTime"`
		DelState       int64     `db:"del_state" json:"delState"`             // 是否已经删除
		Version        int64     `db:"version" json:"version"`                // 版本号
		FlightNumber   string    `db:"flight_number" json:"flightNumber"`     // 对应的航班号
		SetOutDate     time.Time `db:"set_out_date" json:"setOutDate"`        // 出发日期
		Punctuality    int64     `db:"punctuality" json:"punctuality"`        // 准点率(%)
		DepartPosition string    `db:"depart_position" json:"departPosition"` // 起飞地点
		DepartTime     time.Time `db:"depart_time" json:"departTime"`         // 起飞时间
		ArrivePosition string    `db:"arrive_position" json:"arrivePosition"` // 降落地点
		ArriveTime     time.Time `db:"arrive_time" json:"arriveTime"`         // 降落时间
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

// FindPageListByNumberAndDays 根据航班号查询数据
func (m *defaultFlightInfosModel) FindPageListByNumberAndDays(rowBuilder squirrel.SelectBuilder, number string, days, limit int64) ([]*FlightInfos, error) {

	if len(number) > 0 {
		rowBuilder = rowBuilder.Where(" flight_number = ? ", number)
	} else {
		return nil, ErrNotFound
	}

	today := time.Now()
	today, _ = timeTools.Time2TimeYMD000(today)
	todayString := today.Format("2006-01-02 15:04:05")
	lastDate := today.AddDate(0, 0, int(days))
	lastDateString := lastDate.Format("2006-01-02 15:04:05")
	if days < 0 {
		rowBuilder = rowBuilder.Where("set_out_date between ? and ?", lastDateString, todayString)
	} else if days > 0 {
		rowBuilder = rowBuilder.Where("set_out_date between ? and ?", todayString, lastDateString)
	} else {
		rowBuilder = rowBuilder.Where(squirrel.Eq{"set_out_date": todayString})
	}

	if limit > 0 {
		rowBuilder = rowBuilder.Limit(uint64(limit))
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
		sod, _ = timeTools.Time2TimeYMD000(sod)
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
		sod, _ = timeTools.Time2TimeYMD000(sod)
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
		sod, _ = timeTools.Time2TimeYMD000(sod)
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

// FindPageListByPositionAndDays 根据起始地点和距今日的日期差分页查询, limit <= 0 表示不分页
func (m *defaultFlightInfosModel) FindPageListByPositionAndDays(rowBuilder squirrel.SelectBuilder, departPosition, arrivePosition string, days int64, limit int64) ([]*FlightInfos, error) {
	if len(departPosition) == 0 || len(arrivePosition) == 0 {
		return nil, ErrNotFound
	} else {
		today := time.Now()
		today, _ = timeTools.Time2TimeYMD000(today)
		todayString := today.Format("2006-01-02 15:04:05")
		lastDate := today.AddDate(0, 0, int(days))
		lastDateString := lastDate.Format("2006-01-02 15:04:05")
		rowBuilder = rowBuilder.Where(squirrel.Eq{"depart_position": departPosition, "arrive_position": arrivePosition})
		if days < 0 {
			rowBuilder = rowBuilder.Where("set_out_date between ? and ?", lastDateString, todayString)
		} else if days > 0 {
			rowBuilder = rowBuilder.Where("set_out_date between ? and ?", todayString, lastDateString)
		} else {
			rowBuilder = rowBuilder.Where(squirrel.Eq{"set_out_date": todayString})
		}

		if limit > 0 {
			rowBuilder = rowBuilder.Limit(uint64(limit))
		}
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

// FindPageListByPositionSODAndDays 根据起始地点和距选定的日期差分页查询, limit <= 0 表示不分页
func (m *defaultFlightInfosModel) FindPageListByPositionSODAndDays(rowBuilder squirrel.SelectBuilder, departPosition, arrivePosition string, selectedDate time.Time, days int64, limit int64) ([]*FlightInfos, error) {
	if len(departPosition) == 0 || len(arrivePosition) == 0 {
		return nil, ErrNotFound
	} else {
		sod, _ := timeTools.Time2TimeYMD000(selectedDate)
		sodString := sod.Format("2006-01-02 15:04:05")
		lastDate := sod.AddDate(0, 0, int(days))
		lastDateString := lastDate.Format("2006-01-02 15:04:05")
		rowBuilder = rowBuilder.Where(squirrel.Eq{"depart_position": departPosition, "arrive_position": arrivePosition})
		if days < 0 {
			rowBuilder = rowBuilder.Where("set_out_date between ? and ?", lastDateString, sodString)
		} else if days > 0 {
			rowBuilder = rowBuilder.Where("set_out_date between ? and ?", sodString, lastDateString)
		} else {
			rowBuilder = rowBuilder.Where(squirrel.Eq{"set_out_date": sodString})
		}

		if limit > 0 {
			rowBuilder = rowBuilder.Limit(uint64(limit))
		}
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

// FindTransferFlightsByPlace 根据出发地、目的地、出发日期查询中转航班(中转航班)
func (m *defaultFlightInfosModel) FindTransferFlightsByPlace(rowBuilder squirrel.SelectBuilder, departPosition, arrivePosition string, sod time.Time) ([]*Transfer, error) {
	if len(departPosition) == 0 || len(arrivePosition) == 0 || sod.IsZero() {
		return nil, ErrNotFound
	}

	sod, _ = timeTools.Time2TimeYMD000(sod)

	//符合出发地要求
	var departSlice []*FlightInfos
	//符合目的地要求
	var arriveSlice []*FlightInfos
	err := mr.Finish(func() (err error) {
		// 构建出发地点是 departPosition 的 sql
		rowBuilder1 := rowBuilder.Where(squirrel.Eq{"set_out_date": sod.Format("2006-01-02 15:04:05")})

		departRowBuilder := rowBuilder1.Where(squirrel.Eq{"depart_position": departPosition})

		departQuery, departValues, err := departRowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
		if err != nil {
			return err
		}

		err = m.QueryRowsNoCache(&departSlice, departQuery, departValues...)
		return
	}, func() (err error) {
		// 构建出发地点是 arrivePosition 的 sql
		// 中转航班允许跨一天
		lateSOD := sod.AddDate(0, 0, 1)
		rowBuilder2 := rowBuilder.Where("(set_out_date = ? or set_out_date = ?)", sod.Format("2006-01-02 15:04:05"), lateSOD.Format("2006-01-02 15:04:05"))
		arriveRowBuilder := rowBuilder2.Where(squirrel.Eq{"arrive_position": arrivePosition})

		arriveQuery, arriveValues, err := arriveRowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
		if err != nil {
			return err
		}

		err = m.QueryRowsNoCache(&arriveSlice, arriveQuery, arriveValues...)
		if err != nil {
			return err
		}

		return
	})

	if err != nil {
		return nil, err
	}

	if departSlice == nil || arriveSlice == nil {
		return nil, ErrNotFound
	}

	// 存放 departSlice 中的 arrivePosition 与 arriveTime 的对应关系
	departMap := make(map[string]struct {
		t time.Time
		f *FlightInfos
	})
	// 存放 arriveSlice 中的 departPosition 与 departTime 的对应关系
	arriveMap := make(map[string]struct {
		t time.Time
		f *FlightInfos
	})

	for _, infos := range departSlice {
		// 过滤数据库中可能的错误记录
		if infos.DepartTime.After(infos.ArriveTime) {
			continue
		}

		departMap[infos.ArrivePosition] = struct {
			t time.Time
			f *FlightInfos
		}{t: infos.ArriveTime, f: infos}
	}
	for _, infos := range arriveSlice {
		// 过滤数据库中可能的错误记录
		if infos.DepartTime.After(infos.ArriveTime) {
			continue
		}

		arriveMap[infos.DepartPosition] = struct {
			t time.Time
			f *FlightInfos
		}{t: infos.ArriveTime, f: infos}
	}

	var transfers []*Transfer
	for depo, deti := range departMap {
		if arti, exist := arriveMap[depo]; exist {
			// 存在中转地点 depo, 比较到达中转地点的时间和终端地点的起飞时间(至少需要预留 1 个小时)
			reservedTime := datetime.AddHour(deti.t, 1)
			if reservedTime.Equal(arti.t) || reservedTime.Before(arti.t) && arti.t.Sub(reservedTime).Hours() < 24 {
				// 聚合
				transfers = append(transfers, &Transfer{[]*FlightInfos{departMap[depo].f, arriveMap[depo].f}})
			}
		}
	}
	return transfers, nil
}
