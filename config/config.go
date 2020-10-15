package config

import "github.com/ian-kent/gofigure"

type Config struct {
	gofigure             interface{} `order:"env,flag"`
	BindAddress          string      `env:"BIND_ADDRESS"            flag:"bind-address"`
	CertFile             string      `env:"CERT_FILE"               flag:"cert-file"`
	KeyFile              string      `env:"KEY_FILE"                flag:"key-file"`
	BackEndUrl           string      `env:"STREAMING_BACKEND_URL"   flag:"streaming_backend_url"`
	RedisUrl             string      `env:"REDIS_URL"               flag:"redis-url"`
	RedisPoolSize        int         `env:"REDIS_POOL_SIZE"         flag:"redis_pool_size"`
	CacheExpiryInSeconds int64       `env:"CACHE_EXPIRY_IN_SECONDS" flag:"cache-expiry-in-seconds"`
}

// ServiceConfig returns a ServiceConfig interface for Config.
func (c Config) ServiceConfig() ServiceConfig {
	return ServiceConfig{c}
}

// ServiceConfig wraps Config to implement service.Config.
type ServiceConfig struct {
	Config
}

func (cfg ServiceConfig) BindAddr() string {
	return cfg.Config.BindAddress
}

func (cfg ServiceConfig) CertFile() string {
	return cfg.Config.CertFile
}

func (cfg ServiceConfig) KeyFile() string {
	return cfg.Config.KeyFile
}

func (cfg ServiceConfig) Namespace() string {
	return "chs-streaming-api-cache"
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
