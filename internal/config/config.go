package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPServerConfig *HTTPServerConfig `env:",init"`
	DBConfig         *DBConfig         `env:",init"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: true,
	})
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadDBConfig() (*DBConfig, error) {
	var dbCfg DBConfig
	err := env.ParseWithOptions(&dbCfg, env.Options{
		RequiredIfNoDef: true,
		Prefix:          "DB_",
	})
	if err != nil {
		return nil, err
	}
	return &dbCfg, nil
}
