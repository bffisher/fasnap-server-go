package db

import (
	"database/sql"
	"errors"
	"log"

	//Need to import sqlite3
	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlCreateSnapshotVer = "CREATE TABLE IF NOT EXISTS snapshot_version (id INTEGER PRIMARY KEY, user TEXT, date TEXT, version INTEGER, UNIQUE(user, date))"
	sqlSelectSnapshotVer = "SELECT version FROM snapshot_version WHERE user = ? and date = ?"
	sqlUpdateSnapshotVer = "UPDATE snapshot_version SET version = ? WHERE id = ?"
	sqlInsertSnapshotVer = "INSERT INTO snapshot_version (user, date, version) VALUES(?, ?, ?)"
	sqlDeleteSnapshotVer = "DELETE FROM snapshot_version WHERE id = ?"
)

type SQLDB struct {
	impl *sql.DB
}

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
	_, _, err = db.exec(sqlCreateSnapshotVer)
	return err
}

func (db *SQLDB) Close() (err error) {
	if db.impl == nil {
		return nil
	}
	return db.impl.Close()
}

func (db *SQLDB) GetSnapshotVersion(user, date string) (version int64, err error) {
	err = db.impl.QueryRow(sqlSelectSnapshotVer, user, date).Scan(&version)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *SQLDB) InsertSnapshotVersion(user, date string, version int64) (id int64, err error) {
	id, _, err = db.exec(sqlInsertSnapshotVer, user, date, version)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *SQLDB) UpdateSnapshotVersion(id, version int64) (rowCnt int64, err error) {
	_, rowCnt, err = db.exec(sqlUpdateSnapshotVer, version, id)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *SQLDB) DeleteSnapshotVersion(id int64) (rowCnt int64, err error) {
	_, rowCnt, err = db.exec(sqlDeleteSnapshotVer, id)
	if err != nil {
		log.Println(err)
	}
	return
}

func (db *SQLDB) exec(sqlText string, args ...interface{}) (lastID, rowCnt int64, err error) {
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

	lastID, err = res.LastInsertId()
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
