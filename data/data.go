package data

var sqldb *sqldb_t
var kvdb *kvdb_t

func Open(rootPath string) error {
	var err error

	sqldb = &sqldb_t{}
	err = sqldb.open(rootPath + "/sqldb")
	if err != nil {
		return err
	}

	kvdb = &kvdb_t{}
	return kvdb.open(rootPath + "/kvdb")
}

func Close() {
	sqldb.close()
	kvdb.close()
}
