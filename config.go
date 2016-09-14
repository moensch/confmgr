package confmgr

import (
	"github.com/BurntSushi/toml"
	"log"
)

type ConfigMgrConfig struct {
	Listen   listenConfig `toml:"listen"`
	Main     mainConfig   `toml:"main"`
	Backends map[string]backendConfig
}

type backendConfig struct {
	Port    int
	Address string
}

type listenConfig struct {
	Port    int
	Address string
}

type mainConfig struct {
	KeyPaths  []string `toml:"key_paths"`
	KeyPrefix string   `toml:"key_prefix"`
	HdrPrefix string   `toml:"hdr_prefix"`
}

func (c *ConfMgr) LoadConfig(path string) error {
	log.Printf("Reading config from: '%s'", path)

	_, err := toml.DecodeFile(path, &c.Config)

	return err
}
