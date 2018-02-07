package db

import (
	"database/sql"
	"errors"
	"fasnap-server-go/entities"

	//Need to import sqlite3
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlCreateSnapshot = `CREATE TABLE IF NOT EXISTS snapshot (ver INTEGER PRIMARY KEY, user TEXT, date TEXT, operation INTEGER); 
											 CREATE INDEX IF NOT EXISTS snapshot_user_index ON snapshot (user)`
	sqlGetSnapshotVer  = `SELECT ver FROM snapshot WHERE user = ? ORDER BY ver DESC LIMIT 1`
	sqlGetSnapshot     = `SELECT ver, operation FROM snapshot WHERE user = ? and date = ? ORDER BY ver DESC LIMIT 1`
	sqlGetSnapshotList = `SELECT ver, date, operation FROM snapshot WHERE user = ? and ver > ?`
	sqlInsertSnapshot  = `INSERT INTO snapshot (user, date, operation) VALUES(?, ?, ?)`
)

//SQLDB Package sqlite3 instance
type SQLDB struct {
	impl *sql.DB
}

//Open Open sqlite3 database
func (db *SQLDB) Open(path string) error {
	if db == nil {
		return errors.New("db is nil")
	}
	if db.impl != nil {
		db.Close()
	}

	var err error
	db.impl, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	//Create table
	_, _, err = db.exec(sqlCreateSnapshot)
	return err
}

//Close Close sqlite3 database
func (db *SQLDB) Close() (err error) {
	if db.impl == nil {
		return nil
	}
	return db.impl.Close()
}

//GetCurSnapshotVersion Return current version of snapshot for user
func (db *SQLDB) GetCurSnapshotVersion(user string) (int64, error) {
	rows, err := db.impl.Query(sqlGetSnapshotVer, user)
	if err != nil {
		return -1, err
	}

	if rows.Next() {
		var ver int64
		rows.Scan(&ver)
		return ver, nil
	}

	return 0, nil
}

//GetSnapshot Get a snapshot by user and date, Return snapshot
func (db *SQLDB) GetSnapshot(user, date string) (snapshot entities.Snapshot, err error) {
	err = db.impl.QueryRow(sqlGetSnapshot, user, date).Scan(&snapshot.Version, &snapshot.Operation)
	return
}

//GetSnapshotList Get some snapshots by user and min version, Return [id, date, version, operation] list
func (db *SQLDB) GetSnapshotList(user string, minVersion int64) ([]entities.Snapshot, error) {
	snapshots := make([]entities.Snapshot, 0)
	rows, err := db.impl.Query(sqlGetSnapshotList, user, minVersion)
	if err != nil {
		return nil, err
	}

	for i := 0; rows.Next(); i++ {
		snapshots = append(snapshots, entities.Snapshot{})
		err = rows.Scan(&snapshots[i].Version, &snapshots[i].Date, &snapshots[i].Operation)
		if err != nil {
			return nil, err
		}
	}
	return snapshots, nil
}

//SaveSnapshot Save a snapshot
func (db *SQLDB) SaveSnapshot(user, date string) (id int64, err error) {
	id, _, err = db.exec(sqlInsertSnapshot, user, date, entities.OPSSave)
	return
}

//DeleteSnapshot Delete a snapshot(Save a snapshot which operation is set DEL)
func (db *SQLDB) DeleteSnapshot(user, date string) (id int64, err error) {
	id, _, err = db.exec(sqlInsertSnapshot, user, date, entities.OPSDel)
	return
}

func (db *SQLDB) exec(sqlText string, args ...interface{}) (lastID, rowCnt int64, err error) {
	var stmt *sql.Stmt
	stmt, err = db.impl.Prepare(sqlText)
	if err != nil {
		return
	}

	var res sql.Result
	res, err = stmt.Exec(args...)
	if err != nil {
		return
	}

	lastID, err = res.LastInsertId()
	if err != nil {
		return
	}

	rowCnt, err = res.RowsAffected()
	return
}
