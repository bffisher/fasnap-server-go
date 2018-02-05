package data

var sqldb *sqldbType
var kvdb *kvdbType

//Open and init database
func Open(rootPath string) error {
	var err error

	sqldb = &sqldbType{}
	err = sqldb.open(rootPath + "/sqldb")
	if err != nil {
		return err
	}

	kvdb = &kvdbType{}
	return kvdb.open(rootPath + "/kvdb")
}

//Close database
func Close() {
	sqldb.close()
	kvdb.close()
}
