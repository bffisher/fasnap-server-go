package db

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/golang/leveldb"
)

type KVDB struct {
	impl *leveldb.DB
}

const (
	sysKeyCategory      byte = 0
	snapshotKeyCategory byte = 1
)

var versionKey = []byte{sysKeyCategory, 0}

func (db *KVDB) Open(path string) (err error) {
	if db == nil {
		return errors.New("db is nil")
	}
	if db.impl != nil {
		db.Close()
	}

	db.impl, err = leveldb.Open(path, nil)
	return
}

func (db *KVDB) Close() error {
	if db.impl == nil {
		return nil
	}

	return db.impl.Close()
}

func (db *KVDB) GetDataVersion() (uint64, error) {
	res, err := db.impl.Get(versionKey, nil)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(res), nil
}

func (db *KVDB) SetDataVersion(version uint64) error {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, version)
	return db.impl.Set(versionKey, buf, nil)
}

func (db *KVDB) GetSnapshot(id, version int64) (string, error) {
	key := genSnapshotKey(id, version)
	res, err := db.impl.Get(key, nil)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (db *KVDB) SetSnapshot(id, version int64, content string) error {
	key := genSnapshotKey(id, version)
	return db.impl.Set(key, []byte(content), nil)
}

func genSnapshotKey(id, version int64) []byte {
	buf := &bytes.Buffer{}
	buf.WriteByte(snapshotKeyCategory)
	binary.Write(buf, binary.BigEndian, id)
	binary.Write(buf, binary.BigEndian, version)
	return buf.Bytes()
}
