package confmgr

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/moensch/confmgr/backends"
	"github.com/moensch/confmgr/backends/redis"
	"log"
	"net/http"
	"os"
)

type ConfMgr struct {
	Config       ConfigMgrConfig
	Backend      backend.ConfigBackend
	Router       *mux.Router
	RequestScope map[string]string
}

func NewConfMgr() (*ConfMgr, error) {
	confmgr := &ConfMgr{
		Config: ConfigMgrConfig{
			Listen: listenConfig{
				Port:    8080,
				Address: "0.0.0.0",
			},
		},
	}

	configLocations := []string{
		"/etc/confmgr.toml",
		"/confmgr.toml",
		"confmgr.toml",
	}

	var err error
	// Parse config if exists in any of our search locations
	for _, configpath := range configLocations {
		log.Printf("Checking for config in %s\n", configpath)
		if _, err := os.Stat(configpath); err == nil {
			err = confmgr.LoadConfig(configpath)
			if err != nil {
				log.Fatalf("Cannot load config: %s", err)
			}
			break
		}
	}

	confmgr.Backend = redis.Init()
	confmgr.Router = confmgr.NewRouter()

	return confmgr, err
}

func (c *ConfMgr) Run() {
	//c.Backend = redis.Init(c.Config.Backends["redis"])
	listenAddr := fmt.Sprintf("%s:%d", c.Config.Listen.Address, c.Config.Listen.Port)
	log.Printf("Listening on: %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, c.Router))
}
