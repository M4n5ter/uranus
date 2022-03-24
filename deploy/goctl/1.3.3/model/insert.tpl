// Insert 新增数据
func (m *default{{.upperStartCamelObject}}Model) Insert(session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {

	data.DeleteTime = time.Unix(0,0)
	
	{{if .withCache}}{{.keys}}
    return m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		if session != nil{
			return session.Exec(query,{{.expressionValues}})	
		}
		return conn.Exec(query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}
	query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
	if session != nil{
		return session.Exec(query,{{.expressionValues}})	
	}
	return m.conn.Exec(query, {{.expressionValues}})
	
	{{end}}
	
}
