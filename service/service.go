package service

import (
	"github.com/companieshouse/chs-streaming-api-cache/broker"
	"github.com/companieshouse/chs-streaming-api-cache/cache"
	backendclient "github.com/companieshouse/chs-streaming-api-cache/client"
	"github.com/companieshouse/chs-streaming-api-cache/config"
	"github.com/companieshouse/chs-streaming-api-cache/handlers"
	"github.com/companieshouse/chs-streaming-api-cache/logger"
	"github.com/companieshouse/chs-streaming-api-cache/mapper"
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
	username   string
	redisCfg   RedisConfig
	myMapper   *mapper.ConfigurationPathMapper
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
		username:   cfg.Configuration.ChsApiKey,
		redisCfg: RedisConfig{
			redisUrl:        cfg.Configuration.RedisUrl,
			expiryInSeconds: cfg.Configuration.CacheExpiryInSeconds,
			poolSize:        cfg.Configuration.RedisPoolSize,
		},
		myMapper: mapper.New(cfg.Configuration),
	}
}

func (s *CacheService) WithTopic(topic string) *CacheService {
	s.topic = topic
	return s
}

func (s *CacheService) WithPath(path string) *CacheService {
	s.path = path
	return s
}

func (s *CacheService) Initialise() *CacheService {
	cfg := s.redisCfg

	cacheClient := cache.NewRedisCacheService(
		network,
		cfg.redisUrl,
		cfg.poolSize,
		cfg.expiryInSeconds,
	)

	backendPath, err := s.myMapper.GetBackendPathForPath(s.path)
	if err != nil {
		// default to same Url
		backendPath = s.path
	}
	s.client = backendclient.NewClient(
		s.backendURL,
		backendPath,
		s.broker,
		http.DefaultClient,
		s.username,
		cacheClient,
		s.topic,
		logger.NewLogger())

	s.router.Path(s.path).Methods("GET").HandlerFunc(handlers.NewRequestHandler(s.broker, cacheClient, logger.NewLogger(), s.topic).HandleRequest)
	return s
}

func (s *CacheService) Start() {
	go s.client.Run()
	go s.broker.Run()
}
