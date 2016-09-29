package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

type ConfigMgrConfig struct {
	Listen   ListenConfig `toml:"listen"`
	Main     MainConfig   `toml:"main"`
	Backends map[string]BackendConfig
}

type BackendConfig struct {
	Port    int
	Address string
}

type ListenConfig struct {
	Port    int
	Address string
}

type MainConfig struct {
	KeyPaths  []string `toml:"key_paths"`
	KeyPrefix string   `toml:"key_prefix"`
	HdrPrefix string   `toml:"hdr_prefix"`
}

func LoadConfig(c *ConfigMgrConfig, path string) error {
	log.Infof("Reading config from: '%s'", path)

	_, err := toml.DecodeFile(path, &c)

	return err
}
