package main

import (
	"fasnap-server-go/auth"

	"github.com/gin-gonic/gin"
)

func Run(addr ...string) {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/authorization/:user/:password", auth.Authorize())

	router.Use(auth.Validate())

	router.GET("/snapshot-data-log/:version", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"version": "1"})
	})

	router.GET("/snapshot-list/:date/:version", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"version": ctx.Param("version")})
	})

	router.GET("/snapshot/:date/:version", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"version": ctx.Param("version"), "date": ctx.Param("date")})
	})

	router.PUT("/snapshot/:date/:version", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"version": ctx.Param("version")})
	})

	router.DELETE("/snapshot/:date/:version", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"version": ctx.Param("version")})
	})

	router.Run(addr...)
}
