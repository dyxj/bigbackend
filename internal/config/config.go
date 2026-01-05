package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPServerConfig *HTTPServerConfig `env:",init"`
	DBConfig         *DBConfig         `env:",init" envPrefix:"DB_"`
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
