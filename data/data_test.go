package data

import (
	"fasnap-server-go/entities"
	"fasnap-server-go/test"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	dbTestPath    = "../db_files"
	dbTestContent = `{"date":"2018-01-11", "items":[{"rist":"low", amount:1222}, {"rist":"high", amount:2221}]}`
)

var testEngine *gin.Engine

func init() {
	os.Remove(dbTestPath + "/sqldb")
	os.RemoveAll(dbTestPath + "/kvdb")

	testEngine = gin.Default()

	testEngine.Use(func(ctx *gin.Context) {
		ctx.Set("user", "admin")
		ctx.Next()
	})

	testEngine.GET("/snapshot-version", GetSnapshotVersion)

	testEngine.GET("/snapshot-list/:version", GetSnapshotList)

	testEngine.PUT("/snapshot/:date", SaveSnapshot)

	testEngine.DELETE("/snapshot/:date", DeleteSnapshot)
}

func Test_open(t *testing.T) {
	err := Open(dbTestPath)
	if err != nil {
		t.Error(err)
	}
}

func Test_snapshot_version(t *testing.T) {
	res, err := test.HttpGetJson("/snapshot-version", testEngine, nil)
	if err != nil {
		t.Error(err)
	}

	test.VerifyRespondNoError(t, res)

	if ver := int64(res["version"].(float64)); ver != 0 {
		t.Error("version should 0", res)
	}
}

func Test_GetSnapshotList_0(t *testing.T) {
	res, err := test.HttpGetJson("/snapshot-list/0", testEngine, nil)
	if err != nil {
		t.Error(err)
	}

	test.VerifyRespondNoError(t, res)

	if snapshots := res["snapshots"].([]interface{}); len(snapshots) != 0 {
		t.Error("snapshots count should 0", snapshots)
	}
}

func Test_SaveSnapshot(t *testing.T) {
	res, err := test.HttpPutJson("/snapshot/2018-02-06", testEngine, dbTestContent, nil)
	if err != nil {
		t.Error(err)
	}
	test.VerifyRespondNoError(t, res)

	ver := int64(res["version"].(float64))
	if ver == 0 {
		t.Error("version should >0", ver)
	}

	verifySnapshot(t, ver, entities.OPSSave, dbTestContent)
}

func Test_DeleteSnapshot(t *testing.T) {
	res, err := test.HttpDeleteJson("/snapshot/2018-02-06", testEngine, nil)
	if err != nil {
		t.Error(err)
	}
	ver := int64(res["version"].(float64))
	if ver == 0 {
		t.Error("version should >0", ver)
	}
	verifySnapshot(t, ver, entities.OPSDel, "")
}

func Test_close(t *testing.T) {
	Close()
}

func verifySnapshot(t *testing.T, ver int64, operation int, content string) {
	res, err := test.HttpGetJson("/snapshot-list/"+strconv.FormatInt(ver-1, 10), testEngine, nil)
	if err != nil {
		t.Error(err)
	}
	test.VerifyRespondNoError(t, res)

	snapshots := res["snapshots"].([]interface{})
	if len(snapshots) == 0 {
		t.Error("snapshots count should >0", snapshots)
	}
	for _, item := range snapshots {
		snapshot := gin.H(item.(map[string]interface{}))

		if gVer, ok := snapshot["Version"]; !ok || int64(gVer.(float64)) != ver {
			t.Error("snapshot version error", gVer)
		}

		if gOps, ok := snapshot["Operation"]; !ok || int(gOps.(float64)) != operation {
			t.Error("snapshot operation error", gOps)
		}

		if content != "" {
			if gContent, ok := snapshot["Content"]; !ok || gContent.(string) != content {
				t.Error("snapshot content error", gContent)
			}
		}
	}
}
