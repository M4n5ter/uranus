
type (
	{{.upperStartCamelObject}}Model interface{
		// Insert 新增数据
		Insert(session sqlx.Session,data *{{.upperStartCamelObject}}) (sql.Result,error)
		{{.method}}
		// DeleteSoft 软删除数据
		DeleteSoft(session sqlx.Session, data *{{.upperStartCamelObject}}) error
		// Update 更新数据
		Update(session sqlx.Session,data *{{.upperStartCamelObject}})  (sql.Result, error)
		// UpdateWithVersion 更新数据，使用乐观锁
		UpdateWithVersion(session sqlx.Session,data *{{.upperStartCamelObject}}) error
		// FindOneByQuery 根据条件查询一条数据，不走缓存
		FindOneByQuery(rowBuilder squirrel.SelectBuilder) (*{{.upperStartCamelObject}},error) 
		// FindSum sum某个字段
		FindSum(sumBuilder squirrel.SelectBuilder) (float64,error) 
		// FindCount 根据条件统计条数
		FindCount(countBuilder squirrel.SelectBuilder) (int64,error) 
		// FindAll 查询所有数据不分页
		FindAll(rowBuilder squirrel.SelectBuilder,orderBy string) ([]*{{.upperStartCamelObject}},error) 
		// FindPageListByPage 根据页码分页查询分页数据
		FindPageListByPage(rowBuilder squirrel.SelectBuilder,page ,pageSize int64,orderBy string) ([]*{{.upperStartCamelObject}},error) 
		// FindPageListByIdDESC 根据id倒序分页查询分页数据
		FindPageListByIdDESC(rowBuilder squirrel.SelectBuilder ,preMinId ,pageSize int64) ([]*{{.upperStartCamelObject}},error)
		// FindPageListByIdASC 根据id升序分页查询分页数据
		FindPageListByIdASC(rowBuilder squirrel.SelectBuilder,preMaxId ,pageSize int64) ([]*{{.upperStartCamelObject}},error)  
		// Trans 暴露给logic，开启事务
		Trans(fn func(session sqlx.Session) error) error
		// RowBuilder 暴露给logic，查询数据的builder
		RowBuilder() squirrel.SelectBuilder
		// CountBuilder 暴露给logic，查询count的builder
		CountBuilder(field string) squirrel.SelectBuilder
		// SumBuilder 暴露给logic，查询sum的builder
		SumBuilder(field string) squirrel.SelectBuilder
	}

	default{{.upperStartCamelObject}}Model struct {
		{{if .withCache}}sqlc.CachedConn{{else}}conn sqlx.SqlConn{{end}}
		table string
	}

	{{.upperStartCamelObject}} struct {
		{{.fields}}
	}
)
