package main

import (
	"fasnap-server-go/data"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := data.Open("./db_files")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	port := flag.Int("port", 8017, "server port")
	flag.Parse()

	engine := gin.Default()
	bindRouters(engine)
	go engine.Run(":" + strconv.Itoa(*port))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	log.Println("Got signal:", s)
}
