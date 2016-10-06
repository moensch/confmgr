package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/moensch/confmgr"
	"io/ioutil"
	"os"
	"strings"
)

var (
	logLevel     string
	defaultsPath string
)

func init() {
	flag.StringVar(&logLevel, "d", "warn", "Log level (debug|info|warn|error|fatal)")
	flag.StringVar(&defaultsPath, "p", "", "Directory containing defaults data")
}

func main() {
	flag.Parse()
	if defaultsPath == "" {
		fmt.Println("Must provide defaultsPath")
		flag.Usage()
		os.Exit(1)
	}
	defaultsPath = strings.TrimSuffix(defaultsPath, "/")
	lvl, _ := log.ParseLevel(logLevel)
	log.SetLevel(lvl)

	srv, err := confmgr.NewConfMgr()
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}

	log.Infof("Scanning: %s", defaultsPath)
	files, err := ioutil.ReadDir(defaultsPath)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	b := confmgr.BackendFactory.NewBackend()
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		log.Infof("Loading file: %s", f.Name())
		keyName := f.Name()
		if !strings.HasPrefix(keyName, srv.Config.Main.KeyPrefix) {
			keyName = srv.Config.Main.KeyPrefix + keyName
		}
		log.Infof("Storing key: %s", keyName)
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", defaultsPath, f.Name()))
		if err != nil {
			log.Fatalf("Cannot read file: %s", err)
		}
		err = srv.SaveKeyFromJSON(keyName, data, b)
		if err != nil {
			log.Warnf("Cannot store key: %s", err)
			continue
		}
		log.Info(" Stored!")
	}

}
