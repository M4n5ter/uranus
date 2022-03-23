// FindOne 根据主键查询一条数据，走缓存
FindOne({{.lowerStartCamelPrimaryKey}} {{.dataType}}) (*{{.upperStartCamelObject}}, error)