package main

import (
	"fasnap-server-go/data"
	"log"
	"os"
	"os/signal"
)

func main() {
	err := data.Open("./db_files")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	var port string
	if arg, ok := findArg("port"); ok {
		port = ":" + arg
	} else {
		port = ":8017"
	}

	go Run(port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	log.Println("Got signal:", s)
}

func findArg(name string) (string, bool) {
	optionKey := "-" + name
	for i, arg := range os.Args {
		if arg == optionKey && (i+1) < len(os.Args) {
			return os.Args[i+1], true
		}
	}
	return "", false
}
