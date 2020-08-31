package config

import (
	"fmt"
	"os"

	"github.com/rfashwal/scs-utilities/config"
)

type registerConfig struct {
	config.Manager
	dbPath string
}

var instance *registerConfig

func Config() *registerConfig {
	if instance == nil {
		instance = new(registerConfig)
		instance.Init()
		if dbPath, err := os.LookupEnv("DB_PATH"); !err {
			panic(fmt.Sprintf("set DP_PATH and try again"))
		} else {
			instance.dbPath = dbPath
		}
	}
	return instance
}

func (c *registerConfig) DBPath() string {
	return c.dbPath
}
