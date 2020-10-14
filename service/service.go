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

type CacheService struct {
	broker       *broker.Broker
	client		 *backendclient.Client
	router       *pat.Router
}

type Router interface {
	Get(path string, handler http.HandlerFunc) *mux.Route
}

type CacheConfiguration struct {
	Configuration *config.Config
	Router        *pat.Router
}

func NewCacheService(cfg *CacheConfiguration) *CacheService {
	return &CacheService{
		broker:       broker.NewBroker(),
		router:       cfg.Router,
	}
}

func (s *CacheService) WithTopic(topic string) *CacheService {
	s.client = backendclient.NewClient(
		"url",
		s.broker,
		http.DefaultClient,
		cache.NewRedisCacheService(
			"tcp",
			"localhost:32768",
			10,
			int64(3600),
		),
		topic,
		logger.NewLogger())
	return s
}

func (s *CacheService) WithPath(path string) *CacheService {
	s.router.Path(path).Methods("GET").HandlerFunc(handlers.NewRequestHandler(s.broker, logger.NewLogger()).HandleRequest)
	return s
}

func (s *CacheService) Start() {
	go s.broker.Run()
}
