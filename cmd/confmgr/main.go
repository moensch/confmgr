package main

import (
	"github.com/moensch/confmgr"
	"log"
)

func main() {
	srv, err := confmgr.NewConfMgr()
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
	log.Println("initialized")

	srv.Run()
}
