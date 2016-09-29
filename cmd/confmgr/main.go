package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/moensch/confmgr"
)

var (
	logLevel string
)

func init() {
	flag.StringVar(&logLevel, "d", "warn", "Log level (debug|info|warn|error|fatal)")
}

func main() {
	flag.Parse()

	lvl, _ := log.ParseLevel(logLevel)
	log.SetLevel(lvl)
	srv, err := confmgr.NewConfMgr()
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
	log.Info("initialized")

	srv.Run()
}
