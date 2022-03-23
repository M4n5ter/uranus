package casbinTools

type Conf struct {
	DB struct {
		DataSourceWithoutDBName string
		DBName                  string
	}
	Model string
}
