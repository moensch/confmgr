package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/moensch/confmgr"
)

func main() {
	log.SetLevel(log.WarnLevel)
	srv, err := confmgr.NewConfMgr()
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
	log.Info("initialized")

	srv.Run()
}
