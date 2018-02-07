package db

import (
	"fasnap-server-go/entities"
	"log"
	"os"
	"testing"
)

const sqdbTestPath = "../db_files/testsqldb"

var sqldbTestDB *SQLDB
var sqldbtestLastID int64

func Test_SQLDB_Open(t *testing.T) {
	err := os.Remove(sqdbTestPath)
	if err != nil && !os.IsNotExist(err) {
		t.Error("Delete sql db file.", err)
		return
	}

	sqldbTestDB = &SQLDB{}
	err = sqldbTestDB.Open(sqdbTestPath)
	if err != nil {
		t.Error(err)
	}
}

func Test_SQLDB_Insert(t *testing.T) {
	lastID, err := sqldbTestDB.SaveSnapshot("admin", "2018-01-19")
	if err != nil {
		t.Error(err)
	} else if lastID <= 0 {
		t.Error("lastID <= 0 ", lastID)
	} else {
		sqldbtestLastID = lastID
		log.Println("lastID", lastID)
	}
}
func Test_SQLDB_Update(t *testing.T) {
	lastID, err := sqldbTestDB.SaveSnapshot("admin", "2018-01-19")
	if err != nil {
		t.Error(err)
	} else if lastID <= 0 {
		t.Error("lastID <= 0 ", lastID)
	} else {
		sqldbtestLastID = lastID
		log.Println("lastID", lastID)
	}
}

func Test_SQLDB_Get_Exist(t *testing.T) {
	snapshots, err := sqldbTestDB.GetSnapshotList("admin", 0)
	if err != nil {
		t.Error(err)
	} else if l := len(snapshots); l != 2 {
		t.Error("len != 2 ", l)
	}
}

func Test_SQLDB_Delete(t *testing.T) {
	lastID, err := sqldbTestDB.DeleteSnapshot("admin", "2018-01-19")
	if err != nil {
		t.Error(err)
		return
	} else if lastID <= 0 {
		t.Error("lastID <= 0 ", lastID)
	} else {
		sqldbtestLastID = lastID
		log.Println("lastID", lastID)
	}

	snapshot, err := sqldbTestDB.GetSnapshot("admin", "2018-01-19")
	if err != nil {
		t.Error(err)
		return
	}

	if snapshot.Operation != entities.OPSDel {
		t.Error("Delete fail, ", snapshot)
	}
}

func Test_SQLDB_Close(t *testing.T) {
	err := sqldbTestDB.Close()
	if err != nil {
		t.Error(err)
	}
}
