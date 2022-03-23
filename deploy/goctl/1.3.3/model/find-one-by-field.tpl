
// FindOneBy 根据唯一索引查询一条数据，走缓存
func (m *default{{.upperStartCamelObject}}Model) FindOneBy{{.upperField}}({{.in}}) (*{{.upperStartCamelObject}}, error) {
	{{if .withCache}}{{.cacheKey}}
	var resp {{.upperStartCamelObject}}
	err := m.QueryRowIndex(&resp, {{.cacheKeyVariable}}, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where {{.originalField}} and del_state = ?  limit 1", {{.lowerStartCamelObject}}Rows, m.table)
		if err := conn.QueryRow(&resp, query, {{.lowerStartCamelField}},globalkey.DelStateNo); err != nil {
			return nil, err
		}
		return resp.{{.upperStartCamelPrimaryKey}}, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		if resp.DelState == globalkey.DelStateYes{
			return nil, ErrNotFound
		}
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{else}}var resp {{.upperStartCamelObject}}
	query := fmt.Sprintf("select %s from %s where {{.originalField}} and del_state = ? limit 1", {{.lowerStartCamelObject}}Rows, m.table )
	err := m.conn.QueryRow(&resp, query, {{.lowerStartCamelField}},globalkey.DelStateNo)
	switch err {
	case nil:
		if resp.DelState == globalkey.DelStateYes{
			return nil, ErrNotFound
		}
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}{{end}}
