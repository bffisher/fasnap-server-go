package data

import (
	"fasnap-server-go/db"
	"fasnap-server-go/entities"
	"fasnap-server-go/errors"
	"io/ioutil"
	"net/http"
	"strconv"

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

//GetSnapshotVersion Get current version for current user
func GetSnapshotVersion(ctx *gin.Context) {
	if user, ok := checkUser(ctx); ok {
		ver, err := sqldb.GetCurSnapshotVersion(user)
		if err != nil {
			respondJSON(ctx, errors.CanntGetData())
		} else {
			respondJSON(ctx, gin.H{"version": ver})
		}
	}
}

//GetSnapshotList Get Snapshot List
func GetSnapshotList(ctx *gin.Context) {
	if user, ok := checkUser(ctx); ok {
		ver, err := strconv.ParseInt(ctx.Param("version"), 10, 64)
		if err != nil {
			respondJSON(ctx, errors.CanntGetData())
			return
		}

		snapshots, err := sqldb.GetSnapshotList(user, ver)
		if err != nil {
			respondJSON(ctx, errors.CanntGetData())
			return
		}

		for index, snapshot := range snapshots {
			if snapshot.Operation == entities.OPSSave {
				content, err := kvdb.GetSnapshot(snapshot.Version)
				if err != nil {
					respondJSON(ctx, errors.CanntGetData())
					return
				}

				snapshots[index].Content = content
			}
		}
		respondJSON(ctx, gin.H{"snapshots": snapshots})
	}
}

//SaveSnapshot Save Snapshot
func SaveSnapshot(ctx *gin.Context) {
	if user, ok := checkUser(ctx); ok {
		date := ctx.Param("date")
		body := ctx.Request.Body
		content, err := ioutil.ReadAll(body)
		if err != nil {
			respondJSON(ctx, errors.CanntGetParam("body"))
			return
		}

		ver, err := sqldb.SaveSnapshot(user, date)
		if err != nil {
			respondJSON(ctx, errors.CanntSaveData())
			return
		}

		err = kvdb.SetSnapshot(ver, string(content))
		if err != nil {
			respondJSON(ctx, errors.CanntSaveData())
			return
		}

		respondJSON(ctx, gin.H{"version": ver})
	}
}

//DeleteSnapshot Delete Snapshot
func DeleteSnapshot(ctx *gin.Context) {
	if user, ok := checkUser(ctx); ok {
		date := ctx.Param("date")
		ver, err := sqldb.DeleteSnapshot(user, date)
		if err != nil {
			respondJSON(ctx, errors.CanntSaveData())
			return
		}

		respondJSON(ctx, gin.H{"version": ver})
	}
}

//Close database
func Close() {
	sqldb.Close()
	kvdb.Close()
}

func checkUser(ctx *gin.Context) (string, bool) {
	if user, ok := ctx.Get("user"); ok {
		return user.(string), ok
	}

	respondJSON(ctx, errors.InvalidUser())
	return "", false
}

func respondJSON(ctx *gin.Context, json gin.H) {
	ctx.JSON(http.StatusOK, json)
}
