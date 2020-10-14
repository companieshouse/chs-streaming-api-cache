package service

import (
	"github.com/companieshouse/chs-streaming-api-cache/broker"
	"github.com/companieshouse/chs-streaming-api-cache/cache"
	backendclient "github.com/companieshouse/chs-streaming-api-cache/client"
	"github.com/companieshouse/chs-streaming-api-cache/config"
	"github.com/companieshouse/chs-streaming-api-cache/handlers"
	"github.com/companieshouse/chs-streaming-api-cache/logger"
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	"net/http"
)

const network = "tcp"

type CacheService struct {
	broker     *broker.Broker
	client     *backendclient.Client
	router     *pat.Router
	topic      string
	path       string
	backendURL string
	redisCfg   RedisConfig
}

type Router interface {
	Get(path string, handler http.HandlerFunc) *mux.Route
}

type CacheConfiguration struct {
	Configuration *config.Config
	Router        *pat.Router
}

type RedisConfig struct {
	redisUrl        string
	expiryInSeconds int64
	poolSize        int
}

func NewCacheService(cfg *CacheConfiguration) *CacheService {
	return &CacheService{
		broker:     broker.NewBroker(),
		router:     cfg.Router,
		backendURL: cfg.Configuration.BackEndUrl,
		redisCfg: RedisConfig{
			redisUrl:        cfg.Configuration.RedisUrl,
			expiryInSeconds: cfg.Configuration.CacheExpiryInSeconds,
			poolSize:        cfg.Configuration.RedisPoolSize,
		},
	}
}

func (s *CacheService) WithTopic(topic string) *CacheService {
	s.topic = topic
	return s
}

func (s *CacheService) WithPath(path string) *CacheService {
	s.path = path
	s.router.Path(path).Methods("GET").HandlerFunc(handlers.NewRequestHandler(s.broker, logger.NewLogger()).HandleRequest)
	return s
}

func (s *CacheService) Initialise() *CacheService {
	cfg := s.redisCfg
	s.client = backendclient.NewClient(
		s.backendURL+s.path,
		s.broker,
		http.DefaultClient,
		cache.NewRedisCacheService(
			network,
			cfg.redisUrl,
			cfg.poolSize,
			cfg.expiryInSeconds,
		),
		s.topic,
		logger.NewLogger())
	return s
}

func (s *CacheService) Start() {
	go s.broker.Run()
}
