package service

import (
	"github.com/companieshouse/chs-streaming-api-cache/config"
	cachehandlers "github.com/companieshouse/chs-streaming-api-cache/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/pat"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"testing"
)


type mockBackendClient struct {
	mock.Mock
}

func (s *mockBackendClient) Create(key string, delta string, offset int64) error {
	args := s.Called(key, delta, offset)
	return args.Error(0)
}

func TestCreateNewService(t *testing.T) {
	Convey("When a new service instance is constructed", t, func() {
		configuration := &CacheConfiguration{
			Configuration: &config.Config{
				RedisUrl:             "localhost:6379",
				CacheExpiryInSeconds: 2,
			},
			Router: pat.New(),
		}
		actual := NewCacheService(configuration)
		Convey("Then a new service instance reference should be returned", func() {
			So(actual, ShouldNotBeNil)
			So(actual.broker, ShouldNotBeNil)
			So(actual.router, ShouldEqual, configuration.Router)
		})
	})
}

func TestBindTopic(t *testing.T) {
	Convey("Given a new service instance has been constructed", t, func() {
		configuration := &CacheConfiguration{
			Configuration: &config.Config{
				RedisUrl:             "localhost:6379",
				CacheExpiryInSeconds: 2,
			},
			Router: pat.New(),
		}

		service := NewCacheService(configuration)
		Convey("When a topic is bound to it", func() {
			actual := service.WithTopic("topic")
			Convey("Then the topic should be added to the service", func() {
				So(actual, ShouldEqual, service)
				So(actual.topic, ShouldEqual, "topic")
			})
		})
	})
}

func TestAttachRequestHandler(t *testing.T) {
	Convey("Given a new service instance has been constructed", t, func() {
		configuration := &CacheConfiguration{
			Configuration: &config.Config{
				RedisUrl:             "localhost:6379",
				CacheExpiryInSeconds: 2,
			},
			Router: pat.New(),
		}
		service := NewCacheService(configuration)
		Convey("When a request handler is attached to it", func() {
			actual := service.WithPath("/path")
			Convey("Then a new request handler should be allocated to the service", func() {
				So(actual, ShouldEqual, service)
				_ = configuration.Router.Walk(func(r *mux.Route, o *mux.Router, u []*mux.Route) error {
					path, _ := r.GetPathTemplate()
					methods, _ := r.GetMethods()
					handler := r.GetHandler()
					So(path, ShouldEqual, "/path")
					So(methods, ShouldResemble, []string{"GET"})
					So(handler, ShouldEqual, (&cachehandlers.RequestHandler{}).HandleRequest)
					return nil
				})
			})
		})
	})
}
