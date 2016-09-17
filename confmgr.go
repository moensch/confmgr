package confmgr

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/moensch/confmgr/backends"
	"github.com/moensch/confmgr/backends/redis"
	"net/http"
	"os"
)

type ConfMgr struct {
	Config       ConfigMgrConfig
	Backend      backend.ConfigBackend
	Router       *mux.Router
	RequestScope map[string]string
}

var (
	backendFactory backend.ConfigBackendFactory
)

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
		log.Debugf("Checking for config in %s", configpath)
		if _, err := os.Stat(configpath); err == nil {
			err = confmgr.LoadConfig(configpath)
			if err != nil {
				log.Fatalf("Cannot load config: %s", err)
			}
			break
		}
	}

	backendFactory = redis.NewFactory()
	confmgr.Router = confmgr.NewRouter()

	return confmgr, err
}

func (c *ConfMgr) Run() {
	listenAddr := fmt.Sprintf("%s:%d", c.Config.Listen.Address, c.Config.Listen.Port)
	log.Infof("Listening on: %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, c.Router))
}
