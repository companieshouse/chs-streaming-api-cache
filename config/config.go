package config

import "github.com/ian-kent/gofigure"

type Config struct {
	gofigure             interface{} `order:"env,flag"`
	BindAddress          string      `env:"BIND_ADDRESS" flag:"bind-address"`
	RedisUrl             string      `env:"REDIS_URL" flag:"redis-url"`
	CacheExpiryInSeconds string      `env:"CACHE_EXPIRY_IN_SECONDS" flag:"cache-expiry-in-seconds"`
}

var config *Config

func Get() (*Config, error) {
	if config == nil {
		config = &Config{}
		if err := gofigure.Gofigure(config); err != nil {
			return nil, err
		}
	}
	return config, nil
}
