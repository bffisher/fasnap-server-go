package data

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/golang/leveldb"
)

const _KEY_PREFIXS_SYSTEM = "0"
const _KEY_PREFIXS_SNAPSHOT = "1"

var _KEY_VERSION = []byte(_KEY_PREFIXS_SYSTEM + "_data_version")

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

func (db *kvdb_t) getVersion() (int64, error) {
	res, err := db.impl.Get(_KEY_VERSION, nil)
	if err != nil {
		return 0, err
	}

	return bytesToInt64(res), nil
}

func (db *kvdb_t) setVersion(version int64) error {
	return db.impl.Set(_KEY_VERSION, toBytes(version), nil)
}

func (db *kvdb_t) getSnapshot(id, version int64) (string, error) {
	key := kvdb_key_t{id, version}
	res, err := db.impl.Get(toBytes(key), nil)

	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (db *kvdb_t) setSnapshot(id, version int64, content string) error {
	key := kvdb_key_t{id, version}
	return db.impl.Set(toBytes(key), []byte(content), nil)
}

func toBytes(val interface{}) []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, val)
	return buf.Bytes()
}

func bytesToInt64(val []byte) int64 {
	var res int64
	buf := bytes.NewBuffer(val)
	binary.Read(buf, binary.BigEndian, &res)
	return res
}
