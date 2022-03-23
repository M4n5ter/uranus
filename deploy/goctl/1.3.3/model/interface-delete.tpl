// Delete 删除数据
Delete(session sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error