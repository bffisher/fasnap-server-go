package data

import (
	"fasnap-server-go/db"
	"fasnap-server-go/errors"

	"github.com/gin-gonic/gin"
)

var sqldb *db.SQLDB
var kvdb *db.KVDB

//Open and init database
func Open(rootPath string) error {
	var err error

	sqldb = &db.SQLDB{}
	err = sqldb.Open(rootPath + "/sqldb")
	if err != nil {
		return err
	}

	kvdb = &db.KVDB{}
	return kvdb.Open(rootPath + "/kvdb")
}

//get data version
func GetDataVersion(ctx *gin.Context) {
	var res gin.H
	ver, err := kvdb.GetDataVersion()
	if err != nil {
		res = errors.CanntGetData()
	} else {
		res = gin.H{"version": ver}
	}

	ctx.JSON(200, res)
}

func GetSnapshotList(ctx *gin.Context) {
}

func GetSnapshot(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"version": ctx.Param("version"), "date": ctx.Param("date")})
}

func SaveSnapshot(ctx *gin.Context) {
}

func DeleteSnapshot(ctx *gin.Context) {
}

//Close database
func Close() {
	sqldb.Close()
	kvdb.Close()
}
