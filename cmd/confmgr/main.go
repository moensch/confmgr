package main

import (
	"github.com/moensch/confmgr"
	"log"
)

func main() {
	srv, err := confmgr.NewConfMgr()
	if err != nil {
		log.Printf("Cannot start server: %s", err)
		log.Fatal("exiting")
	}
	log.Println("initialized")

	srv.Run()
}
