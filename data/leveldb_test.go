package data

import (
	"os"
	"testing"
)

const KVDB_PATH = "../db_files/testkvdb"
const KVDB_TEST_CONTENT = `{"date":"2018-01-11", "items":[{"rist":"low", amount:1222}, {"rist":"high", amount:2221}]}`

var kvdbtest_db *kvdb_t

func Test_KVDB_Open(t *testing.T) {
	err := os.RemoveAll(KVDB_PATH)
	if err != nil {
		t.Error("Delete kv db folder fail.", err)
		return
	}

	kvdbtest_db = &kvdb_t{}
	err = kvdbtest_db.open(KVDB_PATH)
	if err != nil {
		t.Error(err)
	}
}

func Test_KVDB_Set(t *testing.T) {
	err := kvdbtest_db.setSnapshot(1234, 12, KVDB_TEST_CONTENT)
	if err != nil {
		t.Error(err)
	}
}

func Test_KVDB_Get(t *testing.T) {
	res, err := kvdbtest_db.getSnapshot(1234, 12)
	if err != nil {
		t.Error(err)
	} else if res != KVDB_TEST_CONTENT {
		t.Error("Result is incrrect.", res, KVDB_TEST_CONTENT)
	}
}

func Test_KVDB_Close(t *testing.T) {
	err := kvdbtest_db.close()
	if err != nil {
		t.Error(err)
	}
}
