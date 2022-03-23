// Delete 删除数据
func (m *default{{.upperStartCamelObject}}Model) Delete(session sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne({{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}{{end}}

	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
		if session!=nil{
			return session.Exec(query, {{.lowerStartCamelPrimaryKey}})	
		}
		return conn.Exec(query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
		if session!=nil{
			_,err:= session.Exec(query, {{.lowerStartCamelPrimaryKey}})	
			return err
		}
		_,err:=m.conn.Exec(query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}

// DeleteSoft 软删除数据
func (m *default{{.upperStartCamelObject}}Model) DeleteSoft(session sqlx.Session,data *{{.upperStartCamelObject}}) error {
	data.DelState = globalkey.DelStateYes
	data.DeletedAt.Time = time.Now()
	if err:= m.UpdateWithVersion(session, data);err!= nil{
		return errors.Wrapf(xerr.NewErrMsg("删除数据失败"),"{{.upperStartCamelObject}}Model delete err : %+v",err)
	}
	return nil
}

// Trans 暴露给logic开启事务
func (m *default{{.upperStartCamelObject}}Model) Trans(fn func(session sqlx.Session) error) error {
	{{if .withCache}}
		err := m.Transact(func(session sqlx.Session) error {
			return  fn(session)
		})
		return err
	{{else}}
		err := m.conn.Transact(func(session sqlx.Session) error {
			return  fn(session)
		})
		return err
	{{end}}
}


