package data

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/golang/leveldb"
)

type kvdb_t struct {
	impl *leveldb.DB
}

type kvdb_key_t struct {
	id, version int64
}

func (db *kvdb_t) open(path string) (err error) {
	if db == nil {
		return errors.New("db is nil")
	}
	if db.impl != nil {
		db.close()
	}

	db.impl, err = leveldb.Open(path, nil)
	return
}

func (db *kvdb_t) close() error {
	if db.impl == nil {
		return nil
	}

	return db.impl.Close()
}

func (db *kvdb_t) getSnapshot(id, version int64) (string, error) {
	key := kvdb_key_t{id, version}

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, key)

	res, err := db.impl.Get(buf.Bytes(), nil)

	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (db *kvdb_t) setSnapshot(id, version int64, content string) error {
	key := kvdb_key_t{id, version}

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, key)

	return db.impl.Set(buf.Bytes(), []byte(content), nil)
}
