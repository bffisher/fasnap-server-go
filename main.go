package main

import (
	"fasnap-server-go/data"
	"log"
	"os"
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

	Run(port)
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
