package data

var sqldb sqlDb

func Open(rootPath string) (err error) {
	var sqldbPath = rootPath + "/sqldb"
	sqldb, err = openSqlDb(sqldbPath)
	return
}

func Close() (err error) {
	err = sqldb.close()
	return
}
