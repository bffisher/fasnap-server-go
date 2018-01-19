package data

import (
	"os"
	"testing"
)

const DB_PATH = "./testsqldb"

func openDB() (db sqlDb, err error) {
	os.Remove(DB_PATH)
	db, err = openSqlDb(DB_PATH)
	return
}

func closeDB(db sqlDb) {
	db.close()
	os.Remove(DB_PATH)
}

func Test_Snapshot(t *testing.T) {
	db, err := openDB()
	defer closeDB(db)

	if err != nil {
		t.Error("Open fail", err)
		return
	}

	t.Log("Open sql db sucessfully.")

	var lastID, rowCnt int64
	lastID, err = db.insertSnapshotVersion("admin", "2018-01-19", int64(1429))
	if err != nil {
		t.Error("Insert fail", err)
		return
	} else if lastID <= 0 {
		t.Error("Insert fail, lastID = ", lastID)
		return
	}

	t.Log("Insert sucessfully.")

	rowCnt, err = db.updateSnapshotVersion(lastID, int64(1430))
	if err != nil {
		t.Error("Update fail", err)
		return
	} else if rowCnt != 1 {
		t.Error("Update fail, rowCnt = ", rowCnt)
		return
	}

	t.Log("Update sucessfully.")

	var version int64
	version, err = db.getSnapshotVersion("admin", "2018-01-19")
	if err != nil {
		t.Error("Get fail(exist)", err)
		return
	} else if version != 1430 {
		t.Error("Get fail(exist), version != 1430, version is ", version)
		return
	}

	version, err = db.getSnapshotVersion("admin1", "2018-01-19")
	if err == nil {
		t.Error("Get fail(not exist) version=", version)
		return
	}

	t.Log("Get sucessfully.")

	rowCnt, err = db.deleteSnapshotVersion(lastID)
	if err != nil {
		t.Error("Delete fail", err)
		return
	}

	version, err = db.getSnapshotVersion("admin", "2018-01-19")
	if err == nil || version != 0 {
		t.Error("Delete fail")
		return
	}

	t.Log("Delete sucessfully.")
}
