package auth

import (
	"fasnap-server-go/errors"
	"fasnap-server-go/test"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()

	router.GET("/authorization/:user/:password", Authorize())

	router.Use(Validate())

	router.GET("/snapshot", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"date": "2018-01-23"})
	})
}

func Test_Authorization_NoPwd(t *testing.T) {
	test.VerifyHttpGetJsonFail(t, router, "/authorization/:admin", errors.CanntFindToken(), nil)
}

func Test_Authorization_NoUse(t *testing.T) {
	test.VerifyHttpGetJsonFail(t, router, "/authorization", errors.CanntFindToken(), nil)
}

func Test_Authorization_UseError(t *testing.T) {
	test.VerifyHttpGetJsonFail(t, router, "/authorization/adf/admin", errors.IncorrectUserPwd(), nil)
}

func Test_Authorization_PwdError(t *testing.T) {
	test.VerifyHttpGetJsonFail(t, router, "/authorization/admin/admi", errors.IncorrectUserPwd(), nil)
}

func Test_Authorization_Sucess(t *testing.T) {
	res, err := test.HttpGetJson("/authorization/admin/admin", router, nil)
	if err != nil {
		t.Error(err)
	}

	test.VerifyRespondNoError(t, res)
}

func Test_ValidateAuth_NoToken(t *testing.T) {
	test.VerifyHttpGetJsonFail(t, router, "/snapshot", errors.CanntFindToken(), nil)
}

func Test_ValidateAuth_IncorrectToken(t *testing.T) {
	test.VerifyHttpGetJsonFail(t, router, "/snapshot", errors.NotEvenToken(), func(req *http.Request) {
		req.Header.Set("Authorization", "Bsic.....")
	})
}
func Test_ValidateAuth_ExpiredToken(t *testing.T) {
	tokenStr, err := getSignedTokenStr("admin", time.Now().Add(-time.Hour).Unix())
	if err != nil {
		t.Error("Get Signed Token Fail", err)
	}

	test.VerifyHttpGetJsonFail(t, router, "/snapshot", errors.ExpiredOrNotActiveToken(), func(req *http.Request) {
		req.Header.Set("Authorization", tokenStr)
	})
}

func Test_ValidateAuth_CorrectToken(t *testing.T) {
	tokenStr, err := getSignedTokenStr("admin", time.Now().Add(time.Second).Unix())
	if err != nil {
		t.Error("Get Signed Token Fail", err)
	}

	res, err := test.HttpGetJson("/snapshot", router, func(req *http.Request) {
		req.Header.Set("Authorization", tokenStr)
	})
	if err != nil {
		t.Error(err)
	}

	test.VerifyRespondNoError(t, res)
}
