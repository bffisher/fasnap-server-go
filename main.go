package main

import (
	"fasnap-server-go/data"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	err := data.Open("./db_files")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	router.Use(gin.BasicAuth(gin.Accounts{"admin": "admin"}))

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

	port := flag.String("port", ":8017", "http listen port")
	router.Run(*port)
}
