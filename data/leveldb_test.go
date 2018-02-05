package data

import (
	"os"
	"testing"
)

const (
	kvdbTestDBPath  = "../db_files/testkvdb"
	kvdbTestContent = `{"date":"2018-01-11", "items":[{"rist":"low", amount:1222}, {"rist":"high", amount:2221}]}`
)

var kvdbTestDB *kvdbType

func Test_KVDB_Open(t *testing.T) {
	err := os.RemoveAll(kvdbTestDBPath)
	if err != nil && !os.IsNotExist(err) {
		t.Error("Delete kv db folder fail.", err)
		return
	}

	kvdbTestDB = &kvdbType{}
	err = kvdbTestDB.open(kvdbTestDBPath)
	if err != nil {
		t.Error(err)
	}
}

func Test_KVDB_Set(t *testing.T) {
	err := kvdbTestDB.setVersion(23)
	if err != nil {
		t.Error(err)
	}

	err = kvdbTestDB.setSnapshot(1234, 12, kvdbTestContent)
	if err != nil {
		t.Error(err)
	}
}

func Test_KVDB_Get(t *testing.T) {
	ver, err := kvdbTestDB.getVersion()
	if err != nil {
		t.Error(err)
	} else if ver != 23 {
		t.Error("Result(version) is incrrect.", ver, 22)
	}

	snap, err := kvdbTestDB.getSnapshot(1234, 12)
	if err != nil {
		t.Error(err)
	} else if snap != kvdbTestContent {
		t.Error("Result(snapshot) is incrrect.", snap, kvdbTestContent)
	}
}

func Test_KVDB_Close(t *testing.T) {
	err := kvdbTestDB.close()
	if err != nil {
		t.Error(err)
	}
}
