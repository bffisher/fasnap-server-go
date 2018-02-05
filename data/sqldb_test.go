package data

import (
	"os"
	"testing"
)

const sqdbTestPath = "../db_files/testsqldb"

var sqldbTestDB *sqldbType
var sqldbtestLastID int64

func Test_SQLDB_Open(t *testing.T) {
	err := os.Remove(sqdbTestPath)
	if err != nil && !os.IsNotExist(err) {
		t.Error("Delete sql db file.", err)
		return
	}

	sqldbTestDB = &sqldbType{}
	err = sqldbTestDB.open(sqdbTestPath)
	if err != nil {
		t.Error(err)
	}
}

func Test_SQLDB_Insert(t *testing.T) {
	lastID, err := sqldbTestDB.insertSnapshotVersion("admin", "2018-01-19", int64(1429))
	if err != nil {
		t.Error(err)
	} else if lastID <= 0 {
		t.Error("lastID <= 0 ", lastID)
	} else {
		sqldbtestLastID = lastID
	}
}
func Test_SQLDB_Update(t *testing.T) {
	rowCnt, err := sqldbTestDB.updateSnapshotVersion(sqldbtestLastID, int64(1430))
	if err != nil {
		t.Error(err)
	} else if rowCnt != 1 {
		t.Error("rowCnt != 1", rowCnt)
	}
}

func Test_SQLDB_Get_Exist(t *testing.T) {
	version, err := sqldbTestDB.getSnapshotVersion("admin", "2018-01-19")
	if err != nil {
		t.Error(err)
	} else if version != 1430 {
		t.Error("version != 1430 ", version)
	}
}

func Test_SQLDB_Get_Not_Exist(t *testing.T) {
	version, err := sqldbTestDB.getSnapshotVersion("admin1", "2018-01-19")
	if err == nil {
		t.Error("version=", version)
	}
}

func Test_SQLDB_Delete(t *testing.T) {
	_, err := sqldbTestDB.deleteSnapshotVersion(sqldbtestLastID)
	if err != nil {
		t.Error(err)
		return
	}

	var version int64
	version, err = sqldbTestDB.getSnapshotVersion("admin", "2018-01-19")
	if err == nil || version != 0 {
		t.Error("Delete fail, data exist yet. lastid, version=", sqldbtestLastID, version)
	}
}

func Test_SQLDB_Close(t *testing.T) {
	err := sqldbTestDB.close()
	if err != nil {
		t.Error(err)
	}
}
