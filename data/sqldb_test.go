package data

import (
	"os"
	"testing"
)

const SQDB_PATH = "../db_files/testsqldb"

var sqldbtest_db *sqldb_t
var sqldbtest_lastid int64

func Test_SQLDB_Open(t *testing.T) {
	err := os.Remove(SQDB_PATH)
	if err != nil {
		t.Error("Delete sql db file.", err)
		return
	}

	sqldbtest_db = &sqldb_t{}
	err = sqldbtest_db.open(SQDB_PATH)
	if err != nil {
		t.Error(err)
	}
}

func Test_SQLDB_Insert(t *testing.T) {
	lastID, err := sqldbtest_db.insertSnapshotVersion("admin", "2018-01-19", int64(1429))
	if err != nil {
		t.Error(err)
	} else if lastID <= 0 {
		t.Error("lastID <= 0 ", lastID)
	} else {
		sqldbtest_lastid = lastID
	}
}
func Test_SQLDB_Update(t *testing.T) {
	rowCnt, err := sqldbtest_db.updateSnapshotVersion(sqldbtest_lastid, int64(1430))
	if err != nil {
		t.Error(err)
	} else if rowCnt != 1 {
		t.Error("rowCnt != 1", rowCnt)
	}
}

func Test_SQLDB_Get_Exist(t *testing.T) {
	version, err := sqldbtest_db.getSnapshotVersion("admin", "2018-01-19")
	if err != nil {
		t.Error(err)
	} else if version != 1430 {
		t.Error("version != 1430 ", version)
	}
}

func Test_SQLDB_Get_Not_Exist(t *testing.T) {
	version, err := sqldbtest_db.getSnapshotVersion("admin1", "2018-01-19")
	if err == nil {
		t.Error("version=", version)
	}
}

func Test_SQLDB_Delete(t *testing.T) {
	_, err := sqldbtest_db.deleteSnapshotVersion(sqldbtest_lastid)
	if err != nil {
		t.Error(err)
		return
	}

	var version int64
	version, err = sqldbtest_db.getSnapshotVersion("admin", "2018-01-19")
	if err == nil || version != 0 {
		t.Error("Delete fail, data exist yet. lastid, version=", sqldbtest_lastid, version)
	}
}

func Test_SQLDB_Close(t *testing.T) {
	err := sqldbtest_db.close()
	if err != nil {
		t.Error(err)
	}
}
