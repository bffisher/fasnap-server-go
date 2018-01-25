package main

import (
	"fasnap-server-go/data"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	err := data.Open("./db_files")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	port := flag.Int("port", 8017, "server port")
	flag.Parse()

	go Run(":" + strconv.Itoa(*port))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	log.Println("Got signal:", s)
}
