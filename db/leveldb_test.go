package db

import (
	"os"
	"testing"
)

const (
	kvdbTestDBPath  = "../db_files/testkvdb"
	kvdbTestContent = `{"date":"2018-01-11", "items":[{"rist":"low", amount:1222}, {"rist":"high", amount:2221}]}`
)

var kvdbTestDB *KVDB

func Test_KVDB_Open(t *testing.T) {
	err := os.RemoveAll(kvdbTestDBPath)
	if err != nil && !os.IsNotExist(err) {
		t.Error("Delete kv db folder fail.", err)
		return
	}

	kvdbTestDB = &KVDB{}
	err = kvdbTestDB.Open(kvdbTestDBPath)
	if err != nil {
		t.Error(err)
	}
}

func Test_KVDB_Set(t *testing.T) {
	err := kvdbTestDB.SetSnapshot(1234, kvdbTestContent)
	if err != nil {
		t.Error(err)
	}
}

func Test_KVDB_Get(t *testing.T) {
	snap, err := kvdbTestDB.GetSnapshot(1234)
	if err != nil {
		t.Error(err)
	} else if snap != kvdbTestContent {
		t.Error("Result(snapshot) is incrrect.", snap, kvdbTestContent)
	}
}

func Test_KVDB_Close(t *testing.T) {
	err := kvdbTestDB.Close()
	if err != nil {
		t.Error(err)
	}
}
