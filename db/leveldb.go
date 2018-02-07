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
	snapshotKeyCategory byte = iota
)

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

//GetCurVersion Get current data version of user
// func (db *KVDB) GetCurVersion(user string) (uint64, error) {
// 	key := genVersionKey(user)
// 	buf, err := db.impl.Get(key, nil)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return binary.BigEndian.Uint64(buf), nil
// }

//NewVersion Return a new data version of user
// func (db *KVDB) NewVersion(user string) (uint64, error) {
// 	key := genVersionKey(user)
// 	buf, err := db.impl.Get(key, nil)
// 	if err != nil {
// 		return 0, err
// 	}

// 	ver := binary.BigEndian.Uint64(buf)
// 	ver++
// 	binary.BigEndian.PutUint64(buf, ver)

// 	err = db.impl.Set(key, buf, nil)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return ver, nil
// }

func (db *KVDB) GetSnapshot(version int64) (string, error) {
	key := genSnapshotKey(version)
	res, err := db.impl.Get(key, nil)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (db *KVDB) SetSnapshot(version int64, content string) error {
	key := genSnapshotKey(version)
	return db.impl.Set(key, []byte(content), nil)
}

func genSnapshotKey(version int64) []byte {
	buf := &bytes.Buffer{}
	buf.WriteByte(snapshotKeyCategory)
	binary.Write(buf, binary.BigEndian, version)
	return buf.Bytes()
}

// func genVersionKey(user string) []byte {
// 	buf := &bytes.Buffer{}
// 	buf.WriteByte(snapshotKeyCategory)
// 	buf.Write([]byte(user))
// 	return buf.Bytes()
// }
