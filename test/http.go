package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func HttpGetJson(uri string, engine *gin.Engine, reqHandle func(*http.Request)) (gin.H, error) {
	req := httptest.NewRequest(http.MethodGet, uri, nil)

	if reqHandle != nil {
		reqHandle(req)
	}

	return sendHTTPJSON(engine, req)
}

func HttpPutJson(uri string, engine *gin.Engine, content string, reqHandle func(*http.Request)) (gin.H, error) {
	req := httptest.NewRequest(http.MethodPut, uri, strings.NewReader(content))
	defer req.Body.Close()
	if reqHandle != nil {
		reqHandle(req)
	}
	return sendHTTPJSON(engine, req)
}

func HttpDeleteJson(uri string, engine *gin.Engine, reqHandle func(*http.Request)) (gin.H, error) {
	req := httptest.NewRequest(http.MethodDelete, uri, nil)

	if reqHandle != nil {
		reqHandle(req)
	}

	return sendHTTPJSON(engine, req)
}

func sendHTTPJSON(engine *gin.Engine, req *http.Request) (gin.H, error) {
	w := httptest.NewRecorder()

	engine.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	body, _ := ioutil.ReadAll(result.Body)

	var jsonResult gin.H
	if err := json.Unmarshal(body, &jsonResult); err != nil {

		return nil, err
	}

	return jsonResult, nil
}

func VerifyHttpGetJsonFail(t *testing.T, engine *gin.Engine, url string, errCompared gin.H, reqHandle func(*http.Request)) {
	res, err := HttpGetJson(url, engine, reqHandle)
	if err != nil {
		t.Error(err)
	}

	VerifyRespondError(t, res, errCompared)
}

func VerifyRespondError(t *testing.T, res gin.H, err gin.H) {
	if code, ok := res["error"]; ok && code == err["error"] {
		t.Log("Sucess.")
	} else {
		t.Error("Fail", res)
	}
}

func VerifyRespondNoError(t *testing.T, res gin.H) {
	if _, ok := res["error"]; ok {
		t.Error("Fail", res)
	} else {
		t.Log("Sucess.")
	}
}
