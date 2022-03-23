// Update 修改数据 ,推荐优先使用乐观锁更新
func (m *default{{.upperStartCamelObject}}Model) Update(session sqlx.Session,data *{{.upperStartCamelObject}}) (sql.Result,error) {
        {{if .withCache}}{{.keys}}
    	return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
                query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
				if session != nil{
					return session.Exec(query, {{.expressionValues}})	
				}
                return conn.Exec(query, {{.expressionValues}})
        }, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		if session != nil{
			return session.Exec(query, {{.expressionValues}})
		}else{
			return m.conn.Exec(query, {{.expressionValues}})
		}
		{{end}}
}

// UpdateWithVersion 乐观锁修改数据 ,推荐使用
func (m *default{{.upperStartCamelObject}}Model) UpdateWithVersion(session sqlx.Session,data *{{.upperStartCamelObject}}) error {
		
		oldVersion := data.Version
		data.Version += 1

		var sqlResult sql.Result
		var err error
		
        {{if .withCache}}{{.keys}}
    	sqlResult,err =  m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
                query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} and version = ? ", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
				if session != nil{
					return session.Exec(query, {{.expressionValues}},oldVersion)	
				}
                return conn.Exec(query, {{.expressionValues}},oldVersion)
        }, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} and version = ? ", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		if session != nil{
			sqlResult,err  =  session.Exec(query, {{.expressionValues}},oldVersion)
		}else{
			sqlResult,err  =  m.conn.Exec(query, {{.expressionValues}},oldVersion)
		}
		{{end}}
		if err != nil {
			return err
		}

		updateCount , err := sqlResult.RowsAffected()
		if err != nil{
			return err
		}

		if updateCount == 0 {
			return  xerr.NewErrCode(xerr.DB_UPDATE_AFFECTED_ZERO_ERROR)
		}

		return nil

}



// FindOneByQuery 根据条件查询一条数据
func (m *default{{.upperStartCamelObject}}Model) FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}},error) {

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp {{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowNoCache(&resp, query, values...){{else}}
	err = m.conn.QueryRow(&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

// FindSum 统计某个字段总和
func (m *default{{.upperStartCamelObject}}Model) FindSum(sumBuilder squirrel.SelectBuilder) (float64,error) {

	query, values, err := sumBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp float64
	{{if .withCache}}err = m.QueryRowNoCache(&resp, query, values...){{else}}
	err = m.conn.QueryRow(&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindCount 根据某个字段查询数据数量
func (m *default{{.upperStartCamelObject}}Model) FindCount(countBuilder squirrel.SelectBuilder) (int64,error) {

	query, values, err := countBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	{{if .withCache}}err = m.QueryRowNoCache(&resp, query, values...){{else}}
	err = m.conn.QueryRow(&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindAll 查询所有数据
func (m *default{{.upperStartCamelObject}}Model) FindAll(rowBuilder squirrel.SelectBuilder,orderBy string) ([]*{{.upperStartCamelObject}},error) {

	if orderBy == ""{
		rowBuilder = rowBuilder.OrderBy("id DESC")
	}else{
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCache(&resp, query, values...){{else}}
	err = m.conn.QueryRows(&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}


// FindPageListByPage 按照页码分页查询数据
func (m *default{{.upperStartCamelObject}}Model) FindPageListByPage(rowBuilder squirrel.SelectBuilder,page ,pageSize int64,orderBy string) ([]*{{.upperStartCamelObject}},error) {

	if orderBy == ""{
		rowBuilder = rowBuilder.OrderBy("id DESC")
	}else{
		rowBuilder = rowBuilder.OrderBy(orderBy)
	}
	
	if page < 1{
		page = 1
	}
	offset := (page - 1) * pageSize

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).Offset(uint64(offset)).Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCache(&resp, query, values...){{else}}
	err = m.conn.QueryRows(&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}


// FindPageListByIdDESC 按照id倒序分页查询数据，不支持排序
func (m *default{{.upperStartCamelObject}}Model) FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder ,preMinId ,pageSize int64) ([]*{{.upperStartCamelObject}},error) {

	if preMinId > 0 {
		rowBuilder = rowBuilder.Where(" id < ? " , preMinId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id DESC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCache(&resp, query, values...){{else}}
	err = m.conn.QueryRows(&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}


// FindPageListByIdASC 按照id升序分页查询数据，不支持排序
func (m *default{{.upperStartCamelObject}}Model) FindPageListByIdASC(rowBuilder squirrel.SelectBuilder,preMaxId ,pageSize int64) ([]*{{.upperStartCamelObject}},error)  {


	if preMaxId > 0 {
		rowBuilder = rowBuilder.Where(" id > ? " , preMaxId)
	}

	query, values, err := rowBuilder.Where("del_state = ?", globalkey.DelStateNo).OrderBy("id ASC").Limit(uint64(pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	var resp []*{{.upperStartCamelObject}}
	{{if .withCache}}err = m.QueryRowsNoCache(&resp, query, values...){{else}}
	err = m.conn.QueryRows(&resp, query, values...)
	{{end}}
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

// RowBuilder 暴露给logic查询数据构建条件使用的builder
func (m *default{{.upperStartCamelObject}}Model) RowBuilder() squirrel.SelectBuilder {
	return squirrel.Select({{.lowerStartCamelObject}}Rows).From(m.table)
}

// CountBuilder 暴露给logic查询count构建条件使用的builder
func (m *default{{.upperStartCamelObject}}Model) CountBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("COUNT("+field+")").From(m.table)
}

// SumBuilder 暴露给logic查询构建条件使用的builder
func (m *default{{.upperStartCamelObject}}Model) SumBuilder(field string) squirrel.SelectBuilder {
	return squirrel.Select("IFNULL(SUM("+field+"),0)").From(m.table)
}

