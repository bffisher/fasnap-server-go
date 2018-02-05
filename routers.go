package main

import (
	"fasnap-server-go/auth"
	"fasnap-server-go/data"

	"github.com/gin-gonic/gin"
)

func bindRouters(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	engine.GET("/authorization/:user/:password", auth.Authorize())

	engine.Use(auth.Validate())

	engine.GET("/data-version", data.GetDataVersion)

	engine.GET("/snapshot-list/:version", data.GetSnapshotList)

	engine.GET("/snapshot/:date/:version", data.GetSnapshot)

	engine.PUT("/snapshot/:date/:version", data.SaveSnapshot)

	engine.DELETE("/snapshot/:date/:version", data.DeleteSnapshot)
}
