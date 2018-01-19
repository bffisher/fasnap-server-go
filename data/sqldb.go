package data

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const SQL_CREATE_SNAPSHOT_VERSION = "CREATE TABLE IF NOT EXISTS snapshot_version (id INTEGER PRIMARY KEY, user TEXT, date TEXT, version INTEGER, UNIQUE(user, date))"
const SQL_SELECT_SNAPSHOT_VERSION = "SELECT version FROM snapshot_version WHERE user = ? and date = ?"
const SQL_UPDATE_SNAPSHOT_VERSION = "UPDATE snapshot_version SET version = ? WHERE id = ?"
const SQL_INSERT_SNAPSHOT_VERSION = "INSERT INTO snapshot_version (user, date, version) VALUES(?, ?, ?)"
const SQL_DELETE_SNAPSHOT_VERSION = "DELETE FROM snapshot_version WHERE id = ?"

type sqlDb struct {
	impl *sql.DB
}

func openSqlDb(path string) (db sqlDb, err error) {
	db.impl, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Println(err)
		return
	}

	//Create table
	_, _, err = db.exec(SQL_CREATE_SNAPSHOT_VERSION)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *sqlDb) close() (err error) {
	err = db.impl.Close()
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *sqlDb) getSnapshotVersion(user, date string) (version int64, err error) {
	err = db.impl.QueryRow(SQL_SELECT_SNAPSHOT_VERSION, user, date).Scan(&version)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *sqlDb) insertSnapshotVersion(user, date string, version int64) (id int64, err error) {
	id, _, err = db.exec(SQL_INSERT_SNAPSHOT_VERSION, user, date, version)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *sqlDb) updateSnapshotVersion(id, version int64) (rowCnt int64, err error) {
	_, rowCnt, err = db.exec(SQL_UPDATE_SNAPSHOT_VERSION, version, id)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *sqlDb) deleteSnapshotVersion(id int64) (rowCnt int64, err error) {
	_, rowCnt, err = db.exec(SQL_DELETE_SNAPSHOT_VERSION, id)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *sqlDb) exec(sqlText string, args ...interface{}) (lastId, rowCnt int64, err error) {
	var stmt *sql.Stmt
	stmt, err = db.impl.Prepare(sqlText)
	if err != nil {
		log.Println(err)
		return
	}

	var res sql.Result
	res, err = stmt.Exec(args...)
	if err != nil {
		log.Println(err)
		return
	}

	lastId, err = res.LastInsertId()
	if err != nil {
		log.Println(err)
		return
	}

	rowCnt, err = res.RowsAffected()
	if err != nil {
		log.Println(err)
	}

	return
}
